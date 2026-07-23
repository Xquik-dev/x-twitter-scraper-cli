//go:build ignore

// Copyright the Xquik contributors.
// SPDX-License-Identifier: Apache-2.0

// branch_coverage corrects gocove's cover-profile join before enforcing a gate.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/coveragebranch"
)

func main() {
	if len(os.Args) != 4 {
		exitf("usage: branch_coverage.go REPORT PROFILE MINIMUM")
	}
	minimum, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		exitf("parse minimum: %v", err)
	}

	report, err := os.ReadFile(os.Args[1])
	if err != nil {
		exitf("read report: %v", err)
	}
	profile, err := os.Open(os.Args[2])
	if err != nil {
		exitf("open profile: %v", err)
	}
	defer profile.Close()

	covered, total, err := coveragebranch.MeasureExcluding(
		report,
		profile,
		"github.com/Xquik-dev/x-twitter-scraper-cli/internal/coveragebranch",
	)
	if err != nil {
		exitf("measure branch coverage: %v", err)
	}
	if total == 0 {
		exitf("branch report has no branches")
	}

	percent := 100 * float64(covered) / float64(total)
	fmt.Printf(
		"branch coverage: %d/%d (%.2f%%); minimum %.2f%%\n",
		covered,
		total,
		percent,
		minimum,
	)
	if percent+0.0000001 < minimum {
		os.Exit(1)
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
