#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

# Helpers for the tests that run this repository's scripts as a maintainer
# would. Each test runs a script in a throwaway folder. Its PATH holds fake
# commands in place of the tools that reach other systems, so no test can
# reach GitHub, the Go module proxy, the CodeQL registry or the Keychain.
#
# A test file sources this file, defines functions named test_*, and ends
# with run_tests.

tests_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
readonly tests_root

# Tools that only read and write local files. Each gets a wrapper in the
# throwaway bin folder, wherever this machine installed it.
readonly local_tools=(bash git jq)

# Tools that reach other systems or run long builds. Each call fails unless
# the test fakes it.
readonly remote_tools=(codeql curl gh go goreleaser security shellcheck uv uvx)

# Builds the throwaway folder of 1 test. It holds an empty repository folder,
# a bin folder for fake commands and a logs folder that the fakes write to.
create_sandbox() {
  local tool tool_path
  sandbox="$(mktemp -d "${TMPDIR:-/tmp}/xquik-cli-script-test.XXXXXX")"
  sandbox="$(cd -- "$sandbox" && pwd -P)"
  sandbox_repo="$sandbox/repo"
  sandbox_bin="$sandbox/bin"
  sandbox_logs="$sandbox/logs"
  mkdir -p "$sandbox_repo" "$sandbox_bin" "$sandbox_logs" "$sandbox/tmp" "$sandbox/home"
  for tool in "${local_tools[@]}"; do
    tool_path="$(command -v "$tool")" || {
      echo "$tool is not on PATH. Install it first." >&2
      return 1
    }
    printf '#!/bin/sh\nexec "%s" "$@"\n' "$tool_path" >"$sandbox_bin/$tool"
    chmod 755 "$sandbox_bin/$tool"
  done
  for tool in "${remote_tools[@]}"; do
    printf '#!/bin/sh\necho "%s was called, but this test did not fake it." >&2\nexit 97\n' \
      "$tool" >"$sandbox_bin/$tool"
    chmod 755 "$sandbox_bin/$tool"
  done
  run_env=()
}

remove_sandbox() {
  rm -rf -- "${sandbox:?}"
}

# Puts a fake command on PATH. Its source comes from standard input.
fake() {
  cat >"$sandbox_bin/$1"
  chmod 755 "$sandbox_bin/$1"
}

# Removes a command from the throwaway PATH.
remove_tool() {
  rm -f -- "$sandbox_bin/$1"
}

# Writes an executable file into the throwaway repository from standard input.
stub() {
  mkdir -p "$(dirname -- "$sandbox_repo/$1")"
  cat >"$sandbox_repo/$1"
  chmod 755 "$sandbox_repo/$1"
}

# Writes a plain file into the throwaway repository from standard input.
write_file() {
  mkdir -p "$(dirname -- "$sandbox_repo/$1")"
  cat >"$sandbox_repo/$1"
}

# Copies a file of this repository into the throwaway repository.
copy_file() {
  mkdir -p "$(dirname -- "$sandbox_repo/$1")"
  cp -p "$tests_root/$1" "$sandbox_repo/$1"
}

sandbox_git() {
  env -i \
    GIT_AUTHOR_EMAIL="tester@example.com" \
    GIT_AUTHOR_NAME="Script Tester" \
    GIT_COMMITTER_EMAIL="tester@example.com" \
    GIT_COMMITTER_NAME="Script Tester" \
    GIT_CONFIG_GLOBAL=/dev/null \
    GIT_CONFIG_NOSYSTEM=1 \
    HOME="$sandbox/home" \
    PATH="$sandbox_bin:/usr/bin:/bin" \
    git -C "$sandbox_repo" "$@"
}

# Commits every file of the throwaway repository on main.
commit_repo() {
  sandbox_git init --quiet --initial-branch=main
  sandbox_git add --all
  sandbox_git commit --quiet --message "Test commit"
}

# Gives the throwaway repository a local bare origin that holds main.
add_origin() {
  sandbox_git init --quiet --bare "$sandbox/origin.git"
  sandbox_git remote add origin "$sandbox/origin.git"
  sandbox_git push --quiet origin main
}

# Runs a command in the throwaway repository with a clean environment. The
# run_env array adds variables. The results land in run_status, run_stdout
# and run_stderr.
run_in_repo() {
  run_status=0
  (
    cd -- "$sandbox_repo"
    env -i \
      FAKE_LOGS="$sandbox_logs" \
      GIT_AUTHOR_EMAIL="tester@example.com" \
      GIT_AUTHOR_NAME="Script Tester" \
      GIT_COMMITTER_EMAIL="tester@example.com" \
      GIT_COMMITTER_NAME="Script Tester" \
      GIT_CONFIG_GLOBAL=/dev/null \
      GIT_CONFIG_NOSYSTEM=1 \
      HOME="$sandbox/home" \
      PATH="$sandbox_bin:/usr/bin:/bin" \
      TMPDIR="$sandbox/tmp" \
      ${run_env[@]+"${run_env[@]}"} \
      "$@"
  ) >"$sandbox/stdout" 2>"$sandbox/stderr" </dev/null || run_status=$?
  run_stdout="$(cat "$sandbox/stdout")"
  run_stderr="$(cat "$sandbox/stderr")"
}

# Prints the lines a fake or a stub wrote to its log, or nothing when it
# never ran.
log_of() {
  if [[ -f "$sandbox_logs/$1" ]]; then
    cat "$sandbox_logs/$1"
  fi
}

# Ends the running test as failed, and shows what the last run printed. Each
# test runs in its own subshell.
fail_test() {
  echo "  $1" >&2
  if [[ -n "${run_status:-}" ]]; then
    echo "  The last run exited with $run_status. Its stdout and stderr follow." >&2
    echo "${run_stdout:-}" >&2
    echo "${run_stderr:-}" >&2
  fi
  exit 1
}

assert_equal() {
  if [[ "$1" != "$2" ]]; then
    fail_test "$3 differs.
  expected
$1
  actual
$2"
  fi
}

assert_contains() {
  if [[ "$1" != *"$2"* ]]; then
    fail_test "$3 lacks the text below.
  wanted
$2
  actual
$1"
  fi
}

assert_not_contains() {
  if [[ "$1" == *"$2"* ]]; then
    fail_test "$3 holds text it must not hold.
  unwanted
$2"
  fi
}

assert_ends_with() {
  if [[ "$1" != *"$2" ]]; then
    fail_test "$3 does not end with the text below.
  wanted
$2
  actual
$1"
  fi
}

# Runs every test_* function of the calling file, each in its own throwaway
# folder, and fails when any test fails.
run_tests() {
  local name status failed=0 total=0
  while IFS= read -r name; do
    total=$((total + 1))
    # The subshell is a plain command. Inside a condition, bash would ignore
    # set -e and a failing step could not fail its test.
    set +e
    (
      set -euo pipefail
      create_sandbox
      trap remove_sandbox EXIT
      "$name"
    )
    status=$?
    set -e
    if ((status == 0)); then
      echo "ok $name"
    else
      echo "not ok $name"
      failed=$((failed + 1))
    fi
  done < <(declare -F | sed -n 's/^declare -f \(test_.*\)$/\1/p')
  if ((total == 0)); then
    echo "$0 defines no tests." >&2
    return 1
  fi
  if ((failed > 0)); then
    echo "$failed of $total tests failed in $(basename -- "$0")." >&2
    return 1
  fi
  echo "$total tests passed in $(basename -- "$0")."
}
