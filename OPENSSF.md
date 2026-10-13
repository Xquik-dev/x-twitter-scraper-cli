# OpenSSF Best Practices Evidence

This file records the repository's Gold assessment.

The official project is [bestpractices.dev project 13732][badge].

The assessment date is 2026-07-23.

## Eligibility

This active public CLI qualifies for the OpenSSF Best Practices badge.
No OpenSSF exclusion applies.

## Verified Technical Controls

| Area | Evidence |
| --- | --- |
| License | Apache-2.0, REUSE 3.3 metadata, & bundled BSD-3-Clause terms |
| Contribution process | `CONTRIBUTING.md`, DCO sign-off, & `bin/check-local` before every merge |
| Security reporting | Private GitHub reporting, response targets, & safe harbor in `SECURITY.md` |
| Build | Pinned Go modules, checksum verification, & `./scripts/build` |
| Tests | Race-enabled tests & regression coverage through `./scripts/test` |
| Statement coverage | `./scripts/coverage` enforces 90%. The current result is 91.2%. |
| Branch coverage | `./scripts/branch-coverage` enforces 80%. The current result is 879/1,096 (80.20%). |
| Static analysis | `go vet`, `shellcheck` & CodeQL security-extended queries through `bin/check-local` |
| Dynamic analysis | Race detection in `./scripts/test` & a native Go fuzz target that `bin/check-local` runs for 15 seconds |
| Dependency review | `go mod verify` & `govulncheck` through `./scripts/audit`, with reviewed exceptions in `audit-exceptions.json` |
| Local checks | `bin/check-local` runs licensing, lint, script, test, audit, build, fuzzing & CodeQL checks, then posts the `xquik-local-checks` status |
| Releases | `bin/release-local` publishes checksums & an annotated tag from `main` |
| Reproducibility | `bin/check-local` builds each GoReleaser snapshot twice & compares checksums |
| Portability | `./scripts/test` compiles every package for Windows AMD64 |
| Two-factor authentication | The Xquik-dev organization requires 2FA |

Run every local evidence command at once.

```sh
bin/check-local
```

It runs `./scripts/bootstrap`, `reuse lint`, `./scripts/lint`, `shellcheck`, `./scripts/test-scripts`,
`./scripts/test`, `./scripts/audit`, 2 GoReleaser snapshot builds, fuzzing & `./scripts/codeql`.

## Branch Coverage Evidence

The branch gate pins [`gocove`][gocove], rebuilds every shipped package, and
validates observations with native Go coverage. Only verifier tests are excluded.

## Outstanding Gold Blockers

Repository automation cannot supply these human or organizational controls.

| Gold Criterion | Current Evidence | Required Action |
| --- | --- | --- |
| Access continuity | No evidence proves two people hold required release access | Grant and verify access for another maintainer |
| Bus factor | Git history shows one significant contributor | Add another significant contributor |
| Unassociated contributors | Fewer than two contributors are organization-independent | Accept qualifying external contributions |
| Independent modification review | Historical changes do not prove the required 50% threshold | Require another person to review qualifying pull requests |
| Human security review | No completed independent review exists within five years | Commission and publish a scoped human review |

The repository must not claim Gold while any mandatory criterion remains unmet.

## Maintenance

`bin/check-local` runs licensing, tests, coverage, security, reproducibility, fuzzing & CodeQL checks.
Run it before each merge. `bin/release-local` runs it again before each release.
The repository runs no GitHub Actions workflows or Dependabot updates.
Reassess this file before every major release.
Update bestpractices.dev only with public, verifiable evidence.

[badge]: https://www.bestpractices.dev/projects/13732
[gocove]: https://pkg.go.dev/github.com/srvgit/gocove

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
