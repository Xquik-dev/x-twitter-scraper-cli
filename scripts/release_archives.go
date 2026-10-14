//go:build ignore

// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/releasearchive"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Run go run ./scripts/release_archives.go <dist/artifacts.json>.")
		os.Exit(2)
	}
	problems, checked, err := releasearchive.ArchiveOwners(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "The release archives could not be read. %v. Build them again, then rerun.\n", err)
		os.Exit(1)
	}
	if checked == 0 {
		fmt.Fprintln(os.Stderr, "The manifest lists no archive. Build the release assets, then rerun.")
		os.Exit(1)
	}
	for _, problem := range problems {
		fmt.Fprintln(os.Stderr, problem)
	}
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "%d archive entries name a user or group other than root. Set owner and group to root for every archive entry in .goreleaser.yml, then rerun.\n", len(problems))
		os.Exit(1)
	}
	fmt.Printf("Every entry of the %d release archives belongs to root.\n", checked)
}
