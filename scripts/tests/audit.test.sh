#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

# Runs scripts/audit in a throwaway repository. A fake go command answers go
# mod verify and go list, and serves a canned govulncheck report, so no test
# reaches the Go module proxy or the vulnerability database.

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

readonly cli_module="github.com/Xquik-dev/x-twitter-scraper-cli"
readonly shipped_module="github.com/urfave/cli/v3"
readonly test_module="github.com/stretchr/testify"

day() {
  jq -nr --argjson offset "$1" 'now + $offset * 86400 | strftime("%Y-%m-%d")'
}

# Builds a repository with scripts/audit, an exceptions file from standard
# input and an empty govulncheck report.
prepare_audit() {
  copy_file scripts/audit
  write_file audit-exceptions.json
  : >"$sandbox/report.json"
  fake go <<EOF
#!/bin/sh
# Stands in for go. It knows 3 packages and the modules that ship, and serves
# the canned govulncheck report.
case "\$1 \$2" in
  "mod verify")
    echo "all modules verified"
    ;;
  "list -f")
    echo "$cli_module/cmd/x-twitter-scraper main"
    echo "$cli_module/pkg/cmd cmd"
    echo "$cli_module/internal/mocktest mocktest"
    ;;
  "list -deps")
    echo "\$*" >>"\$FAKE_LOGS/go-deps"
    echo "$cli_module"
    echo ""
    echo "$shipped_module"
    ;;
  "run golang.org/x/vuln/cmd/govulncheck@v1.9.0")
    echo "\$*" >>"\$FAKE_LOGS/govulncheck"
    [ -z "\${FAKE_SCANNER_FAILS:-}" ] || exit 1
    cat "$sandbox/report.json"
    ;;
  *)
    echo "unexpected go \$*" >&2
    exit 97
    ;;
esac
EOF
}

# Adds a finding to the report. Arguments are the advisory ID, the module, the
# fixed version or an empty string, and an optional called function. Without
# a function, govulncheck traced no call.
add_finding() {
  local trace
  if [[ -n "${4:-}" ]]; then
    trace="{\"module\": \"$2\", \"package\": \"$2/pkg\", \"function\": \"$4\"}"
  else
    trace="{\"module\": \"$2\"}"
  fi
  {
    echo "{\"osv\": {\"id\": \"$1\", \"database_specific\": {\"severity\": \"LOW\"}}}"
    echo "{\"finding\": {\"osv\": \"$1\", \"fixed_version\": \"$3\", \"trace\": [$trace]}}"
  } >>"$sandbox/report.json"
}

exception() {
  printf '{"id": "%s", "module": "%s", "reason": "%s", "recheck": "%s"}' "$1" "$2" "$3" "$4"
}

test_a_finding_without_a_traced_call_in_a_shipped_module_blocks() {
  prepare_audit <<<"[]"
  add_finding GO-2026-1000 "$shipped_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1000 is in $shipped_module, which ships in the CLI. It has no fix yet, so replace the module." "stdout"
  assert_contains "$run_stderr" "The dependency audit blocked 1 findings." "stderr"
}

test_a_low_finding_in_a_shipped_module_blocks_even_when_listed_and_unfixable() {
  prepare_audit <<<"[$(exception GO-2026-1001 "$shipped_module" "build-time only" "$(day 10)")]"
  add_finding GO-2026-1001 "$shipped_module" "" "Parse"

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1001 is in $shipped_module, which ships in the CLI. It has no fix yet, so replace the module. Remove its audit-exceptions.json entry." "stdout"
}

test_a_standard_library_finding_blocks() {
  prepare_audit <<<"[]"
  add_finding GO-2026-1002 stdlib "v1.27.3" "net/http.Serve"

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1002 is in stdlib, which ships in the CLI. Upgrade to v1.27.3." "stdout"
}

test_a_test_only_finding_with_a_fix_blocks_even_when_listed() {
  prepare_audit <<<"[$(exception GO-2026-1003 "$test_module" "build-time only" "$(day 10)")]"
  add_finding GO-2026-1003 "$test_module" "v1.13.0"

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1003 in $test_module is fixed in v1.13.0. Upgrade it." "stdout"
}

test_an_unfixable_test_only_finding_blocks_when_no_exception_lists_it() {
  prepare_audit <<<"[]"
  add_finding GO-2026-1004 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1004 in $test_module has no fix. Review it and list it in audit-exceptions.json." "stdout"
}

test_an_exception_for_another_module_excuses_nothing() {
  prepare_audit <<<"[$(exception GO-2026-1005 "example.com/other" "build-time only" "$(day 10)")]"
  add_finding GO-2026-1005 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. GO-2026-1005 in $test_module has no fix." "stdout"
}

test_a_listed_unfixable_test_only_finding_passes_until_its_recheck_date() {
  prepare_audit <<<"[$(exception GO-2026-1006 "$test_module" "build-time only" "$(day 10)")]"
  add_finding GO-2026-1006 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Accepted the reviewed exception for GO-2026-1006 in $test_module." "stdout"
  assert_ends_with "$run_stdout" "Dependency audit passed with 1 reviewed exception." "stdout"
}

test_a_listed_finding_blocks_again_on_its_recheck_date() {
  local today
  today="$(day 0)"
  prepare_audit <<<"[$(exception GO-2026-1007 "$test_module" "build-time only" "$today")]"
  add_finding GO-2026-1007 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. The exception for GO-2026-1007 in $test_module expired on $today. Recheck the finding." "stdout"
}

test_an_exception_with_a_recheck_date_more_than_30_days_out_stops_the_audit() {
  prepare_audit <<<"[$(exception GO-2026-1008 "$test_module" "build-time only" "$(day 45)")]"
  add_finding GO-2026-1008 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. The exception for GO-2026-1008 in $test_module rechecks on $(day 45). Use $(day 30) or earlier." "stdout"
}

test_an_exception_with_an_unknown_reason_stops_the_audit() {
  prepare_audit <<<"[$(exception GO-2026-1009 "$test_module" "low risk" "$(day 10)")]"
  add_finding GO-2026-1009 "$test_module" ""

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. The exception for GO-2026-1009 in $test_module needs a reason of \"build-time only\" or \"peer dependency provided by the host\"." "stdout"
}

test_a_report_without_findings_passes_with_an_empty_exceptions_list() {
  prepare_audit <<<"[]"

  run_in_repo scripts/audit

  assert_equal 0 "$run_status" "the exit code"
  assert_equal "" "$run_stderr" "stderr"
  assert_ends_with "$run_stdout" "Dependency audit passed." "stdout"
  assert_equal "run golang.org/x/vuln/cmd/govulncheck@v1.9.0 -format json -test ./..." "$(log_of govulncheck)" "the govulncheck call"
}

test_the_modules_that_ship_come_from_the_binary_and_importable_packages() {
  prepare_audit <<<"[]"

  run_in_repo scripts/audit

  assert_equal 0 "$run_status" "the exit code"
  assert_equal "list -deps -f {{with .Module}}{{.Path}}{{end}} $cli_module/cmd/x-twitter-scraper $cli_module/pkg/cmd" \
    "$(log_of go-deps)" "the go list call"
}

test_a_scanner_error_blocks() {
  prepare_audit <<<"[]"
  run_env=("FAKE_SCANNER_FAILS=1")

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_ends_with "$run_stderr" "govulncheck did not finish. Check the network and the Go toolchain, then rerun." "stderr"
}

test_an_exceptions_file_that_is_not_a_list_stops_the_audit() {
  prepare_audit <<<'{"id": "GO-2026-1010"}'

  run_in_repo scripts/audit

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stdout" "Blocked. audit-exceptions.json must hold an array of entries with id, module, reason and recheck." "stdout"
}

run_tests
