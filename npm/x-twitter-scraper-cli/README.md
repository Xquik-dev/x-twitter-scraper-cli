# X API & Twitter API CLI by Xquik

Use this Twitter CLI by Xquik to search tweets & export followers as JSON.
It runs on the [Xquik REST API](https://docs.xquik.com/api-reference/overview). Public reads need no X developer account.

Xquik is the world's fastest and cheapest X Twitter scraper service with the most complete X data.
Most other Apify Actors charge before filtering or deduplication.
Xquik charges only for delivered, unique, filter-matching results and never pads rows.
Xquik is independent and not affiliated with X Corp.

<!-- trust:win:start -->
The Xquik REST API read 45 tweets per second at $0.15 per 1,000 on 2026-10-15. The fastest of 4 rival Actors read 12.5 per second. The cheapest cost $0.173 per 1,000. [Compare each run](https://xquik.com/benchmarks#benchmark-x-tweet-scraper).
<!-- trust:win:end -->

## Install

```sh
npm install --global x-twitter-scraper-cli
```

Or run it once without installing.

```sh
npx x-twitter-scraper-cli --help
```

npm installs the build for your platform from a matching `x-twitter-scraper-cli-<os>-<cpu>` package.
Nothing downloads when you run the CLI.

## Use it

Set your API key from the [Xquik dashboard](https://dashboard.xquik.com).

```sh
export X_TWITTER_SCRAPER_API_KEY="your-api-key"
```

This command returns 10 recent posts from 1 account.

```sh
x-twitter-scraper x:tweets search \
  --q from:elonmusk \
  --limit 10
```

Next, read the [CLI guide](https://docs.xquik.com/sdks/cli) or the [CLI source on GitHub](https://github.com/Xquik-dev/x-twitter-scraper-cli) for every command.

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

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
