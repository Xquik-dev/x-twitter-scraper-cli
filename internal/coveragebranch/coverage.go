// Copyright the Xquik contributors.
// SPDX-License-Identifier: Apache-2.0

// Package coveragebranch joins gocove branch targets with native Go coverage.
package coveragebranch

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type position struct {
	line int
	col  int
}

type profileEntry struct {
	start position
	end   position
}

type coverageReport struct {
	Packages []reportPackage `json:"packages"`
}

type reportPackage struct {
	ImportPath string       `json:"import_path"`
	Files      []reportFile `json:"files"`
}

type reportFile struct {
	Name     string         `json:"name"`
	Blocks   []reportBlock  `json:"blocks"`
	Branches []reportBranch `json:"branches"`
}

type reportBlock struct {
	Index   int `json:"idx"`
	Line    int `json:"line"`
	Column  int `json:"col"`
	EndLine int `json:"end_line"`
	EndCol  int `json:"end_col"`
}

type reportBranch struct {
	BlockIndex int `json:"block_idx"`
}

var profileLine = regexp.MustCompile(
	`^(.*):([0-9]+)\.([0-9]+),([0-9]+)\.([0-9]+)[[:space:]]+[0-9]+[[:space:]]+([0-9]+)$`,
)

// Measure reports covered and total branch targets.
//
// The pinned gocove beta can associate a nested branch with an outer block.
// Measure instead checks each target against positive native profile ranges.
// Strict overlap rejects ranges that only touch a target boundary.
func Measure(reportJSON []byte, profile io.Reader) (int, int, error) {
	return MeasureExcluding(reportJSON, profile)
}

// MeasureExcluding reports branch coverage outside test-only packages.
func MeasureExcluding(
	reportJSON []byte,
	profile io.Reader,
	excludedImportPaths ...string,
) (int, int, error) {
	var report coverageReport
	if err := json.Unmarshal(reportJSON, &report); err != nil {
		return 0, 0, fmt.Errorf("parse report: %w", err)
	}

	entries, err := readProfile(profile)
	if err != nil {
		return 0, 0, err
	}
	return measure(report, entries, excludedImportPaths)
}

func readProfile(profile io.Reader) (map[string][]profileEntry, error) {
	entries := make(map[string][]profileEntry)
	scanner := bufio.NewScanner(profile)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		match := profileLine.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("parse profile line %q", line)
		}

		count, err := parseInt(match[6])
		if err != nil {
			return nil, err
		}
		if count == 0 {
			continue
		}

		startLine, err := parseInt(match[2])
		if err != nil {
			return nil, err
		}
		startCol, err := parseInt(match[3])
		if err != nil {
			return nil, err
		}
		endLine, err := parseInt(match[4])
		if err != nil {
			return nil, err
		}
		endCol, err := parseInt(match[5])
		if err != nil {
			return nil, err
		}
		entry := profileEntry{
			start: position{line: startLine, col: startCol},
			end:   position{line: endLine, col: endCol},
		}
		entries[match[1]] = append(entries[match[1]], entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read profile: %w", err)
	}
	return entries, nil
}

func measure(
	report coverageReport,
	entries map[string][]profileEntry,
	excludedImportPaths []string,
) (int, int, error) {
	covered := 0
	total := 0
	for _, pkg := range report.Packages {
		if slices.Contains(excludedImportPaths, pkg.ImportPath) {
			continue
		}
		for _, file := range pkg.Files {
			blocks := make(map[int]reportBlock, len(file.Blocks))
			for _, block := range file.Blocks {
				blocks[block.Index] = block
			}
			path := pkg.ImportPath + "/" + file.Name
			for _, branch := range file.Branches {
				total++
				block, ok := blocks[branch.BlockIndex]
				if !ok {
					return 0, 0, fmt.Errorf(
						"%s references missing block %d",
						path,
						branch.BlockIndex,
					)
				}
				if blockCovered(block, entries[path]) {
					covered++
				}
			}
		}
	}
	return covered, total, nil
}

func blockCovered(block reportBlock, entries []profileEntry) bool {
	start := position{line: block.Line, col: block.Column}
	end := position{line: block.EndLine, col: block.EndCol}
	for _, entry := range entries {
		if before(entry.start, end) && before(start, entry.end) {
			return true
		}
	}
	return false
}

func before(left position, right position) bool {
	return left.line < right.line ||
		left.line == right.line && left.col < right.col
}

func parseInt(value string) (int, error) {
	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("parse profile number %q: %w", value, err)
	}
	return number, nil
}
