# OpenSSF Best Practices Evidence

This repository tracks [OpenSSF Best Practices project 13732][badge].

## Automated Evidence

- `./scripts/test` runs race detection and all tests.
- `./scripts/coverage` enforces 90% statement coverage.
- CI verifies licenses, dependencies, vulnerabilities, and builds.
- CodeQL runs extended security queries.
- OpenSSF Scorecard publishes code-scanning results.
- GoReleaser builds receive GitHub artifact attestations.
- CI compares 2 consecutive release checksum manifests.

The merged Go profile currently reports 91.1% statement coverage.

## Decision Coverage Assessment

Go's standard tool does not report decision coverage.
The assessment identified the MIT-licensed [`gocove`][gocove] beta.
Its latest module has no tagged release.

The 2026-07-23 assessment used commit `a0dceb9dadca`.
Instrumentation failed for 2 of 10 packages.
The tool reported only 2 packages after that failure.
Its partial condition result was 64.4%.
That incomplete result cannot establish the required 80%.

The decision-coverage criterion remains unresolved.
Do not mark it satisfied until a complete FLOSS measurement passes.

## Organization-Level Constraints

Organization governance evidence lives in the shared `.github` repository.
Gold remains blocked until every mandatory organization criterion passes.
No repository may claim Gold before that evidence exists.

[badge]: https://www.bestpractices.dev/projects/13732
[gocove]: https://github.com/srvgit/gocove

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
