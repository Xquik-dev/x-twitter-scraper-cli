// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package releasearchive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var zipOwnerFields = map[uint16]string{0x7855: "Unix", 0x7875: "Unix owner"}

type artifact struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func ArchiveOwners(manifestPath string) ([]string, int, error) {
	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, 0, err
	}
	var artifacts []artifact
	if err := json.Unmarshal(contents, &artifacts); err != nil {
		return nil, 0, fmt.Errorf("read %s: %w", manifestPath, err)
	}
	root := filepath.Dir(filepath.Dir(manifestPath))
	var problems []string
	checked := 0
	for _, item := range artifacts {
		if item.Type != "Archive" {
			continue
		}
		path := item.Path
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		found, err := owners(path)
		if err != nil {
			return nil, 0, fmt.Errorf("read %s: %w", item.Name, err)
		}
		for _, problem := range found {
			problems = append(problems, item.Name+" "+problem)
		}
		checked++
	}
	return problems, checked, nil
}

func owners(path string) ([]string, error) {
	switch {
	case strings.HasSuffix(path, ".tar.gz"):
		return tarOwners(path)
	case strings.HasSuffix(path, ".zip"):
		return zipOwners(path)
	default:
		return nil, errors.New("the archive is neither .tar.gz nor .zip")
	}
}

func tarOwners(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	unzipped, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	reader := tar.NewReader(unzipped)
	var problems []string
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			return problems, nil
		}
		if err != nil {
			return nil, err
		}
		if header.Uid != 0 || header.Gid != 0 || header.Uname != "root" || header.Gname != "root" {
			problems = append(problems, fmt.Sprintf(
				"%s belongs to %s (%d) and group %s (%d)",
				header.Name, header.Uname, header.Uid, header.Gname, header.Gid,
			))
		}
	}
}

func zipOwners(path string) ([]string, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var problems []string
	for _, file := range reader.File {
		extra := file.Extra
		for len(extra) >= 4 {
			id := binary.LittleEndian.Uint16(extra)
			size := int(binary.LittleEndian.Uint16(extra[2:]))
			if name, ok := zipOwnerFields[id]; ok {
				problems = append(problems, fmt.Sprintf("%s holds a %s extra field", file.Name, name))
			}
			if len(extra) < 4+size {
				break
			}
			extra = extra[4+size:]
		}
	}
	return problems, nil
}
