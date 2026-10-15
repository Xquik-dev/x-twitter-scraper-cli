#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

node_path="$(command -v node)" || {
  echo "node is not on PATH. Install Node.js first." >&2
  exit 1
}
readonly node_path

install_launcher() {
  printf '#!/bin/sh\nexec "%s" "$@"\n' "$node_path" | fake node
  mkdir -p "$sandbox_repo/node_modules/x-twitter-scraper-cli/bin"
  cp "$tests_root/npm/x-twitter-scraper-cli/bin/x-twitter-scraper.js" \
    "$sandbox_repo/node_modules/x-twitter-scraper-cli/bin/x-twitter-scraper.js"
  target="$("$node_path" -p 'process.platform + "-" + process.arch')"
}

install_platform_binary() {
  local folder="$sandbox_repo/node_modules/x-twitter-scraper-cli-$target/bin"
  mkdir -p "$folder"
  cat >"$folder/x-twitter-scraper"
  chmod 755 "$folder/x-twitter-scraper"
}

test_the_launcher_runs_the_platform_binary_with_every_argument_and_its_exit_code() {
  install_launcher
  install_platform_binary <<'EOF'
#!/bin/sh
printf '%s|' "$@"
echo
echo "to stderr" >&2
exit 3
EOF

  run_in_repo node node_modules/x-twitter-scraper-cli/bin/x-twitter-scraper.js \
    x:tweets search --q "from:nasa moon" --limit 10

  assert_equal 3 "$run_status" "the exit code"
  assert_equal "x:tweets|search|--q|from:nasa moon|--limit|10|" "$run_stdout" "stdout"
  assert_equal "to stderr" "$run_stderr" "stderr"
}

test_the_launcher_exits_as_a_shell_does_when_the_binary_dies_from_sigpipe() {
  install_launcher
  install_platform_binary <<'EOF'
#!/bin/sh
kill -PIPE $$
EOF

  run_in_repo node node_modules/x-twitter-scraper-cli/bin/x-twitter-scraper.js x:tweets search

  assert_equal 141 "$run_status" "the exit code"
}

test_the_launcher_explains_a_missing_platform_build() {
  install_launcher

  run_in_repo node node_modules/x-twitter-scraper-cli/bin/x-twitter-scraper.js --help

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "" "$run_stdout" "stdout"
  assert_equal "x-twitter-scraper-cli found no binary for $target.
Reinstall it without --omit=optional or --no-optional.
On another platform, download a build from https://github.com/Xquik-dev/x-twitter-scraper-cli/releases." \
    "$run_stderr" "stderr"
}

run_tests
