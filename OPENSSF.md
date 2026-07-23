# OpenSSF Best Practices Evidence

This file tracks the Gold assessment for this repository.

The official project is [bestpractices.dev project 13732][badge].

Assessment date: 2026-07-23.

## Eligibility

This public command-line client is active and released.

It is eligible for the OpenSSF Best Practices badge.

No OpenSSF-defined ineligibility applies.

## Verified Technical Controls

| Area | Evidence |
| --- | --- |
| License | Apache-2.0 with REUSE 3.3 metadata |
| Contribution process | `CONTRIBUTING.md`, DCO sign-off, and independent review policy |
| Security reporting | Private GitHub reporting, response targets, and safe harbor in `SECURITY.md` |
| Build | Pinned Go modules, checksum verification, and `./scripts/build` |
| Tests | Race-enabled tests and regression coverage through `./scripts/test` |
| Statement coverage | `./scripts/coverage` enforces 90%; current result is 91.2% |
| Branch coverage | `./scripts/branch-coverage` enforces 80%; current result is 873/1,090 (80.09%) |
| Static analysis | `go vet` and CodeQL security-extended queries |
| Dynamic analysis | Scheduled native Go fuzzing and race detection |
| Dependency review | Dependabot, `go mod verify`, and `govulncheck` |
| CI | Same-repository and fork pull requests run all required jobs |
| Releases | Signed checksums, exact tags, and GitHub artifact attestations |
| Reproducibility | CI builds each GoReleaser snapshot twice and compares checksums |
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

It rebuilds and instruments every shipped Go package.

Native Go coverage confirms which target ranges executed.

The local verifier corrects the beta reporter's nested-block associations.

It excludes only the verifier package from the shipped-code denominator.

The current result is 873 of 1,090 branches, or 80.09%.

## Outstanding Gold Blockers

These criteria need verifiable human or organizational evidence.

They cannot be satisfied through repository automation alone.

| Gold Criterion | Current Evidence | Required Action |
| --- | --- | --- |
| Access continuity | No evidence proves two people hold required release access | Grant and verify access for another maintainer |
| Bus factor | Git history shows one significant contributor | Add another significant contributor |
| Unassociated contributors | Fewer than two contributors are organization-independent | Accept qualifying external contributions |
| Independent modification review | Historical changes do not prove the required 50% threshold | Require another person to review qualifying pull requests |
| Human security review | No completed independent review exists within five years | Commission and publish a scoped human review |

The repository must not claim Gold while any mandatory criterion remains unmet.

The current remediation pull request also requires a different human reviewer.

## Maintenance

CI runs licensing, tests, coverage, security, and reproducibility checks.

Scheduled CodeQL, fuzzing, Scorecard, and dependency updates detect drift.

Reassess this file before every major release.

Update bestpractices.dev only with public, verifiable evidence.

[badge]: https://www.bestpractices.dev/projects/13732
[gocove]: https://pkg.go.dev/github.com/srvgit/gocove

Xquik is an independent third-party service. Not affiliated with X Corp. "Twitter" and "X" are trademarks of X Corp.
