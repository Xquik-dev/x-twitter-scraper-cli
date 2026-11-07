// Copyright the Xquik contributors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/coveragebranch"
)

func main() {
	if len(os.Args) != 4 {
		exit("Usage: branch-coverage REPORT PROFILE MINIMUM")
	}
	minimum, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		exit("parse minimum: " + err.Error())
	}

	report, err := os.ReadFile(os.Args[1])
	if err != nil {
		exit("read report: " + err.Error())
	}
	profile, err := os.Open(os.Args[2])
	if err != nil {
		exit("open profile: " + err.Error())
	}
	defer profile.Close()

	covered, total, err := coveragebranch.MeasureExcluding(
		report,
		profile,
		"github.com/Xquik-dev/x-twitter-scraper-cli/internal/coveragebranch",
	)
	if err != nil {
		exit("measure branch coverage: " + err.Error())
	}
	if total == 0 {
		exit("Branch report has no branches. Check gocove output.")
	}

	percent := 100 * float64(covered) / float64(total)
	fmt.Printf(
		"Branch coverage: %d/%d (%.2f%%); minimum %.2f%%\n",
		covered,
		total,
		percent,
		minimum,
	)
	if percent+0.0000001 < minimum {
		os.Exit(1)
	}
}

func exit(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
