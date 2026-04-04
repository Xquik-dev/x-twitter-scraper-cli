# Changelog

## 0.3.0 (2026-04-04)

Full Changelog: [v0.2.0...v0.3.0](https://github.com/Xquik-dev/x-twitter-scraper-cli/compare/v0.2.0...v0.3.0)

### Features

* allow `-` as value representing stdin to binary-only file parameters in CLIs ([712746c](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/712746c9735661862689301f04ff1a90ebac1d8a))
* **api:** api update ([dcbfdc4](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/dcbfdc4f2bef587eb67db7d75bde62aaac8cb586))
* **api:** api update ([2939f30](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/2939f308c9443e7e0fe53e1ac47f43413b6d78a5))
* **api:** api update ([600dade](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/600dadedbb67605a0619ff305b6859ac2d2b22f2))
* better error message if scheme forgotten in CLI `*_BASE_URL`/`--base-url` ([07a4005](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/07a40054367cdbd2ebba05fc271f0f8924adc3b3))
* binary-only parameters become CLI flags that take filenames only ([ce4ff99](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/ce4ff99c2fd457bd041238b36255282dec76aa1a))


### Bug Fixes

* handle empty data set using `--format explore` ([fb087bb](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/fb087bb19f7783f5c6b0d34a6c800da8d6ff6337))
* use `RawJSON` when iterating items with `--format explore` in the CLI ([9b49da1](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/9b49da128f9110d987a9174a6fa84419ae337e9b))


### Chores

* mark all CLI-related tests in Go with `t.Parallel()` ([273bb76](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/273bb76c4ca2923aaf4ca99ad527b1ebb93f6be9))
* modify CLI tests to inject stdout so mutating `os.Stdout` isn't necessary ([2da8d2f](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/2da8d2f9f31119189bbf53ffae8b92817959bb2d))
* switch some CLI Go tests from `os.Chdir` to `t.Chdir` ([445ffb4](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/445ffb42a38c6b8cb90a18dbf377e1949e1b0fde))

## 0.2.0 (2026-04-01)

Full Changelog: [v0.1.0...v0.2.0](https://github.com/Xquik-dev/x-twitter-scraper-cli/compare/v0.1.0...v0.2.0)

### Features

* **api:** api update ([c92310f](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/c92310f443f869ecf343b9f470ef38bbb5b3b4ca))


### Chores

* update SDK settings ([1f68d8e](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/1f68d8e5a490ef84be5e50f612e06cb6d780150c))
* update SDK settings ([ae4c4da](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/ae4c4dadc4371af61306bc35ce86857676f3e348))

## 0.1.0 (2026-03-30)

Full Changelog: [v0.0.1...v0.1.0](https://github.com/Xquik-dev/x-twitter-scraper-cli/compare/v0.0.1...v0.1.0)

### Features

* **api:** api update ([117f17e](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/117f17ed4d15b4a76ac3d7306cf5cdc42122e51a))


### Chores

* update SDK settings ([eb5bf2a](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/eb5bf2a929db79e75e059985033c4c225baf331f))
* update SDK settings ([e20e024](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/e20e024081537cc77b42f7e119715625e9aee512))
* update SDK settings ([f5f5697](https://github.com/Xquik-dev/x-twitter-scraper-cli/commit/f5f569769f18cb84abcf8797f99e321630202e81))
