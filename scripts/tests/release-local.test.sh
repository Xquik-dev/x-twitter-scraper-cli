#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2026 Xquik contributors
#
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

# shellcheck source=scripts/tests/support.sh
source "$(dirname -- "${BASH_SOURCE[0]}")/support.sh"

readonly test_token="gho_test_token_for_release_local"
readonly repository="Xquik-dev/x-twitter-scraper-cli"

prepare_release() {
  local version="${1:-1.2.3}"
  copy_file bin/release-local
  copy_file scripts/release-listing
  stub bin/check-local <<'EOF'
#!/usr/bin/env bash
echo "check-local${*:+ $*}" >>"$FAKE_LOGS/order"
EOF
  stub scripts/audit <<'EOF'
#!/usr/bin/env bash
echo "audit" >>"$FAKE_LOGS/order"
EOF
  stub bin/check-release-debug-info <<'EOF'
#!/usr/bin/env bash
echo "check-release-debug-info" >>"$FAKE_LOGS/order"
EOF
  stub bin/check-release-archives <<'EOF'
#!/usr/bin/env bash
echo "check-release-archives" >>"$FAKE_LOGS/order"
if [[ -n "${FAKE_LOCAL_OWNER:-}" ]]; then
  echo "1 archive entries name a user or group other than root. Set owner and group to root for every archive entry in .goreleaser.yml, then rerun." >&2
  exit 1
fi
EOF
  write_file .release-please-manifest.json <<EOF
{".": "$version"}
EOF
  write_file pkg/cmd/version.go <<EOF
package cmd

const Version = "$version"
EOF
  write_file CHANGELOG.md <<EOF
# Changelog

## $version (2026-10-13)

### Bug Fixes

* keep every flag on the field it filled before

## 0.0.1 (2026-01-01)

* first release
EOF
  write_file .gitignore <<'EOF'
dist/
EOF
  fake go <<'EOF'
#!/bin/sh
exit 0
EOF
  fake goreleaser <<'EOF'
#!/bin/sh
{
  printf 'goreleaser %s' "$*"
  [ -n "${GORELEASER_CURRENT_TAG:-}" ] && printf ' with GORELEASER_CURRENT_TAG=%s' "$GORELEASER_CURRENT_TAG"
  printf '\n'
} >>"$FAKE_LOGS/order"
echo "$*" >>"$FAKE_LOGS/goreleaser-arguments"
version="$(sed -n 's/^const Version = "\(.*\)"$/\1/p' pkg/cmd/version.go)"
rm -rf dist
mkdir -p dist/control
prefix="x-twitter-scraper_${version}"
stamp="$(git describe --tags --exact-match HEAD 2>/dev/null || echo "pseudo-version of $(git rev-parse --short HEAD)")"
printf 'archive %s built as %s\n' "$version" "$stamp" >"dist/${prefix}_linux_amd64.tar.gz"
printf 'Package: x-twitter-scraper\nVersion: %s\nMaintainer: support@xquik.com\nHomepage: https://xquik.com\nDescription: %s\n' \
  "$version" "${FAKE_DESCRIPTION:-X API & Twitter API CLI by Xquik for tests. Not affiliated with X Corp.}" >dist/control/control
(cd dist/control && touch -t 202601010000 control && tar -cf - ./control | gzip -n >control.tar.gz &&
  echo "2.0" >debian-binary && : >data.tar.gz && touch -t 202601010000 debian-binary control.tar.gz data.tar.gz &&
  ar -rc "../${prefix}_linux_amd64.deb" debian-binary control.tar.gz data.tar.gz)
rm -rf dist/control
(cd dist && shasum -a 256 "${prefix}_linux_amd64.tar.gz" "${prefix}_linux_amd64.deb" >"${prefix}_checksums.txt")
cat >dist/artifacts.json <<JSON
[
  {"name": "x-twitter-scraper", "path": "dist/linux_amd64/x-twitter-scraper", "type": "Binary", "goos": "linux", "goarch": "amd64"},
  {"name": "${prefix}_linux_amd64.tar.gz", "path": "dist/${prefix}_linux_amd64.tar.gz", "type": "Archive", "goarch": "amd64"},
  {"name": "${prefix}_linux_amd64.deb", "path": "dist/${prefix}_linux_amd64.deb", "type": "Linux Package", "goarch": "amd64"},
  {"name": "${prefix}_checksums.txt", "path": "dist/${prefix}_checksums.txt", "type": "Checksum"}
]
JSON
[ -z "${FAKE_DIRTY_BUILD:-}" ] || echo "changed" >>CHANGELOG.md
case " $* " in
  *" --skip="*) ;;
  *)
    if [ "${GITHUB_TOKEN:-}" = "gho_test_token_for_release_local" ]; then
      echo "the expected token" >"$FAKE_LOGS/goreleaser-token"
    else
      echo "no token or a wrong token" >"$FAKE_LOGS/goreleaser-token"
    fi
    printf '%s\n' "${prefix}_checksums.txt" "${prefix}_linux_amd64.deb" "${prefix}_linux_amd64.tar.gz" >"$FAKE_LOGS/uploaded"
    cp "dist/${prefix}_checksums.txt" "$FAKE_LOGS/uploaded-checksums"
    ;;
esac
EOF
  fake gh <<'EOF'
#!/bin/sh
echo "gh $*" >>"$FAKE_LOGS/order"
case "$1 $2" in
  "api repos/Xquik-dev/x-twitter-scraper-cli/releases/tags/"*)
    if [ -n "${FAKE_TAG_LOOKUP_FAILS:-}" ]; then
      echo "gh: Bad Gateway (HTTP 502)" >&2
      exit 1
    elif [ -n "${FAKE_PUBLISHED:-}" ]; then
      echo '{"tag_name": "v1.2.3", "draft": false}'
    else
      echo "gh: Not Found (HTTP 404)" >&2
      exit 1
    fi
    ;;
  "auth status") exit 0 ;;
  "auth token") echo "gho_test_token_for_release_local" ;;
  "api --paginate") echo "${FAKE_RELEASES:-[[]]}" ;;
  "release create") : >"$FAKE_LOGS/draft" ;;
  "release view")
    if [ -n "${FAKE_VIEW_ASSETS:-}" ]; then
      echo "$FAKE_VIEW_ASSETS"
    else
      jq -R . "$FAKE_LOGS/uploaded" | jq -cs 'sort'
    fi
    ;;
  "release download")
    if [ -n "${FAKE_DOWNLOADED_CHECKSUMS:-}" ]; then
      echo "$FAKE_DOWNLOADED_CHECKSUMS"
    else
      cat "$FAKE_LOGS/uploaded-checksums"
    fi
    ;;
  "release edit") ;;
  *)
    echo "unexpected gh $*" >&2
    exit 97
    ;;
esac
EOF
  commit_repo
  add_origin
}

remote_tags() {
  sandbox_git ls-remote --tags origin
}

order_without_paths() {
  log_of order | sed "s|$sandbox|SANDBOX|g; s|--notes-file [^ ]*|--notes-file NOTES|g"
}

test_a_dry_run_builds_every_asset_prints_the_listing_and_releases_nothing() {
  local digest
  prepare_release

  run_in_repo bin/release-local --dry-run

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  digest="$(cd "$sandbox_repo/dist" && shasum -a 256 x-twitter-scraper_1.2.3_linux_amd64.deb | cut -d ' ' -f 1)"
  assert_contains "$run_stdout" "Version  1.2.3
Tag  v1.2.3
Release assets
  x-twitter-scraper_1.2.3_checksums.txt  " "stdout"
  assert_contains "$run_stdout" "  x-twitter-scraper_1.2.3_linux_amd64.deb  " "stdout"
  assert_contains "$run_stdout" " bytes  sha256 $digest" "stdout"
  assert_contains "$run_stdout" "Asset count  3
Linux package fields
  Package  x-twitter-scraper
  Version  1.2.3
  Maintainer  support@xquik.com
  Homepage  https://xquik.com
  Description  X API & Twitter API CLI by Xquik for tests. Not affiliated with X Corp.
GitHub release title  v1.2.3
GitHub release notes from CHANGELOG.md

  ### Bug Fixes

  * keep every flag on the field it filled before" "stdout"
  assert_ends_with "$run_stdout" "Dry run finished. Nothing was published, tagged or released." "stdout"
  assert_equal "gh auth status --hostname github.com
gh api repos/$repository/releases/tags/v1.2.3
gh api --paginate --slurp repos/$repository/releases?per_page=100
check-local --no-status
audit
goreleaser release --clean --skip=publish,validate with GORELEASER_CURRENT_TAG=v1.2.3
check-release-debug-info
check-release-archives" "$(order_without_paths)" "the steps"
  assert_equal "" "$(remote_tags)" "the tags on origin"
  assert_equal "" "$(sandbox_git tag)" "the local tags"
}

test_a_dry_run_lists_the_same_checksums_as_a_release_of_the_same_commit() {
  local dry_run_checksums
  prepare_release

  run_in_repo bin/release-local --dry-run
  dry_run_checksums="$(cat "$sandbox_repo/dist/x-twitter-scraper_1.2.3_checksums.txt")"
  run_in_repo bin/release-local

  assert_equal 0 "$run_status" "the exit code of the release"
  assert_contains "$(cat "$sandbox_repo/dist/x-twitter-scraper_1.2.3_linux_amd64.tar.gz")" "built as v1.2.3" "the released archive"
  assert_equal "$dry_run_checksums" "$(log_of uploaded-checksums)" "the uploaded checksums"
}

test_a_dry_run_keeps_a_local_tag_that_was_there_before() {
  prepare_release
  sandbox_git tag --annotate v1.2.3 --message "Release v1.2.3"

  run_in_repo bin/release-local --dry-run

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "==> Reusing the local annotated tag v1.2.3" "stdout"
  assert_equal "tag" "$(sandbox_git cat-file -t v1.2.3)" "the local tag type"
  assert_equal "" "$(remote_tags)" "the tags on origin"
}

test_a_failed_dry_run_removes_the_tag_it_created() {
  prepare_release
  run_env=("FAKE_LOCAL_OWNER=1")

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "" "$(sandbox_git tag)" "the local tags"
}

test_an_archive_entry_owned_by_a_local_account_stops_a_release_before_the_draft() {
  prepare_release
  run_env=("FAKE_LOCAL_OWNER=1")

  run_in_repo bin/release-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "1 archive entries name a user or group other than root. Set owner and group to root for every archive entry in .goreleaser.yml, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of order)" "gh release create" "the steps"
  assert_equal "" "$(remote_tags)" "the tags on origin"
}

test_a_dry_run_reports_a_release_that_github_already_published_and_goes_on() {
  prepare_release
  run_env=("FAKE_PUBLISHED=1")

  run_in_repo bin/release-local --dry-run

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "GitHub release v1.2.3 is already published. A real release would stop here." "stdout"
  assert_ends_with "$run_stdout" "Dry run finished. Nothing was published, tagged or released." "stdout"
}

test_a_release_uploads_with_the_token_in_the_environment_then_tags_and_publishes() {
  prepare_release

  run_in_repo bin/release-local

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_ends_with "$run_stdout" "Released v1.2.3 at https://github.com/$repository/releases/tag/v1.2.3" "stdout"
  assert_equal "gh auth status --hostname github.com
gh api repos/$repository/releases/tags/v1.2.3
gh api --paginate --slurp repos/$repository/releases?per_page=100
check-local
audit
goreleaser release --clean --skip=publish,validate with GORELEASER_CURRENT_TAG=v1.2.3
check-release-debug-info
check-release-archives
gh release create v1.2.3 --repo $repository --draft --title v1.2.3 --target $(sandbox_git rev-parse HEAD) --notes-file NOTES
gh auth token --hostname github.com
goreleaser release --clean
gh release view v1.2.3 --repo $repository --json assets --jq [.assets[].name] | sort
gh release download v1.2.3 --repo $repository --pattern x-twitter-scraper_1.2.3_checksums.txt --output -
gh release edit v1.2.3 --repo $repository --draft=false --verify-tag --title v1.2.3 --notes-file NOTES --latest" \
    "$(order_without_paths)" "the steps"
  assert_equal "the expected token" "$(log_of goreleaser-token)" "the token GoReleaser saw"
  assert_not_contains "$(log_of goreleaser-arguments)" "$test_token" "the GoReleaser arguments"
  assert_not_contains "$run_stdout$run_stderr" "$test_token" "the output"
  assert_contains "$(remote_tags)" "refs/tags/v1.2.3" "the tags on origin"
  assert_equal "tag" "$(sandbox_git cat-file -t v1.2.3)" "the tag type"
}

test_a_release_skips_the_upload_when_the_draft_already_has_every_asset() {
  prepare_release
  printf '%s\n' x-twitter-scraper_1.2.3_checksums.txt x-twitter-scraper_1.2.3_linux_amd64.deb \
    x-twitter-scraper_1.2.3_linux_amd64.tar.gz >"$sandbox_logs/uploaded"
  run_env=('FAKE_RELEASES=[[{"tag_name": "v1.2.3", "draft": true, "assets": [{"name": "x-twitter-scraper_1.2.3_checksums.txt"}, {"name": "x-twitter-scraper_1.2.3_linux_amd64.deb"}, {"name": "x-twitter-scraper_1.2.3_linux_amd64.tar.gz"}]}]]')
  sandbox_git tag --annotate v1.2.3 --message "Release v1.2.3"
  mkdir -p "$sandbox/setup-logs"
  (
    cd "$sandbox_repo"
    env PATH="$sandbox_bin:/usr/bin:/bin" FAKE_LOGS="$sandbox/setup-logs" goreleaser release --skip=publish >/dev/null
    cp dist/x-twitter-scraper_1.2.3_checksums.txt "$sandbox_logs/uploaded-checksums"
    rm -rf dist
  )

  run_in_repo bin/release-local

  assert_equal "" "$run_stderr" "stderr"
  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$run_stdout" "==> Skipping the upload because draft release v1.2.3 already has every asset" "stdout"
  assert_not_contains "$(log_of order)" "gh release create" "the steps"
  assert_not_contains "$(log_of order)" "gh auth token" "the steps"
  assert_contains "$(log_of order)" "gh release edit v1.2.3" "the steps"
}

test_a_published_release_stops_a_real_run_before_the_checks() {
  prepare_release
  run_env=("FAKE_PUBLISHED=1")

  run_in_repo bin/release-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "GitHub release v1.2.3 is already published. Bump the version, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of order)" "check-local" "the steps"
}

test_a_failed_release_lookup_stops_the_run_before_the_checks() {
  prepare_release
  run_env=("FAKE_TAG_LOOKUP_FAILS=1")

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "GitHub did not answer for release v1.2.3. Check the network and gh auth status, then rerun." "$run_stderr" "stderr"
  assert_not_contains "$(log_of order)" "check-local" "the steps"
}

test_two_drafts_for_the_tag_stop_the_run() {
  prepare_release
  run_env=('FAKE_RELEASES=[[{"tag_name": "v1.2.3", "draft": true, "assets": []}, {"tag_name": "v1.2.3", "draft": true, "assets": []}]]')

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Several draft releases use v1.2.3. Delete the extra drafts on GitHub, then rerun." "$run_stderr" "stderr"
}

test_a_draft_missing_an_asset_after_the_upload_stops_before_the_tag() {
  prepare_release
  run_env=('FAKE_VIEW_ASSETS=["x-twitter-scraper_1.2.3_checksums.txt"]')

  run_in_repo bin/release-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Draft release v1.2.3 does not hold exactly the built assets. Fix the draft on GitHub, then rerun." "$run_stderr" "stderr"
  assert_equal "" "$(remote_tags)" "the tags on origin"
  assert_not_contains "$(log_of order)" "gh release edit" "the steps"
}

test_checksums_that_differ_on_the_draft_stop_before_the_tag() {
  prepare_release
  run_env=("FAKE_DOWNLOADED_CHECKSUMS=0000  x-twitter-scraper_1.2.3_linux_amd64.deb")

  run_in_repo bin/release-local

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The checksums on draft v1.2.3 differ from this build. Delete the draft assets, then rerun." "$run_stderr" "stderr"
  assert_equal "" "$(remote_tags)" "the tags on origin"
}

test_a_prerelease_is_published_without_the_latest_label() {
  prepare_release 1.3.0-beta.1

  run_in_repo bin/release-local

  assert_equal 0 "$run_status" "the exit code"
  assert_contains "$(log_of order)" "--draft=false --verify-tag --title v1.3.0-beta.1" "the steps"
  assert_ends_with "$(log_of order)" "--prerelease --latest=false" "the steps"
}

test_a_package_description_without_the_listing_words_stops_the_run() {
  prepare_release
  run_env=("FAKE_DESCRIPTION=no description given")

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The package description must start with \"X API & Twitter API\" and end with \"Not affiliated with X Corp.\" Fix the nfpms section of .goreleaser.yml." "$run_stderr" "stderr"
  assert_not_contains "$run_stdout" "Dry run finished" "stdout"
}

test_a_build_that_changes_tracked_files_stops_the_run() {
  prepare_release
  run_env=("FAKE_DIRTY_BUILD=1")

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The build changed tracked files. Fix the build so it leaves the tree alone, then rerun." "$run_stderr" "stderr"
}

test_a_tag_that_origin_already_has_stops_the_run() {
  prepare_release
  sandbox_git tag v1.2.3
  sandbox_git push --quiet origin v1.2.3
  sandbox_git tag -d v1.2.3 >/dev/null

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "Tag v1.2.3 already exists on origin. Bump the version, then rerun." "$run_stderr" "stderr"
}

test_a_working_tree_with_changes_stops_the_run() {
  prepare_release
  write_file notes.txt <<'EOF'
Uncommitted work.
EOF

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The working tree has changes. Commit or remove them, then rerun." "$run_stderr" "stderr"
}

test_a_branch_other_than_main_stops_the_run() {
  prepare_release
  sandbox_git switch --quiet -c feature

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "The current branch is not main. Switch to main, then rerun." "$run_stderr" "stderr"
}

test_a_main_that_origin_lacks_stops_the_run() {
  prepare_release
  write_file extra.txt <<'EOF'
Not pushed.
EOF
  sandbox_git add extra.txt
  sandbox_git commit --quiet --message "Local only"

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "main differs from origin/main. Pull or push main, then rerun." "$run_stderr" "stderr"
}

test_a_version_without_a_changelog_section_stops_the_run() {
  prepare_release
  write_file CHANGELOG.md <<'EOF'
# Changelog

## 0.0.1 (2026-01-01)

* first release
EOF
  sandbox_git commit --quiet --all --message "Drop the section"
  sandbox_git push --quiet origin main

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "CHANGELOG.md has no section for 1.2.3. Merge the release pull request first." "$run_stderr" "stderr"
}

test_a_version_file_that_disagrees_stops_the_run() {
  prepare_release
  write_file pkg/cmd/version.go <<'EOF'
package cmd

const Version = "1.2.2"
EOF
  sandbox_git commit --quiet --all --message "Old version"
  sandbox_git push --quiet origin main

  run_in_repo bin/release-local --dry-run

  assert_equal 1 "$run_status" "the exit code"
  assert_equal "pkg/cmd/version.go does not say 1.2.3. Merge the release pull request first." "$run_stderr" "stderr"
}

test_an_unknown_option_stops_with_the_usage_line() {
  prepare_release

  run_in_repo bin/release-local --publish

  assert_equal 2 "$run_status" "the exit code"
  assert_equal "Unknown option --publish. Run bin/release-local or bin/release-local --dry-run." "$run_stderr" "stderr"
  assert_equal "" "$(log_of order)" "the steps"
}

run_tests
