#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

command -v zip >/dev/null 2>&1 || {
  echo "zip is not on PATH. Install it first." >&2
  exit 1
}

readonly test_token="npm_test_token_for_release_npm"
readonly platforms=(
  "darwin x64 macos_amd64.zip"
  "darwin arm64 macos_arm64.zip"
  "linux ia32 linux_386.tar.gz"
  "linux x64 linux_amd64.tar.gz"
  "linux arm linux_armv6.tar.gz"
  "linux arm64 linux_arm64.tar.gz"
  "win32 ia32 windows_386.zip"
  "win32 x64 windows_amd64.zip"
  "win32 arm64 windows_arm64.zip"
)

build_release() {
  local version="$1" os cpu suffix folder binary
  mkdir -p "$sandbox/release"
  for platform in "${platforms[@]}"; do
    read -r os cpu suffix <<<"$platform"
    folder="$sandbox/archive/$os-$cpu"
    binary=x-twitter-scraper
    if [[ "$os" == win32 ]]; then
      binary=x-twitter-scraper.exe
    fi
    mkdir -p "$folder/LICENSES" "$folder/completions"
    printf '#!/bin/sh\necho "%s %s build $*"\n' "$os" "$cpu" >"$folder/$binary"
    chmod 755 "$folder/$binary"
    echo "Apache License" >"$folder/LICENSE"
    echo "BSD 3-Clause" >"$folder/LICENSES/BSD-3-Clause.txt"
    echo "Notices" >"$folder/THIRD_PARTY_NOTICES.md"
    echo "complete" >"$folder/completions/x-twitter-scraper.bash"
    if [[ "$suffix" == *.zip ]]; then
      (cd "$folder" && zip -q -r "$sandbox/release/x-twitter-scraper_${version}_${suffix}" .)
    else
      tar -czf "$sandbox/release/x-twitter-scraper_${version}_${suffix}" -C "$folder" .
    fi
  done
  (cd "$sandbox/release" && shasum -a 256 -- * >"x-twitter-scraper_${version}_checksums.txt")
}

prepare_npm() {
  local version="${1:-1.2.3}" tag_before_npm="${2:-}"
  copy_file bin/release-npm
  copy_file LICENSE
  write_file .release-please-manifest.json <<EOF
{".": "$version"}
EOF
  if [[ -n "$tag_before_npm" ]]; then
    commit_repo
    sandbox_git tag --annotate "v$version" --message "Release v$version"
  fi
  copy_file npm/x-twitter-scraper-cli/package.json
  copy_file npm/x-twitter-scraper-cli/README.md
  copy_file npm/x-twitter-scraper-cli/bin/x-twitter-scraper.js
  build_release "$version"
  fake gh <<'EOF'
#!/bin/sh
echo "gh $*" >>"$FAKE_LOGS/order"
case "$1 $2" in
  "api repos/Xquik-dev/x-twitter-scraper-cli/releases/tags/"*)
    if [ -n "${FAKE_NO_RELEASE:-}" ]; then
      echo "gh: Not Found (HTTP 404)" >&2
      exit 1
    fi
    printf '{"draft": %s}\n' "${FAKE_DRAFT:-false}"
    ;;
  "release download")
    dir=""
    previous=""
    for argument in "$@"; do
      [ "$previous" = "--dir" ] && dir="$argument"
      previous="$argument"
    done
    mkdir -p "$dir"
    cp "$FAKE_RELEASE"/* "$dir/"
    if [ -n "${FAKE_TAMPER:-}" ]; then
      echo "changed" >>"$dir/$FAKE_TAMPER"
    fi
    ;;
  *)
    echo "unexpected gh $*" >&2
    exit 97
    ;;
esac
EOF
  fake bun <<'EOF'
#!/bin/sh
case "$1 $2" in
  "pm pack")
    filename=""
    previous=""
    for argument in "$@"; do
      [ "$previous" = "--filename" ] && filename="$argument"
      previous="$argument"
    done
    staging="$(mktemp -d)"
    cp -Rp . "$staging/package"
    [ -z "${FAKE_NO_EXEC:-}" ] || chmod a-x "$staging"/package/bin/*
    tar -czf "$filename" -C "$staging" package
    rm -rf "$staging"
    mkdir -p "$FAKE_LOGS/tarballs"
    cp "$filename" "$FAKE_LOGS/tarballs/"
    ;;
  "publish "*)
    name="$(tar -xzOf "$2" package/package.json | jq -r .name)"
    echo "publish $name ${*#publish }" >>"$FAKE_LOGS/publish-arguments"
    if [ "${NPM_CONFIG_TOKEN:-}" = "npm_test_token_for_release_npm" ]; then
      echo "$name with the expected token" >>"$FAKE_LOGS/published"
    else
      echo "$name with no token or a wrong token" >>"$FAKE_LOGS/published"
    fi
    ;;
  *)
    echo "unexpected bun $*" >&2
    exit 97
    ;;
esac
EOF
  fake curl <<'EOF'
#!/bin/sh
for argument in "$@"; do url="$argument"; done
if [ -n "${FAKE_NPM_STATUS:-}" ]; then
  printf '%s' "$FAKE_NPM_STATUS"
  exit 0
fi
name="$(echo "$url" | sed 's|^https://registry.npmjs.org/||; s|/[^/]*$||')"
if grep -q "^$name with" "$FAKE_LOGS/published" 2>/dev/null; then
  printf 200
  exit 0
fi
for held in ${FAKE_NPM_HAS:-}; do
  if [ "$held" = "$name" ]; then
    printf 200
    exit 0
  fi
done
printf 404
EOF
  fake security <<'EOF'
#!/bin/sh
echo "security $*" >>"$FAKE_LOGS/order"
[ -z "${FAKE_NO_ITEM:-}" ] || exit 44
case " $* " in
  *" -w "*) echo "npm_test_token_for_release_npm" ;;
  *) echo "keychain: item attributes" ;;
esac
EOF
  if [[ -n "$tag_before_npm" ]]; then
    sandbox_git add --all
    sandbox_git commit --quiet --message "Add the npm packages"
  else
    commit_repo
    sandbox_git tag --annotate "v$version" --message "Release v$version"
  fi
  add_origin
  sandbox_git push --quiet origin --tags
  run_env=("FAKE_RELEASE=$sandbox/release")
}

publish_env() {
  run_env+=("NPM_KEYCHAIN_SERVICE=test-npm-service" "NPM_KEYCHAIN_ACCOUNT=test-npm-account")
}

package_json() {
  tar -xzOf "$sandbox_logs/tarballs/$1-1.2.3.tgz" package/package.json
}

test_a_dry_run_packs_1_package_per_platform_and_publishes_nothing() {
  local unpacked
  prepare_npm

  run_in_repo bin/release-npm --dry-run

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "  x-twitter-scraper-cli-darwin-arm64@1.2.3  " "stdout"
  assert_contains "$run_stdout" "  x-twitter-scraper-cli@1.2.3  " "stdout"
  assert_contains "$run_stdout" "Package count  10" "stdout"
  assert_ends_with "$run_stdout" "Dry run finished. Nothing was published." "stdout"
  assert_equal '{
  "x-twitter-scraper-cli-darwin-x64": "1.2.3",
  "x-twitter-scraper-cli-darwin-arm64": "1.2.3",
  "x-twitter-scraper-cli-linux-ia32": "1.2.3",
  "x-twitter-scraper-cli-linux-x64": "1.2.3",
  "x-twitter-scraper-cli-linux-arm": "1.2.3",
  "x-twitter-scraper-cli-linux-arm64": "1.2.3",
  "x-twitter-scraper-cli-win32-ia32": "1.2.3",
  "x-twitter-scraper-cli-win32-x64": "1.2.3",
  "x-twitter-scraper-cli-win32-arm64": "1.2.3"
}' "$(package_json x-twitter-scraper-cli | jq '.optionalDependencies')" "the optional dependencies"
  assert_equal '{"name":"x-twitter-scraper-cli-linux-arm","version":"1.2.3","os":["linux"],"cpu":["arm"],"files":["bin","LICENSE","LICENSES","THIRD_PARTY_NOTICES.md"]}' \
    "$(package_json x-twitter-scraper-cli-linux-arm | jq -c '{name, version, os, cpu, files}')" "the linux arm manifest"
  unpacked="$sandbox/unpacked"
  mkdir -p "$unpacked"
  tar -xzf "$sandbox_logs/tarballs/x-twitter-scraper-cli-linux-arm-1.2.3.tgz" -C "$unpacked"
  assert_equal "linux arm build --version" "$("$unpacked/package/bin/x-twitter-scraper" --version)" "the packed binary"
  assert_equal "LICENSE
LICENSES/BSD-3-Clause.txt
README.md
THIRD_PARTY_NOTICES.md
bin/x-twitter-scraper
package.json" "$(cd "$unpacked/package" && find . -type f | sed 's|^\./||' | sort)" "the linux arm files"
  assert_equal "" "$(log_of publish-arguments)" "the publish calls"
  assert_not_contains "$(log_of order)" "security" "the steps"
}

test_a_windows_package_holds_the_exe() {
  prepare_npm

  run_in_repo bin/release-npm --dry-run

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$(tar -tzf "$sandbox_logs/tarballs/x-twitter-scraper-cli-win32-x64-1.2.3.tgz")" \
    "package/bin/x-twitter-scraper.exe" "the win32 x64 files"
}

test_an_archive_that_differs_from_the_checksums_stops_before_packing() {
  prepare_npm
  run_env+=("FAKE_TAMPER=x-twitter-scraper_1.2.3_linux_amd64.tar.gz")

  run_in_repo bin/release-npm --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "x-twitter-scraper_1.2.3_linux_amd64.tar.gz does not match x-twitter-scraper_1.2.3_checksums.txt. Rerun, and check release v1.2.3 on GitHub if it fails again." \
    "$run_stderr" "stderr"
  assert_equal "" "$(ls "$sandbox_logs/tarballs" 2>/dev/null)" "the packed tarballs"
}

test_a_draft_release_stops_the_run() {
  prepare_npm
  run_env+=("FAKE_DRAFT=true")

  run_in_repo bin/release-npm --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "GitHub release v1.2.3 is still a draft. Publish it with bin/release-local, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of order)" "gh release download" "the steps"
}

test_a_missing_release_stops_the_run() {
  prepare_npm
  run_env+=("FAKE_NO_RELEASE=1")

  run_in_repo bin/release-npm --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "GitHub shows no release v1.2.3. Publish it with bin/release-local, then rerun." "$run_stderr" "stderr"
}

test_a_release_publishes_the_platform_packages_first_with_the_token_in_the_environment() {
  prepare_npm
  publish_env

  run_in_repo bin/release-npm

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "Published x-twitter-scraper-cli@1.2.3 at https://www.npmjs.com/package/x-twitter-scraper-cli" "stdout"
  assert_equal "x-twitter-scraper-cli-darwin-x64 with the expected token
x-twitter-scraper-cli-darwin-arm64 with the expected token
x-twitter-scraper-cli-linux-ia32 with the expected token
x-twitter-scraper-cli-linux-x64 with the expected token
x-twitter-scraper-cli-linux-arm with the expected token
x-twitter-scraper-cli-linux-arm64 with the expected token
x-twitter-scraper-cli-win32-ia32 with the expected token
x-twitter-scraper-cli-win32-x64 with the expected token
x-twitter-scraper-cli-win32-arm64 with the expected token
x-twitter-scraper-cli with the expected token" "$(log_of published)" "the published packages"
  assert_contains "$(log_of publish-arguments)" "--access public --tag latest --registry https://registry.npmjs.org" "the publish arguments"
  assert_equal "security find-generic-password -s test-npm-service -a test-npm-account
security find-generic-password -s test-npm-service -a test-npm-account -w" "$(log_of order | grep '^security')" "the Keychain reads"
  assert_not_contains "$(log_of publish-arguments)" "$test_token" "the publish arguments"
  assert_not_contains "$run_stdout$run_stderr" "$test_token" "the output"
}

test_a_rerun_skips_the_versions_npm_already_has() {
  prepare_npm
  publish_env
  run_env+=("FAKE_NPM_HAS=x-twitter-scraper-cli-darwin-x64 x-twitter-scraper-cli-darwin-arm64")

  run_in_repo bin/release-npm

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "npm already has x-twitter-scraper-cli-darwin-x64@1.2.3." "stdout"
  assert_not_contains "$(log_of published)" "x-twitter-scraper-cli-darwin-arm64" "the published packages"
  assert_contains "$(log_of published)" "x-twitter-scraper-cli-linux-x64 with the expected token" "the published packages"
}

test_a_prerelease_publishes_with_the_next_tag() {
  prepare_npm 1.3.0-beta.1
  publish_env

  run_in_repo bin/release-npm

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$(log_of publish-arguments)" "--tag next" "the publish arguments"
  assert_not_contains "$(log_of publish-arguments)" "--tag latest" "the publish arguments"
}

test_a_release_without_the_keychain_names_stops_before_any_download() {
  prepare_npm

  run_in_repo bin/release-npm

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Set NPM_KEYCHAIN_SERVICE and NPM_KEYCHAIN_ACCOUNT to the Keychain item of the npm token, then rerun." "$run_stderr" "stderr"
  assert_equal "" "$(log_of order)" "the steps"
}

test_a_release_from_a_branch_other_than_main_stops_the_run() {
  prepare_npm
  publish_env
  sandbox_git checkout --quiet -b feature

  run_in_repo bin/release-npm

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The current branch is not main. Switch to main, then rerun." "$run_stderr" "stderr"
}

test_an_unknown_option_stops_with_the_usage_line() {
  prepare_npm

  run_in_repo bin/release-npm --publish

  assert_equal 2 "$run_status" "the exit code"
  assert_equal "Unknown option --publish. Run bin/release-npm, bin/release-npm --dry-run or bin/release-npm --check." "$run_stderr" "stderr"
}

test_a_rerun_after_a_full_publish_reads_no_token() {
  prepare_npm
  publish_env
  run_env+=("FAKE_NPM_HAS=$(printf 'x-twitter-scraper-cli-%s ' darwin-x64 darwin-arm64 linux-ia32 linux-x64 linux-arm linux-arm64 win32-ia32 win32-x64 win32-arm64)x-twitter-scraper-cli")

  run_in_repo bin/release-npm

  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "npm already has every package for 1.2.3." "stdout"
  assert_not_contains "$(log_of order)" " -w" "the steps"
  assert_equal "" "$(log_of published)" "the published packages"
}

test_an_unexpected_npm_answer_stops_before_publishing() {
  prepare_npm
  publish_env
  run_env+=("FAKE_NPM_STATUS=500")

  run_in_repo bin/release-npm

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "npm answered 500 for x-twitter-scraper-cli-darwin-x64@1.2.3. Retry later." "$run_stderr" "stderr"
  assert_equal "" "$(log_of published)" "the published packages"
}

test_a_binary_without_the_executable_bit_stops_before_publishing() {
  prepare_npm
  run_env+=("FAKE_NO_EXEC=1")

  run_in_repo bin/release-npm --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "x-twitter-scraper-cli-darwin-x64 packs its binary without the executable bit. Fix the pack step, then rerun." \
    "$run_stderr" "stderr"
}

test_npm_files_changed_after_the_tag_stop_the_run() {
  prepare_npm
  publish_env
  echo "Changed after the release." >>"$sandbox_repo/npm/x-twitter-scraper-cli/README.md"
  sandbox_git commit --quiet --all --message "Change the npm README"
  sandbox_git push --quiet origin main

  run_in_repo bin/release-npm

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "npm/ or LICENSE changed after v1.2.3. Release a new version, then publish it." "$run_stderr" "stderr"
  assert_not_contains "$(log_of order)" "gh release download" "the steps"
}

test_a_release_from_before_the_npm_packages_uses_the_launcher_on_main() {
  prepare_npm 1.2.3 tag-before-npm

  run_in_repo bin/release-npm --dry-run

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "v1.2.3 predates the npm packages, so they use the launcher on this commit." "stdout"
  assert_contains "$(package_json x-twitter-scraper-cli | jq -r .bin.\"x-twitter-scraper\")" "bin/x-twitter-scraper.js" "the launcher bin"
}

test_check_passes_when_the_keychain_holds_the_npm_token_item() {
  prepare_npm
  publish_env

  run_in_repo bin/release-npm --check

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "npm publishing is ready." "stdout"
  assert_equal "security find-generic-password -s test-npm-service -a test-npm-account" "$(log_of order)" "the steps"
}

test_check_refuses_a_keychain_without_the_item() {
  prepare_npm
  publish_env
  run_env+=("FAKE_NO_ITEM=1")

  run_in_repo bin/release-npm --check

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The Keychain holds no item for NPM_KEYCHAIN_SERVICE and NPM_KEYCHAIN_ACCOUNT. Add the npm token, then rerun." "$run_stderr" "stderr"
}

run_tests
