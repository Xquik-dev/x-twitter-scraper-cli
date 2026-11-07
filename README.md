# X API & Twitter API CLI by Xquik

Use this Twitter CLI by Xquik to search tweets & export followers as JSON.
It runs in a terminal or a script. Public reads need no X developer account.

Xquik is the world's fastest and cheapest X Twitter scraper service with the most complete X data.
Most other Apify Actors charge before filtering or deduplication.
Xquik charges only for delivered, unique, filter-matching results and never pads rows.
Xquik is independent and not affiliated with X Corp.

<!-- trust:win:start -->
The Xquik REST API read 45 tweets per second at $0.15 per 1,000 on 2026-10-15. The fastest of 4 rival Actors read 12.5 per second. The cheapest cost $0.173 per 1,000. [Compare each run](https://xquik.com/benchmarks#benchmark-x-tweet-scraper).
<!-- trust:win:end -->

## Install the Twitter CLI

Install [Node.js](https://nodejs.org) 18 or later, then install the [npm package](https://www.npmjs.com/package/x-twitter-scraper-cli).

```sh
npm install --global x-twitter-scraper-cli
```

Or run the CLI once without installing it.

```sh
npx x-twitter-scraper-cli --help
```

npm installs only the build for your platform. Nothing downloads when the CLI runs.

## Scrape tweets from the X CLI

Create an API key in the [Xquik dashboard](https://dashboard.xquik.com) and export it.

```sh
export X_TWITTER_SCRAPER_API_KEY="your-api-key"
```

This command returns 10 recent posts from 1 account.

```sh
x-twitter-scraper x:tweets search \
  --q from:elonmusk \
  --limit 10
```

Next, read the [CLI guide](https://docs.xquik.com/sdks/cli), [REST API reference](https://docs.xquik.com/api-reference/overview) or [search operators](https://github.com/Xquik-dev/x-twitter-advanced-search).

<!-- trust:start -->

## Benchmarks & usage

[![Twitter API package downloads](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Fdownloads.json)](https://xquik.com/benchmarks#downloads)
[![Twitter API p99 server read time](https://img.shields.io/endpoint?url=https%3A%2F%2Fxquik.com%2Fbadges%2Frest-p99.json)](https://xquik.com/benchmarks#speed)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13732/badge)](https://www.bestpractices.dev/projects/13732)

- 63,000+ downloads of Xquik's Twitter API packages on npm, PyPI, RubyGems, NuGet, Packagist & Terraform ([source](https://xquik.com/benchmarks#downloads)) <!-- trust:downloads 63178 -->
- 6.4 s p99 REST API server read time over 7 days, failed reads included ([source](https://xquik.com/benchmarks#speed)) <!-- trust:rest-p99 6317 -->
- This repo meets the OpenSSF Best Practices passing criteria ([source](https://www.bestpractices.dev/projects/13732)) <!-- trust:openssf-x-twitter-scraper-cli 1 -->

Checked on 2026-11-06. See each number with its date & source on [xquik.com/benchmarks](https://xquik.com/benchmarks) or as [JSON](https://xquik.com/trust.json).

<!-- trust:end -->

## More ways to install

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

## Twitter scraper commands

Each command follows the same shape.

```sh
x-twitter-scraper [resource] <command> [flags...]
```

Run a command with `--help` before you add flags.

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

## Checks & releases

Run `bin/check-local` before you open a pull request.
It posts the result as the `xquik-local-checks` commit status.
Maintainers publish releases from `main` with `bin/release-local`.
`bin/release-local --dry-run` builds every release asset and publishes nothing.

## Questions

### Do I need an X developer account?

No. Public reads need only an Xquik API key. Private reads & X account actions
need a connected X account.

### Is this the official X API?

No. Xquik runs its own REST API. This CLI calls that API, not the official X API.

## More Xquik tools

- [X (Twitter) scraper API & Twitter MCP server](https://github.com/Xquik-dev/x-twitter-scraper) with REST, MCP & Skill setup.
- Twitter API SDKs for [TypeScript](https://github.com/Xquik-dev/x-twitter-scraper-typescript), [Python](https://github.com/Xquik-dev/x-twitter-scraper-python), [Go](https://github.com/Xquik-dev/x-twitter-scraper-go) & [Java](https://github.com/Xquik-dev/x-twitter-scraper-java).
- Twitter API SDKs for [Ruby](https://github.com/Xquik-dev/x-twitter-scraper-ruby), [Kotlin](https://github.com/Xquik-dev/x-twitter-scraper-kotlin), [C#](https://github.com/Xquik-dev/x-twitter-scraper-csharp) & [PHP](https://github.com/Xquik-dev/x-twitter-scraper-php).
- [Twitter advanced search operators](https://github.com/Xquik-dev/x-twitter-advanced-search) with tested examples.
- Twitter agent plugins for [OpenClaw](https://github.com/Xquik-dev/tweetclaw) & [Hermes Agent](https://github.com/Xquik-dev/hermes-tweet).
- Twitter integrations for [n8n](https://github.com/Xquik-dev/n8n-nodes-xquik), [Prefect](https://github.com/Xquik-dev/prefect-xquik), [Haystack](https://github.com/Xquik-dev/xquik-haystack) & [Paperclip](https://github.com/Xquik-dev/paperclip-plugin-xquik).
- Twitter integrations for [Terraform](https://github.com/Xquik-dev/terraform-provider-x-twitter-scraper), [Convex](https://github.com/Xquik-dev/xquik-convex), [Supabase](https://github.com/Xquik-dev/xquik-supabase) & [Strapi](https://github.com/Xquik-dev/strapi-plugin-xquik).
- Twitter MCP for [AI agents](https://github.com/Xquik-dev/xquik-mcp), [VS Code](https://github.com/Xquik-dev/xquik-vscode) & [Raycast](https://github.com/Xquik-dev/xquik-raycast).
- [Xquik X API & Twitter scraper](https://xquik.com).
- [Twitter API documentation](https://docs.xquik.com/api-reference/overview) and its [source](https://github.com/Xquik-dev/xquik-docs).

## Support & policies

- [Organization support policy](https://github.com/Xquik-dev/.github/blob/main/SUPPORT.md)
- [Security policy](SECURITY.md)
- [Contribution guide](https://github.com/Xquik-dev/.github/blob/main/CONTRIBUTING.md)

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
