#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

prepare_codeql() {
  copy_file scripts/codeql
  write_file main.go <<'EOF'
package main
EOF
  write_file pkg/lib.go <<'EOF'
package lib
EOF
  write_file codeql/no-code-comments/qlpack.yml <<'EOF'
name: xquik/no-code-comments
EOF
  write_file codeql/no-code-comments-test/NoCodeComments/case.go.fixture <<'EOF'
package comments
EOF
  fake codeql <<'EOF'
#!/bin/sh
echo "$*" >>"$FAKE_LOGS/codeql"
option() {
  name="$1"
  shift
  for argument in "$@"; do
    case "$argument" in
      "$name"=*) echo "${argument#"$name"=}" ;;
    esac
  done
}
case "$1 $2" in
  "version --format=terse")
    echo "2.27.2"
    ;;
  "pack download")
    if [ -n "${FAKE_STALL:-}" ]; then
      sh -c 'echo $$ >"$FAKE_LOGS/stalled-pid"; exec sleep 30'
      exit 0
    fi
    [ -z "${FAKE_DOWNLOAD_FAILS:-}" ] || exit 1
    pack="$(option --dir "$@")/codeql/go-queries/1.6.12"
    mkdir -p "$pack/codeql-suites"
    echo "version: 1.6.12" >"$pack/qlpack.yml"
    [ -n "${FAKE_INCOMPLETE_PACK:-}" ] || : >"$pack/codeql-suites/go-security-extended.qls"
    ;;
  "pack install")
    [ -z "${FAKE_INSTALL_FAILS:-}" ] || exit 1
    ;;
  "test run")
    [ -f "$3/NoCodeComments/case.go" ] && [ ! -e "$3/NoCodeComments/case.go.fixture" ] || exit 98
    [ -z "${FAKE_QUERY_TEST_FAILS:-}" ] || exit 1
    ;;
  "database create")
    source_root="$(option --source-root "$@")"
    echo "$source_root" >"$FAKE_LOGS/source-root"
    mkdir -p "$3"
    echo "built" >"$source_root/build-output.txt"
    files=""
    for file in main.go pkg/lib.go; do
      [ "$file" = "${FAKE_SKIPPED_FILE:-}" ] || files="$files ${source_root#/}/$file"
    done
    (cd / && zip -q "$3/src.zip" $files)
    ;;
  "database analyze")
    output="$(option --output "$@")"
    results=""
    notifications=""
    diagnostic="go/diagnostics/extraction-errors"
    if [ -n "${FAKE_RESULTS:-}" ]; then
      results='{"ruleId":"go/sql-injection","locations":[{"physicalLocation":{"artifactLocation":{"uri":"pkg/lib.go"},"region":{"startLine":7}}}]}'
    fi
    if [ -n "${FAKE_UNREAD_FILE:-}" ]; then
      notifications='{"message":{"text":"Extraction failed in '"$FAKE_UNREAD_FILE"' with error undefined: missingName"},"level":"error","descriptor":{"id":"go/diagnostics/extraction-errors"}}'
    fi
    [ -z "${FAKE_NO_EXTRACTION_DIAGNOSTIC:-}" ] || diagnostic="go/diagnostics/successfully-extracted-files"
    echo '{"runs":[{"tool":{"driver":{"name":"CodeQL","notifications":[{"id":"'"$diagnostic"'"}]}},"invocations":[{"toolExecutionNotifications":['"$notifications"'],"executionSuccessful":true}],"results":['"$results"']}]}' >"$output"
    ;;
  *)
    echo "unexpected codeql $*" >&2
    exit 97
    ;;
esac
EOF
  fake go <<'EOF'
#!/bin/sh
[ "$*" = "list -json ./..." ] || exit 97
here="$(pwd -P)"
echo "{\"Dir\": \"$here\", \"GoFiles\": [\"main.go\"]}"
echo "{\"Dir\": \"$here/pkg\", \"GoFiles\": [\"lib.go\"]}"
EOF
  commit_repo
}

test_codeql_scans_a_copy_with_the_pinned_pack_and_leaves_the_working_tree_alone() {
  prepare_codeql

  run_in_repo scripts/codeql

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "CodeQL read all 2 Go files and found no results." "stdout"
  assert_equal "" "$(sandbox_git status --porcelain)" "the working tree changes"
  assert_contains "$(log_of source-root)" "$sandbox/tmp/xquik-cli-codeql." "the scanned folder"
  assert_contains "$(log_of codeql)" "pack download codeql/go-queries@1.6.12 --dir=$sandbox/tmp/" "the pack download"
  assert_contains "$(log_of codeql)" "database analyze" "the CodeQL calls"
  assert_contains "$(log_of codeql)" " codeql/go-queries@1.6.12:codeql-suites/go-security-extended.qls codeql/no-code-comments --additional-packs=$sandbox/tmp/" "the analysis"
  assert_contains "$(log_of codeql)" "pack install codeql/no-code-comments --common-caches=$sandbox/tmp/" "the comment query install"
  assert_contains "$(log_of codeql)" "test run $sandbox/tmp/" "the comment query test"
  assert_equal 6 "$(grep -c -- "--common-caches=$sandbox/tmp/xquik-cli-codeql." "$sandbox_logs/codeql")" "the calls with a private cache"
  assert_equal "" "$(ls "$sandbox/tmp")" "the leftover temporary files"
}

test_a_failing_comment_query_test_stops_the_run() {
  prepare_codeql
  run_env=("FAKE_QUERY_TEST_FAILS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The comment query test failed. Fix codeql/no-code-comments or the expected results in codeql/no-code-comments-test, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of codeql)" "database analyze" "the CodeQL calls"
}

test_a_codeql_result_fails_the_run_and_is_printed() {
  prepare_codeql
  run_env=("FAKE_RESULTS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "go/sql-injection at pkg/lib.go:7
CodeQL found 1 results. Fix each one above, then rerun." "$run_stderr" "stderr"
}

test_a_failed_comment_query_install_fails_the_run() {
  prepare_codeql
  run_env=("FAKE_INSTALL_FAILS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "CodeQL could not install the libraries of codeql/no-code-comments. Check the network and GitHub access, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of codeql)" "database analyze" "the CodeQL calls"
}

test_a_go_file_missing_from_the_database_fails_the_run() {
  prepare_codeql
  run_env=("FAKE_SKIPPED_FILE=pkg/lib.go")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "pkg/lib.go
CodeQL missed the Go files above. Fix the build so it compiles every package." "$run_stderr" "stderr"
  assert_not_contains "$(log_of codeql)" "database analyze" "the CodeQL calls"
}

test_a_go_file_that_codeql_could_not_read_fails_the_run_and_is_named() {
  prepare_codeql
  run_env=("FAKE_UNREAD_FILE=pkg/lib.go")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Extraction failed in pkg/lib.go with error undefined: missingName
CodeQL could not read the Go code named above, so the scan is incomplete. Fix that code so Go compiles it, then rerun." "$run_stderr" "stderr"
}

test_a_scan_without_the_extraction_diagnostic_fails_the_run() {
  prepare_codeql
  run_env=("FAKE_NO_EXTRACTION_DIAGNOSTIC=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The CodeQL scan lacks go/diagnostics/extraction-errors, so it cannot show that it read every file. Make the suite in scripts/codeql run that diagnostic query." "$run_stderr" "stderr"
}

test_a_stalled_pack_download_stops_after_the_limit_and_ends_the_download() {
  local stalled_pid
  prepare_codeql
  run_env=("FAKE_STALL=1" "CODEQL_DOWNLOAD_SECONDS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The CodeQL pack download stalled for 1 seconds. Check the network, then rerun." "$run_stderr" "stderr"
  stalled_pid="$(log_of stalled-pid)"
  if kill -0 "$stalled_pid" 2>/dev/null; then
    kill "$stalled_pid"
    fail_test "The stalled download process $stalled_pid was still running."
  fi
}

test_a_failed_pack_download_stops_the_run() {
  prepare_codeql
  run_env=("FAKE_DOWNLOAD_FAILS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "CodeQL could not download codeql/go-queries 1.6.12. Check the network and GitHub access, then rerun." "$run_stderr" "stderr"
}

test_an_incomplete_pack_stops_the_run() {
  prepare_codeql
  run_env=("FAKE_INCOMPLETE_PACK=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The downloaded codeql/go-queries pack is incomplete. Rerun to download it again." "$run_stderr" "stderr"
}

test_a_download_limit_that_is_not_a_number_stops_the_run() {
  prepare_codeql
  run_env=("CODEQL_DOWNLOAD_SECONDS=soon")

  run_in_repo scripts/codeql

  assert_equal 2 "$run_status" "the exit code"
  assert_equal "CODEQL_DOWNLOAD_SECONDS must be a whole number of seconds, such as 300." "$run_stderr" "stderr"
  assert_equal "" "$(log_of codeql)" "the CodeQL calls"
}

run_tests
