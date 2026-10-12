# OpenSSF Best Practices Evidence

This file records the repository's Gold assessment.

The official project is [bestpractices.dev project 13732][badge].

Assessment date: 2026-07-23.

## Eligibility

This active public CLI qualifies for the OpenSSF Best Practices badge.
No OpenSSF exclusion applies.

## Verified Technical Controls

| Area | Evidence |
| --- | --- |
| License | Apache-2.0, REUSE 3.3 metadata, & bundled BSD-3-Clause terms |
| Contribution process | `CONTRIBUTING.md`, DCO sign-off, & independent review policy |
| Security reporting | Private GitHub reporting, response targets, & safe harbor in `SECURITY.md` |
| Build | Pinned Go modules, checksum verification, & `./scripts/build` |
| Tests | Race-enabled tests & regression coverage through `./scripts/test` |
| Statement coverage | `./scripts/coverage` enforces 90%; current result is 91.2% |
| Branch coverage | `./scripts/branch-coverage` enforces 80%; current result is 881/1,100 (80.09%) |
| Static analysis | `go vet` through `./scripts/lint` |
| Dynamic analysis | Race detection in `./scripts/test` & a native Go fuzz target that maintainers run by hand |
| Dependency review | `go mod verify` & `govulncheck` through `./scripts/audit` |
| Local checks | `bin/check-local` runs licensing, lint, test, audit, & build checks, then posts the `xquik-local-checks` status |
| Releases | `bin/release-local` publishes checksums & an annotated tag from `main` |
| Reproducibility | `bin/check-local` builds each GoReleaser snapshot twice & compares checksums |
| Portability | `./scripts/test` compiles every package for Windows AMD64 |
| Two-factor authentication | The Xquik-dev organization requires 2FA |

Run every local evidence command at once:

```sh
bin/check-local
```

It runs `./scripts/bootstrap`, `reuse lint`, `./scripts/lint`, `./scripts/test`,
`./scripts/audit`, & 2 GoReleaser snapshot builds.

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

`bin/check-local` runs licensing, tests, coverage, security, & reproducibility checks.
Run it before each merge. `bin/release-local` runs it again before each release.
The repository runs no GitHub Actions workflows or Dependabot updates.
Reassess this file before every major release.
Update bestpractices.dev only with public, verifiable evidence.

[badge]: https://www.bestpractices.dev/projects/13732
[gocove]: https://pkg.go.dev/github.com/srvgit/gocove

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
