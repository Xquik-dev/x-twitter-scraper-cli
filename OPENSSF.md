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
| Statement coverage | `./scripts/coverage` enforces 90%; current result is 91.1% |
| Branch coverage | `./scripts/branch-coverage` enforces 80%; current result is 881/1,100 (80.09%) |
| Static analysis | `go vet` & CodeQL security-extended queries |
| Dynamic analysis | Scheduled native Go fuzzing & race detection |
| Dependency review | Dependabot, `go mod verify`, & `govulncheck` |
| CI | Same-repository & fork pull requests run all required jobs |
| Releases | Signed checksums, exact tags, & GitHub artifact attestations |
| Reproducibility | CI builds each GoReleaser snapshot twice & compares checksums |
| Portability | CI compiles every package for Windows AMD64 |
| Two-factor authentication | The Xquik-dev organization requires 2FA |

Run the local evidence commands:

```sh
./scripts/bootstrap
./scripts/lint
./scripts/test
uvx --from reuse==5.1.1 reuse lint
go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...
```

## Branch Coverage Evidence

The branch gate pins an exact FLOSS [`gocove`][gocove] version.
It rebuilds & instruments every shipped Go package.
Native Go coverage confirms each executed range.
The verifier corrects nested-block associations from the beta reporter.
Only the verifier package is excluded from shipped-code coverage.

The current result is 881 of 1,100 branches, or 80.09%.

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

Gold eligibility still requires review by a different human.

## Maintenance

CI runs licensing, tests, coverage, security, & reproducibility checks.
Scheduled CodeQL, fuzzing, Scorecard, & dependency updates detect drift.
Reassess this file before every major release.
Update bestpractices.dev only with public, verifiable evidence.

[badge]: https://www.bestpractices.dev/projects/13732
[gocove]: https://pkg.go.dev/github.com/srvgit/gocove

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
