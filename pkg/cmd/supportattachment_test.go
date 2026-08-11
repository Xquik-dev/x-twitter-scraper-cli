// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func TestSupportAttachmentsDownload(t *testing.T) {
	t.Run("regular flags", func(t *testing.T) {
		mocktest.TestRunMockTestWithFlags(
			t,
			"--api-key", "string",
			"--bearer-token", "string",
			"support:attachments", "download",
			"--id", "att_a1b2c3d4e5f6a1b2c3d4e5f6",
			"--range", "bytes=0-1048575",
			"--output", "/dev/null",
		)
	})
}
