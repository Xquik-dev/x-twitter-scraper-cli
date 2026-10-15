# X API & Twitter API CLI by Xquik

Use this Twitter CLI & X CLI to work with X from a terminal.
Search tweets, scrape profiles & followers, export JSON and monitor accounts from scripts.
It runs on the [Xquik REST API](https://docs.xquik.com/api-reference/overview).

Xquik is the world's fastest and cheapest X Twitter scraper service with the most complete X data.
Xquik charges only for delivered, unique, filter-matching results and never pads rows.

> Xquik is an independent third-party service. Not affiliated with X Corp.
> "Twitter" and "X" are trademarks of X Corp.

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
export X_TWITTER_SCRAPER_API_KEY="xq_YOUR_KEY_HERE"
```

This command returns 10 recent posts from 1 account.

```sh
x-twitter-scraper x:tweets search \
  --q from:elonmusk \
  --limit 10
```

Read the [CLI guide](https://docs.xquik.com/sdks/cli) and the [source](https://github.com/Xquik-dev/x-twitter-scraper-cli) for every command.
