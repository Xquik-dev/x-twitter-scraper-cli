#!/usr/bin/env node
"use strict";

const { spawnSync } = require("node:child_process");
const { constants } = require("node:os");

const target = `${process.platform}-${process.arch}`;
const binary = process.platform === "win32" ? "x-twitter-scraper.exe" : "x-twitter-scraper";

function installedBinary() {
  try {
    return require.resolve(`x-twitter-scraper-cli-${target}/bin/${binary}`);
  } catch {
    return undefined;
  }
}

const path = installedBinary();
if (path === undefined) {
  process.stderr.write(
    [
      `x-twitter-scraper-cli found no binary for ${target}.`,
      "Reinstall it without --omit=optional or --no-optional.",
      "On another platform, download a build from https://github.com/Xquik-dev/x-twitter-scraper-cli/releases.",
      "",
    ].join("\n"),
  );
  process.exit(1);
}

const result = spawnSync(path, process.argv.slice(2), { stdio: "inherit" });
if (result.error !== undefined) {
  process.stderr.write(`x-twitter-scraper could not start. ${result.error.message}\n`);
  process.exit(1);
}
if (result.signal !== null) {
  process.kill(process.pid, result.signal);
  process.exit(128 + (constants.signals[result.signal] ?? 0));
}
process.exit(result.status ?? 1);
