# X API & Twitter API CLI by Xquik

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13732/badge)](https://www.bestpractices.dev/projects/13732)

Use this Twitter CLI & X CLI to work with X from a terminal.
It works as a Twitter scraper & tweet scraper for search, profiles, followers & lists.
Scrape Twitter timelines, communities & trends into JSON, YAML or an interactive view.
Download media, monitor accounts, manage webhooks & run X automation from scripts.
It is a Twitter API alternative built on the [Xquik REST API](https://docs.xquik.com/api-reference/overview).

Xquik is the world's fastest and cheapest X Twitter scraper service with the most complete X data.
Xquik charges only for delivered, unique, filter-matching results and never pads rows.

> Xquik is an independent third-party service. Not affiliated with X Corp.
> "Twitter" and "X" are trademarks of X Corp.

<!-- trust:start -->

## Trust & usage

[![Twitter API package downloads](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fdownloads.json)](https://xquik.com/trust#downloads)
[![Twitter scraper users on Apify](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fapify-users.json)](https://apify.com/xquik)
[![Twitter scraper runs on Apify](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fapify-runs.json)](https://apify.com/xquik)
[![Tweet scraper rating on Apify](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fapify-rating.json)](https://apify.com/xquik/x-tweet-scraper/reviews)
[![Xquik Twitter API accounts](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Faccounts.json)](https://xquik.com/trust#accounts)
[![Twitter Skill installs](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fskill-installs.json)](https://xquik.com/trust#skill-installs)
[![Xquik verified revenue on TrustMRR](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Frevenue.json)](https://trustmrr.com/startup/xquik)
[![Twitter API median response time](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Frest-p50.json)](https://xquik.com/trust#speed)
[![Xquik GitHub stars](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fgithub-stars.json)](https://github.com/Xquik-dev)
[![x-twitter-scraper-cli downloads on GitHub Releases](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Freleases-x-twitter-scraper-cli.json)](https://github.com/Xquik-dev/x-twitter-scraper-cli/releases)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13732/badge)](https://www.bestpractices.dev/projects/13732)

- 51,000+ downloads of Xquik's Twitter API packages on npm, PyPI, RubyGems, NuGet, Packagist & Terraform ([source](https://xquik.com/trust#downloads)) <!-- trust:downloads 51020 -->
- 7,800+ users of Xquik's Twitter scraper Actors on Apify ([source](https://apify.com/xquik)) <!-- trust:apify-users 7800 -->
- 5,700,000+ runs of Xquik's Twitter scraper Actors on Apify ([source](https://apify.com/xquik)) <!-- trust:apify-runs 5704701 -->
- 4.66 of 5 for the X Tweet Scraper from 16 Apify reviews ([source](https://apify.com/xquik/x-tweet-scraper/reviews)) <!-- trust:apify-rating 4.66 16 -->
- 3,900+ accounts on xquik.com ([source](https://xquik.com/trust#accounts)) <!-- trust:accounts 3924 -->
- 1,600+ installs of Xquik's Twitter Skills on skills.sh ([source](https://xquik.com/trust#skill-installs)) <!-- trust:skill-installs 1612 -->
- $38,000+ in revenue verified by TrustMRR ([source](https://trustmrr.com/startup/xquik)) <!-- trust:revenue 38115 -->
- 501 ms median REST API read time over the last 7 days ([source](https://xquik.com/trust#speed)) <!-- trust:rest-p50 501 -->
- 370+ GitHub stars across Xquik's repos ([source](https://github.com/Xquik-dev)) <!-- trust:github-stars 374 -->
- 67 downloads of x-twitter-scraper-cli on GitHub Releases ([source](https://github.com/Xquik-dev/x-twitter-scraper-cli/releases)) <!-- trust:releases-x-twitter-scraper-cli 67 -->
- This repo meets the OpenSSF Best Practices passing criteria ([source](https://www.bestpractices.dev/projects/13732)) <!-- trust:openssf-x-twitter-scraper-cli 1 -->

Checked on 2026-10-15. See each number with its date & source on [xquik.com/trust](https://xquik.com/trust) or as [JSON](https://xquik.com/trust.json).

<!-- trust:end -->

## Installation

### Install with npm

Install [Node.js](https://nodejs.org) 18 or later.

```sh
npm install --global x-twitter-scraper-cli
```

Or run the CLI once without installing it.

```sh
npx x-twitter-scraper-cli --help
```

npm installs only the build for your platform. Nothing downloads when the CLI runs.

### Install with Go

Install [Go](https://go.dev/doc/install) 1.21 or later. Go downloads the toolchain this CLI needs.

```sh
go install 'github.com/Xquik-dev/x-twitter-scraper-cli/cmd/x-twitter-scraper@latest'
```

Go writes the binary to `$GOBIN` or `$(go env GOPATH)/bin`.
If the command is unavailable, add that directory to `PATH`.

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Download a release

Each [GitHub release](https://github.com/Xquik-dev/x-twitter-scraper-cli/releases) has archives for macOS, Linux & Windows.
Linux users can also pick a `.deb`, `.rpm`, `.apk` or Arch Linux package.
Compare each download with the release's checksums file.

From a cloned repository, run `./scripts/run args...` instead.

## Usage

The CLI follows a resource-based command structure.

```sh
x-twitter-scraper [resource] <command> [flags...]
```

This command returns 10 recent posts from 1 account.

```sh
x-twitter-scraper x:tweets search \
  --q from:elonmusk \
  --limit 10
```

## Common commands

Run each command with `--help` before adding flags.

| Task | Start command |
| --- | --- |
| Search tweets | `x-twitter-scraper x:tweets search --help` |
| Get a user's posts | `x-twitter-scraper x:users retrieve-tweets --help` |
| Scrape Twitter followers | `x-twitter-scraper x:users retrieve-followers --help` |
| Scrape following accounts | `x-twitter-scraper x:users retrieve-following --help` |
| Read a home timeline | `x-twitter-scraper x get-home-timeline --help` |
| Export a large dataset | `x-twitter-scraper extractions estimate-cost --help` |
| Monitor an account | `x-twitter-scraper monitors create --help` |
| Post or reply | `x-twitter-scraper x:tweets create --help` |

### Environment variables

| Environment variable             | Description            | Required | Default value |
| -------------------------------- | ---------------------- | -------- | ------------- |
| `X_TWITTER_SCRAPER_API_KEY`      | Xquik API key          | no       | `null`        |
| `X_TWITTER_SCRAPER_BEARER_TOKEN` | OAuth 2.1 access token | no       | `null`        |

### Global flags

- `--api-key` sets the Xquik API key.
- `--bearer-token` sets an OAuth 2.1 access token.
- `--help` shows command-line usage.
- `--debug` adds HTTP request & response details to logs.
- `--version` or `-v` shows the CLI version.
- `--base-url` sends requests to another API base URL.
- `--format` sets the output format.
- `--format-error` sets the error format.
- `--transform` filters output with [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md).
- `--transform-error` filters errors with GJSON syntax.
- `--raw-output` or `-r` prints a string result without JSON quotes.

Credential flags also read their matching environment variables.
Debug logs may contain private response data. Do not share them.

### Passing files as arguments

Prefix a path with `@` to send its contents.

```bash
x-twitter-scraper <command> --arg @abe.jpg
```

Use file references inside JSON or YAML.

```bash
x-twitter-scraper <command> --arg '{image: "@abe.jpg"}'
# The same input as YAML.
x-twitter-scraper <command> <<YAML
arg:
  image: "@abe.jpg"
YAML
```

Escape a leading `@` when sending a literal value.

```bash
x-twitter-scraper <command> --username '\@abe'
```

JSON endpoints detect text and binary files.
Use `@file://` for text, `@data://` for base64, and 3 slashes for absolute paths.

```bash
x-twitter-scraper <command> --arg @data://file.txt
```

## Test another Go SDK version

Pass a module version or local path to `./scripts/link`.

```bash
./scripts/link github.com/org/repo@version
./scripts/link ../path/to/xtwitterscraper-go
```

The script defaults to `../x-twitter-scraper-go`.

## Checks & releases

Run `bin/check-local` before you open a pull request.
It runs REUSE, lint, script tests, Go tests, coverage, `govulncheck`, 2 reproducible builds, fuzzing & CodeQL.
It posts the result as the `xquik-local-checks` commit status.
Maintainers publish releases from `main` with `bin/release-local`.
`bin/release-local --dry-run` builds every release asset and publishes nothing.
After the GitHub release, `bin/release-local` runs `bin/release-npm`.
It publishes the [x-twitter-scraper-cli](https://www.npmjs.com/package/x-twitter-scraper-cli) package and 1 package per platform.
`bin/release-npm --dry-run` packs them from the published release and publishes nothing.

## More Xquik tools

- [X (Twitter) scraper API & Twitter MCP server](https://github.com/Xquik-dev/x-twitter-scraper) with REST, MCP & Skill setup.
- Twitter API SDKs for [TypeScript](https://github.com/Xquik-dev/x-twitter-scraper-typescript), [Python](https://github.com/Xquik-dev/x-twitter-scraper-python), [Go](https://github.com/Xquik-dev/x-twitter-scraper-go), [Ruby](https://github.com/Xquik-dev/x-twitter-scraper-ruby), [Java](https://github.com/Xquik-dev/x-twitter-scraper-java), [Kotlin](https://github.com/Xquik-dev/x-twitter-scraper-kotlin), [C#](https://github.com/Xquik-dev/x-twitter-scraper-csharp) & [PHP](https://github.com/Xquik-dev/x-twitter-scraper-php).
- [Twitter advanced search operators](https://github.com/Xquik-dev/x-twitter-advanced-search) with tested examples.
- [Xquik X API & Twitter scraper](https://xquik.com) and the [Twitter API documentation](https://docs.xquik.com/api-reference/overview).

## Support & policies

- [Organization support policy](https://github.com/Xquik-dev/.github/blob/main/SUPPORT.md)
- [Security policy](SECURITY.md)
- [Contribution guide](https://github.com/Xquik-dev/.github/blob/main/CONTRIBUTING.md)

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
