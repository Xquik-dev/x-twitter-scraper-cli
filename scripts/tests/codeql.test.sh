#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

# Runs scripts/codeql in a throwaway repository. A fake CodeQL CLI downloads
# a fake pack, builds a fake database and writes a SARIF report. A fake go
# command lists the Go files of the copy. No test reaches the CodeQL registry.

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
  fake codeql <<'EOF'
#!/bin/sh
# Stands in for the CodeQL CLI. FAKE_STALL makes the pack download hang,
# FAKE_DOWNLOAD_FAILS makes it fail, FAKE_INCOMPLETE_PACK leaves out the
# suite, FAKE_SKIPPED_FILE leaves a file out of the database and FAKE_RESULTS
# adds 1 result to the report.
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
    if [ -n "${FAKE_RESULTS:-}" ]; then
      echo '{"runs":[{"results":[{"ruleId":"go/sql-injection","locations":[{"physicalLocation":{"artifactLocation":{"uri":"pkg/lib.go"},"region":{"startLine":7}}}]}]}]}' >"$output"
    else
      echo '{"runs":[{"results":[]}]}' >"$output"
    fi
    ;;
  *)
    echo "unexpected codeql $*" >&2
    exit 97
    ;;
esac
EOF
  fake go <<'EOF'
#!/bin/sh
# Stands in for go list -json. It reports the 2 Go files of the copy.
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
  assert_contains "$(log_of codeql)" " codeql/go-queries@1.6.12:codeql-suites/go-security-extended.qls --additional-packs=$sandbox/tmp/" "the analysis"
  # Every call after the version check keeps its caches in the run's folder.
  assert_equal 3 "$(grep -c -- "--common-caches=$sandbox/tmp/xquik-cli-codeql." "$sandbox_logs/codeql")" "the calls with a private cache"
  assert_equal "" "$(ls "$sandbox/tmp")" "the leftover temporary files"
}

test_a_codeql_result_fails_the_run_and_is_printed() {
  prepare_codeql
  run_env=("FAKE_RESULTS=1")

  run_in_repo scripts/codeql

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "go/sql-injection at pkg/lib.go:7
CodeQL found 1 results. Fix each one above, then rerun." "$run_stderr" "stderr"
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
