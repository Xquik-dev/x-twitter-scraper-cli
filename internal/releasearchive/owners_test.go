// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package releasearchive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type entry struct {
	name, owner, group string
	uid, gid           int
}

// release writes dist/artifacts.json and the named archives under a new
// folder, the way GoReleaser lays them out, and returns the manifest path.
func release(t *testing.T, archives map[string]func(path string)) string {
	t.Helper()
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "dist"), 0o755))
	manifest := []map[string]string{{"name": "x-twitter-scraper", "path": "dist/linux/x-twitter-scraper", "type": "Binary"}}
	for name, write := range archives {
		write(filepath.Join(root, "dist", name))
		manifest = append(manifest, map[string]string{"name": name, "path": "dist/" + name, "type": "Archive"})
	}
	contents, err := json.Marshal(manifest)
	require.NoError(t, err)
	path := filepath.Join(root, "dist", "artifacts.json")
	require.NoError(t, os.WriteFile(path, contents, 0o644))
	return path
}

func tarGz(t *testing.T, entries ...entry) func(string) {
	return func(path string) {
		file, err := os.Create(path)
		require.NoError(t, err)
		zipped := gzip.NewWriter(file)
		writer := tar.NewWriter(zipped)
		for _, e := range entries {
			require.NoError(t, writer.WriteHeader(&tar.Header{
				Name: e.name, Mode: 0o644, Size: 2,
				Uname: e.owner, Gname: e.group, Uid: e.uid, Gid: e.gid,
			}))
			_, err := writer.Write([]byte("ok"))
			require.NoError(t, err)
		}
		require.NoError(t, writer.Close())
		require.NoError(t, zipped.Close())
		require.NoError(t, file.Close())
	}
}

func zipFile(t *testing.T, extra []byte) func(string) {
	return func(path string) {
		file, err := os.Create(path)
		require.NoError(t, err)
		writer := zip.NewWriter(file)
		_, err = writer.CreateHeader(&zip.FileHeader{Name: "LICENSE", Extra: extra})
		require.NoError(t, err)
		require.NoError(t, writer.Close())
		require.NoError(t, file.Close())
	}
}

func TestArchivesOwnedByRootPass(t *testing.T) {
	manifest := release(t, map[string]func(string){
		"cli_linux_amd64.tar.gz": tarGz(t, entry{"LICENSE", "root", "root", 0, 0}, entry{"x-twitter-scraper", "root", "root", 0, 0}),
		"cli_macos_arm64.zip":    zipFile(t, nil),
	})

	problems, checked, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Empty(t, problems)
	assert.Equal(t, 2, checked)
}

func TestATarEntryOwnedByALocalUserIsReported(t *testing.T) {
	manifest := release(t, map[string]func(string){
		"cli_linux_amd64.tar.gz": tarGz(t,
			entry{"LICENSE", "burak", "staff", 501, 20},
			entry{"x-twitter-scraper", "root", "staff", 0, 20},
			entry{"man/man1/x-twitter-scraper.1.gz", "root", "root", 0, 0},
		),
	})

	problems, checked, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Equal(t, []string{
		"cli_linux_amd64.tar.gz LICENSE belongs to burak (501) and group staff (20)",
		"cli_linux_amd64.tar.gz x-twitter-scraper belongs to root (0) and group staff (20)",
	}, problems)
	assert.Equal(t, 1, checked)
}

func TestATarEntryWithRootNamesButOtherIDsIsReported(t *testing.T) {
	manifest := release(t, map[string]func(string){
		"cli.tar.gz": tarGz(t, entry{"LICENSE", "root", "root", 501, 0}),
	})

	problems, _, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Equal(t, []string{"cli.tar.gz LICENSE belongs to root (501) and group root (0)"}, problems)
}

func TestAZipEntryWithAUnixOwnerFieldIsReported(t *testing.T) {
	// A 0x7875 field holding uid 501 and gid 20, after a 0x5455 time field.
	timeField := []byte{0x55, 0x54, 0x05, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
	ownerField := []byte{0x75, 0x78, 0x0b, 0x00, 0x01, 0x04, 0xf5, 0x01, 0x00, 0x00, 0x04, 0x14, 0x00, 0x00, 0x00}
	manifest := release(t, map[string]func(string){
		"cli_windows_amd64.zip": zipFile(t, append(timeField, ownerField...)),
	})

	problems, _, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Equal(t, []string{"cli_windows_amd64.zip LICENSE holds a Unix owner extra field"}, problems)
}

func TestAZipExtraFieldThatOverrunsItsEntryEndsTheScan(t *testing.T) {
	manifest := release(t, map[string]func(string){
		"cli.zip": zipFile(t, []byte{0x55, 0x54, 0x09, 0x00, 0x01}),
	})

	problems, _, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Empty(t, problems)
}

func TestUnreadableReleasesReturnAnError(t *testing.T) {
	cases := map[string]func(t *testing.T) string{
		"missing manifest": func(t *testing.T) string {
			return filepath.Join(t.TempDir(), "dist", "artifacts.json")
		},
		"manifest that is not JSON": func(t *testing.T) string {
			path := filepath.Join(t.TempDir(), "artifacts.json")
			require.NoError(t, os.WriteFile(path, []byte("{"), 0o644))
			return path
		},
		"archive of another format": func(t *testing.T) string {
			return release(t, map[string]func(string){"cli.tar.xz": func(path string) {
				require.NoError(t, os.WriteFile(path, []byte("xz"), 0o644))
			}})
		},
		"missing archive": func(t *testing.T) string {
			return release(t, map[string]func(string){"cli.tar.gz": func(string) {}})
		},
		"tar.gz that is not gzip": func(t *testing.T) string {
			return release(t, map[string]func(string){"cli.tar.gz": func(path string) {
				require.NoError(t, os.WriteFile(path, []byte("plain"), 0o644))
			}})
		},
		"gzip that holds no tar": func(t *testing.T) string {
			return release(t, map[string]func(string){"cli.tar.gz": func(path string) {
				file, err := os.Create(path)
				require.NoError(t, err)
				zipped := gzip.NewWriter(file)
				_, err = zipped.Write([]byte("not a tar header, and too short for one"))
				require.NoError(t, err)
				require.NoError(t, zipped.Close())
				require.NoError(t, file.Close())
			}})
		},
		"zip that is not zip": func(t *testing.T) string {
			return release(t, map[string]func(string){"cli.zip": func(path string) {
				require.NoError(t, os.WriteFile(path, []byte("plain"), 0o644))
			}})
		},
	}
	for name, setup := range cases {
		t.Run(name, func(t *testing.T) {
			problems, checked, err := ArchiveOwners(setup(t))

			require.Error(t, err)
			assert.Nil(t, problems)
			assert.Zero(t, checked)
		})
	}
}

func TestAnAbsoluteArchivePathIsReadAsIs(t *testing.T) {
	archive := filepath.Join(t.TempDir(), "cli.tar.gz")
	tarGz(t, entry{"LICENSE", "burak", "staff", 501, 20})(archive)
	manifest := filepath.Join(t.TempDir(), "artifacts.json")
	contents, err := json.Marshal([]map[string]string{{"name": "cli.tar.gz", "path": archive, "type": "Archive"}})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(manifest, contents, 0o644))

	problems, checked, err := ArchiveOwners(manifest)

	require.NoError(t, err)
	assert.Equal(t, []string{"cli.tar.gz LICENSE belongs to burak (501) and group staff (20)"}, problems)
	assert.Equal(t, 1, checked)
}
