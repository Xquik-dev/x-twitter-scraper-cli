# Xquik CLI: Twitter Search, Followers & X Automation

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/13732/badge)](https://www.bestpractices.dev/projects/13732)
[![Skills.sh x-twitter-scraper Skill](https://skills.sh/b/xquik-dev/x-twitter-scraper)](https://skills.sh/xquik-dev/x-twitter-scraper)

Use the Xquik CLI for Twitter search, profile tweets, user lookup, and follower exports.
Download media, monitor X, manage webhooks, and run X automation from a terminal.
It wraps the [Xquik REST API](https://xquik.com) as a command-line Twitter API alternative.

[Stainless](https://www.stainless.com/) generates this CLI.

<!-- x-release-please-start-version -->

## Choose the CLI

Use it for terminals, shell scripts, and one-off automation.
Choose an SDK for long-running apps or custom error handling.
Run `--help` before scripting generated commands and flags.

## Installation

### Install With Go

Install [Go](https://go.dev/doc/install) 1.22 or later.

```sh
go install 'github.com/Xquik-dev/x-twitter-scraper-cli/cmd/x-twitter-scraper@latest'
```

Go writes the binary to `$GOBIN` or `$(go env GOPATH)/bin`.
If the command is unavailable, add that directory to `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

<!-- x-release-please-end -->

### Run Locally

Run `./scripts/run args...` from the cloned repository.

## Usage

The CLI follows a resource-based command structure:

```sh
x-twitter-scraper [resource] <command> [flags...]
```

```sh
x-twitter-scraper x:tweets search \
  --api-key 'My API Key' \
  --q from:elonmusk \
  --limit 10
```

## Common Commands

Run each command with `--help` before adding flags.

| Task | Start Command |
| --- | --- |
| How do I search tweets? | `x-twitter-scraper x:tweets search --help` |
| How do I get a user's posts? | `x-twitter-scraper x:users retrieve-tweets --help` |
| How do I scrape followers? | `x-twitter-scraper x:users retrieve-followers --help` |
| How do I scrape following accounts? | `x-twitter-scraper x:users retrieve-following --help` |
| How do I read my home timeline? | `x-twitter-scraper x get-home-timeline --help` |
| How do I export a large dataset? | `x-twitter-scraper extractions estimate-cost --help` |
| How do I monitor an account? | `x-twitter-scraper monitors create --help` |
| How do I post or reply? | `x-twitter-scraper x:tweets create --help` |

### Environment Variables

| Environment variable             | Description            | Required | Default value |
| -------------------------------- | ---------------------- | -------- | ------------- |
| `X_TWITTER_SCRAPER_API_KEY`      | Xquik API key          | no       | `null`        |
| `X_TWITTER_SCRAPER_BEARER_TOKEN` | OAuth 2.1 access token | no       | `null`        |

### Global Flags

- `--api-key`: Set the Xquik API key.
- `--bearer-token`: Set an OAuth 2.1 access token.
- `--help`: Show command-line usage.
- `--debug`: Include HTTP request & response details in logs.
- `--version`, `-v`: Show the CLI version.
- `--base-url`: Use a custom API base URL.
- `--format`: Set the output format.
- `--format-error`: Set the error format.
- `--transform`: Filter output with [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md).
- `--transform-error`: Filter errors with GJSON syntax.

Credential flags also read their matching environment variables.
Debug logs may contain private response data. Do not share them.

### Passing Files as Arguments

Prefix a path with `@` to send its contents:

```bash
x-twitter-scraper <command> --arg @abe.jpg
```

Use file references inside JSON or YAML:

```bash
x-twitter-scraper <command> --arg '{image: "@abe.jpg"}'
# Equivalent:
x-twitter-scraper <command> <<YAML
arg:
  image: "@abe.jpg"
YAML
```

Escape a leading `@` when sending a literal value:

```bash
x-twitter-scraper <command> --username '\@abe'
```

#### Explicit Encoding

JSON endpoints detect text and binary files.
Use `@file://` for text, `@data://` for base64, and 3 slashes for absolute paths.

```bash
x-twitter-scraper <command> --arg @data://file.txt
```

## Link Another Go SDK Version

Use `./scripts/link` to test another Xquik Go SDK version:

```bash
./scripts/link github.com/org/repo@version
./scripts/link ../path/to/xtwitterscraper-go
```

Without an argument, the script uses `../x-twitter-scraper-go`.

## Support & Policies

- [Organization support policy](https://github.com/Xquik-dev/.github/blob/main/SUPPORT.md)
- [Security policy](SECURITY.md)
- [Contribution guide](https://github.com/Xquik-dev/.github/blob/main/CONTRIBUTING.md)

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
