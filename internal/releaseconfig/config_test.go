// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package releaseconfig

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGoReleaserUsesExistingDraft(t *testing.T) {
	t.Parallel()

	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate release config test")
	}
	configPath := filepath.Join(filepath.Dir(sourceFile), "..", "..", ".goreleaser.yml")
	contents, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read GoReleaser config: %v", err)
	}

	usesDraft, err := usesExistingDraft(contents)
	if err != nil {
		t.Fatalf("parse GoReleaser config: %v", err)
	}
	if !usesDraft {
		t.Fatal("release.use_existing_draft must be true")
	}
}

func TestGoReleaserDraftParsingRejectsMisleadingText(t *testing.T) {
	t.Parallel()

	for name, contents := range map[string]string{
		"comment":   "# use_existing_draft: true\nrelease:\n  use_existing_draft: false\n",
		"unrelated": "other:\n  use_existing_draft: true\nrelease:\n  use_existing_draft: false\n",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			usesDraft, err := usesExistingDraft([]byte(contents))
			if err != nil {
				t.Fatalf("parse config: %v", err)
			}
			if usesDraft {
				t.Fatal("misleading text must not enable existing draft reuse")
			}
		})
	}
}

func TestGoReleaserDraftParsingRejectsInvalidYaml(t *testing.T) {
	t.Parallel()

	if _, err := usesExistingDraft([]byte("release: [")); err == nil {
		t.Fatal("invalid YAML must fail")
	}
}
