#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

# Runs bin/check-local in a throwaway repository. Stubs stand in for the check
# scripts. Fake commands stand in for the GitHub CLI, uv, GoReleaser, Go, the
# CodeQL CLI and shellcheck, so no test reaches GitHub or runs a real check.

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

readonly check_order="bootstrap
uvx
lint
shellcheck
test-scripts
test
audit
goreleaser check
goreleaser release --snapshot --clean
goreleaser release --snapshot --clean
check-release-archives
check-release-debug-info
go test fuzz
codeql"
readonly passed_line="Passed modules, REUSE, lint, shellcheck, script tests, Go tests, audit, reproducible build, fuzzing and CodeQL."
readonly status_endpoint="repos/Xquik-dev/x-twitter-scraper-cli/statuses"

# Builds a repository whose checks all pass, on a machine with every tool.
prepare_checks() {
  local name
  copy_file bin/check-local
  for name in bootstrap lint test-scripts test codeql; do
    stub "scripts/$name" <<EOF
#!/usr/bin/env bash
# Stands in for scripts/$name. It fails when FAKE_FAILING_CHECK names it.
echo "$name" >>"\$FAKE_LOGS/order"
[[ "\${FAKE_FAILING_CHECK:-}" != "$name" ]]
EOF
  done
  stub scripts/audit <<'EOF'
#!/usr/bin/env bash
# Stands in for scripts/audit. FAKE_EXCEPTIONS sets the accepted exceptions.
echo "audit" >>"$FAKE_LOGS/order"
[[ "${FAKE_FAILING_CHECK:-}" != "audit" ]]
case "${FAKE_EXCEPTIONS:-0}" in
  0) echo "Dependency audit passed." ;;
  1) echo "Dependency audit passed with 1 reviewed exception." ;;
  *) echo "Dependency audit passed with $FAKE_EXCEPTIONS reviewed exceptions." ;;
esac
EOF
  stub bin/check-release-archives <<'EOF'
#!/usr/bin/env bash
# Stands in for bin/check-release-archives.
echo "check-release-archives" >>"$FAKE_LOGS/order"
EOF
  stub bin/check-release-debug-info <<'EOF'
#!/usr/bin/env bash
# Stands in for bin/check-release-debug-info.
echo "check-release-debug-info" >>"$FAKE_LOGS/order"
EOF
  stub scripts/notes.txt <<'EOF'
Not a script.
EOF
  write_file release-please-config.json <<'EOF'
{"draft": true}
EOF
  write_file .goreleaser.yml <<'EOF'
release:
  draft: true
EOF
  write_file .gitignore <<'EOF'
dist/
EOF
  fake goreleaser <<'EOF'
#!/bin/sh
# Stands in for GoReleaser. A snapshot build writes a checksums file. It
# differs between builds when FAKE_UNREPRODUCIBLE is set, and a build adds a
# tracked-looking file when FAKE_DIRTY_BUILD is set.
echo "goreleaser $*" >>"$FAKE_LOGS/order"
case "$1" in
  check) exit 0 ;;
  release)
    mkdir -p dist
    count=$(grep -c 'release' "$FAKE_LOGS/order")
    if [ -n "${FAKE_UNREPRODUCIBLE:-}" ]; then
      echo "build $count" >dist/x_checksums.txt
    else
      echo "same build" >dist/x_checksums.txt
    fi
    if [ -n "${FAKE_DIRTY_BUILD:-}" ]; then
      echo "left behind" >leftover.txt
    fi
    ;;
esac
EOF
  fake go <<'EOF'
#!/bin/sh
# Stands in for go. It records the fuzzing call.
echo "go test fuzz" >>"$FAKE_LOGS/order"
echo "$*" >>"$FAKE_LOGS/go"
[ "${FAKE_FAILING_CHECK:-}" != "fuzz" ]
EOF
  fake shellcheck <<'EOF'
#!/bin/sh
# Stands in for shellcheck. It records the scripts it is asked to check.
echo "shellcheck" >>"$FAKE_LOGS/order"
echo "$*" >>"$FAKE_LOGS/shellcheck"
EOF
  fake uvx <<'EOF'
#!/bin/sh
# Stands in for uvx. It records its arguments.
echo "uvx" >>"$FAKE_LOGS/order"
echo "$*" >>"$FAKE_LOGS/uvx"
EOF
  fake codeql <<'EOF'
#!/bin/sh
# Stands in for the CodeQL CLI. bin/check-local only checks that it exists.
exit 0
EOF
  fake gh <<'EOF'
#!/bin/sh
# Stands in for the GitHub CLI. It knows every commit unless
# FAKE_GH_UNKNOWN_COMMIT is set, and records each status it is asked to post.
case "$*" in
  "api repos/Xquik-dev/x-twitter-scraper-cli/commits/"*)
    [ -z "${FAKE_GH_UNKNOWN_COMMIT:-}" ]
    exit $?
    ;;
esac
state=""
description=""
endpoint=""
for argument in "$@"; do
  case "$argument" in
    repos/*) endpoint="$argument" ;;
    state=*) state="${argument#state=}" ;;
    description=*) description="${argument#description=}" ;;
  esac
done
echo "$endpoint $state $description" >>"$FAKE_LOGS/gh"
EOF
  commit_repo
}

test_a_passing_run_lists_the_checks_in_order_and_runs_reuse_through_uvx() {
  prepare_checks

  run_in_repo bin/check-local --no-status

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "$passed_line" "stdout"
  assert_contains "$run_stdout" "Commit status not posted because --no-status was passed." "stdout"
  assert_equal "$check_order" "$(log_of order)" "the order of the checks"
  assert_equal "--from reuse==6.2.0 reuse lint" "$(log_of uvx)" "the uvx call"
  assert_equal "" "$(log_of gh)" "the posted statuses"
}

test_a_run_from_a_terminal_passes_the_same_checks_in_the_same_order() {
  prepare_checks
  copy_file scripts/on-terminal

  run_in_repo scripts/on-terminal bin/check-local --no-status

  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "${run_stdout//$'\r'/}" "$passed_line" "the terminal output"
  assert_equal "$check_order" "$(log_of order)" "the order of the checks"
}

test_fuzzing_runs_for_15_seconds_unless_fuzz_time_says_otherwise() {
  prepare_checks

  run_in_repo bin/check-local --no-status
  run_env=("FUZZ_TIME=2m")
  run_in_repo bin/check-local --no-status

  assert_equal 0 "$run_status" "the exit code"
  assert_equal "test ./internal/apiform -run ^\$ -fuzz ^FuzzMarshalRoundTrip\$ -fuzztime=15s
test ./internal/apiform -run ^\$ -fuzz ^FuzzMarshalRoundTrip\$ -fuzztime=2m" "$(log_of go)" "the fuzzing calls"
}

test_shellcheck_reads_every_bash_script_and_skips_other_files() {
  prepare_checks
  stub scripts/tests/sample.test.sh <<'EOF'
#!/usr/bin/env bash
exit 0
EOF
  sandbox_git add --all
  sandbox_git commit --quiet --message "Add a test file"

  run_in_repo bin/check-local --no-status

  assert_equal 0 "$run_status" "the exit code"
  assert_equal "--external-sources bin/check-local bin/check-release-archives bin/check-release-debug-info scripts/audit scripts/bootstrap scripts/codeql scripts/lint scripts/test scripts/test-scripts scripts/tests/sample.test.sh" \
    "$(log_of shellcheck)" "the shellcheck call"
}

test_a_clean_pushed_head_gets_a_pending_status_then_a_success_status() {
  local head
  prepare_checks
  head="$(sandbox_git rev-parse HEAD)"

  run_in_repo bin/check-local

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_equal "$status_endpoint/$head pending Local checks are running.
$status_endpoint/$head success $passed_line" "$(log_of gh)" "the posted statuses"
}

test_the_success_status_names_the_reviewed_audit_exceptions() {
  local head
  prepare_checks
  head="$(sandbox_git rev-parse HEAD)"
  run_env=("FAKE_EXCEPTIONS=2")

  run_in_repo bin/check-local

  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$(log_of gh)" "$status_endpoint/$head success Passed modules, REUSE, lint, shellcheck, script tests, Go tests, audit with 2 reviewed exceptions, reproducible build, fuzzing and CodeQL." "the success status"
}

test_a_working_tree_with_changes_gets_no_status() {
  prepare_checks
  write_file draft.txt <<'EOF'
Uncommitted work.
EOF

  run_in_repo bin/check-local

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Commit status not posted because the working tree has changes." "stdout"
  assert_equal "" "$(log_of gh)" "the posted statuses"
}

test_a_commit_that_github_lacks_gets_no_status() {
  prepare_checks
  run_env=("FAKE_GH_UNKNOWN_COMMIT=1")

  run_in_repo bin/check-local

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Commit status not posted because GitHub has no commit" "stdout"
  assert_equal "" "$(log_of gh)" "the posted statuses"
}

test_a_run_without_uv_stops_with_the_install_command() {
  prepare_checks
  remove_tool uvx
  fake reuse <<'EOF'
#!/bin/sh
echo "reuse" >>"$FAKE_LOGS/order"
EOF

  run_in_repo bin/check-local --no-status

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Check failed at Required tools. Install uv with brew install uv, then rerun." "$run_stderr" "stderr"
  assert_equal "" "$(log_of order)" "the checks that ran"
}

test_a_run_without_the_codeql_cli_stops_with_the_install_command() {
  local head
  prepare_checks
  remove_tool codeql
  head="$(sandbox_git rev-parse HEAD)"

  run_in_repo bin/check-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Check failed at Required tools. Install the CodeQL CLI with brew install --cask codeql, then rerun." "$run_stderr" "stderr"
  assert_equal "$status_endpoint/$head pending Local checks are running.
$status_endpoint/$head failure Failed at Required tools." "$(log_of gh)" "the posted statuses"
}

test_a_failing_check_stops_the_run_and_posts_a_failure_status_that_names_it() {
  local head
  prepare_checks
  head="$(sandbox_git rev-parse HEAD)"
  run_env=("FAKE_FAILING_CHECK=test")

  run_in_repo bin/check-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Check failed at Go tests and coverage. Fix the failing test or add coverage, then rerun." "$run_stderr" "stderr"
  assert_equal "bootstrap
uvx
lint
shellcheck
test-scripts
test" "$(log_of order)" "the checks that ran"
  assert_equal "$status_endpoint/$head pending Local checks are running.
$status_endpoint/$head failure Failed at Go tests and coverage." "$(log_of gh)" "the posted statuses"
}

test_a_codeql_failure_posts_a_failure_status_and_no_success() {
  local head
  prepare_checks
  head="$(sandbox_git rev-parse HEAD)"
  run_env=("FAKE_FAILING_CHECK=codeql")

  run_in_repo bin/check-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Check failed at CodeQL security-extended queries. Fix the CodeQL error or results above, then rerun." "$run_stderr" "stderr"
  assert_equal "$status_endpoint/$head pending Local checks are running.
$status_endpoint/$head failure Failed at CodeQL security-extended queries." "$(log_of gh)" "the posted statuses"
}

test_snapshot_builds_with_different_checksums_fail_the_run() {
  prepare_checks
  run_env=("FAKE_UNREPRODUCIBLE=1")

  run_in_repo bin/check-local --no-status

  assert_equal 1 "$run_status" "the exit code"
  assert_ends_with "$run_stderr" "Check failed at Reproducible snapshot build. Remove the build input that changes between builds, such as a timestamp." "stderr"
}

test_a_check_that_changes_the_tree_fails_the_run() {
  local head
  prepare_checks
  head="$(sandbox_git rev-parse HEAD)"
  run_env=("FAKE_DIRTY_BUILD=1")

  run_in_repo bin/check-local

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "?? leftover.txt" "stdout"
  assert_equal "Check failed at Unchanged working tree. A check changed or added the files above. Fix that check so it leaves the tree alone, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of gh)" "success" "the posted statuses"
}

test_an_unknown_option_stops_with_the_usage_line() {
  prepare_checks

  run_in_repo bin/check-local --fast

  assert_equal 2 "$run_status" "the exit code"
  assert_equal "Unknown option --fast. Run bin/check-local or bin/check-local --no-status." "$run_stderr" "stderr"
  assert_equal "" "$(log_of order)" "the checks that ran"
}

run_tests
