#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

# Runs scripts/test-on-terminal and scripts/on-terminal in a throwaway
# repository. A fake go lists 2 packages and writes a fake test binary for
# each, so no test compiles Go code. The fake test binary misbehaves only
# when its standard output is a terminal, as the real tests once did.

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

readonly advice="Give the code its colors, input and output as values, or call terminaltest.Run from TestMain."

# FAKE_ON_TERMINAL picks what the tests of the second package do on a
# terminal. They pass when it is unset.
prepare_packages() {
  copy_file scripts/on-terminal
  copy_file scripts/test-on-terminal
  fake go <<'EOF_GO'
#!/bin/sh
# Stands in for go. "go list" prints 2 packages. "go test -c -o <binary>
# <package>" writes a test binary that records where it runs.
case "$1" in
  list) printf 'example.com/cli/first\nexample.com/cli/second\n' ;;
  test)
    binary="$4"
    package="$5"
    cat >"$binary" <<EOF_BINARY
#!/bin/sh
if [ -t 1 ]; then output=terminal; else output=pipe; fi
echo "$package \$output TERM=\$TERM NO_COLOR=\${NO_COLOR:-unset} \$*" >>"\$FAKE_LOGS/tests"
if [ "$package" = example.com/cli/second ] && [ "\$output" = terminal ]; then
  case "\${FAKE_ON_TERMINAL:-}" in
    fail)
      echo "--- FAIL: TestColors"
      exit 1
      ;;
    colors) printf '\033[31mred\033[0m\n' ;;
  esac
fi
echo PASS
EOF_BINARY
    chmod 755 "$binary"
    ;;
esac
EOF_GO
}

test_every_package_runs_on_a_color_terminal_with_a_2_minute_limit() {
  prepare_packages
  run_env=("NO_COLOR=1")

  run_in_repo scripts/test-on-terminal

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_equal "Every package passes on a terminal." "$run_stdout" "stdout"
  assert_equal "example.com/cli/first terminal TERM=screen-256color NO_COLOR=unset -test.count=1 -test.timeout=2m
example.com/cli/second terminal TERM=screen-256color NO_COLOR=unset -test.count=1 -test.timeout=2m" \
    "$(log_of tests)" "the test runs"
}

test_a_test_that_fails_only_on_a_terminal_fails_the_check_and_shows_its_output() {
  prepare_packages
  run_env=("FAKE_ON_TERMINAL=fail")

  run_in_repo scripts/test-on-terminal

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stderr" "--- FAIL: TestColors" "stderr"
  assert_ends_with "$run_stderr" "The tests of example.com/cli/second fail on a terminal. $advice" "stderr"
  assert_equal "" "$run_stdout" "stdout"
}

test_a_test_that_writes_control_codes_on_a_terminal_fails_the_check() {
  prepare_packages
  run_env=("FAKE_ON_TERMINAL=colors")

  run_in_repo scripts/test-on-terminal

  assert_equal 1 "$run_status" "the exit code"
  assert_contains "$run_stderr" "^[[31mred^[[0m" "stderr"
  assert_ends_with "$run_stderr" "The tests of example.com/cli/second write terminal control codes. $advice" "stderr"
  assert_equal "" "$run_stdout" "stdout"
}

test_on_terminal_gives_the_command_a_terminal_and_returns_its_exit_code() {
  copy_file scripts/on-terminal

  run_in_repo scripts/on-terminal sh -c 'if [ -t 1 ]; then echo "a terminal"; fi; exit 7'

  assert_equal 7 "$run_status" "the exit code"
  assert_contains "$run_stdout" "a terminal" "stdout"
}

test_on_terminal_without_a_command_stops_with_the_usage_line() {
  copy_file scripts/on-terminal

  run_in_repo scripts/on-terminal

  assert_equal 2 "$run_status" "the exit code"
  assert_equal "Run scripts/on-terminal <command> [argument...]." "$run_stderr" "stderr"
}

run_tests
