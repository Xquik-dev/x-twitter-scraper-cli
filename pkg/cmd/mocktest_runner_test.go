package cmd

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/Xquik-dev/x-twitter-scraper-cli/internal/mocktest"
)

func init() {
	mocktest.RegisterRunner(func(t *testing.T, ctx context.Context, args []string, pipeData []byte) error {
		command := newCommand()
		command.Writer = io.Discard
		command.ErrWriter = io.Discard

		if pipeData == nil {
			return command.Run(ctx, args)
		}

		stdin, err := os.CreateTemp(t.TempDir(), "stdin-*")
		if err != nil {
			return err
		}
		defer stdin.Close()
		if _, err := stdin.Write(pipeData); err != nil {
			return err
		}
		if _, err := stdin.Seek(0, 0); err != nil {
			return err
		}

		originalStdin := os.Stdin
		os.Stdin = stdin
		defer func() {
			os.Stdin = originalStdin
		}()

		return command.Run(ctx, args)
	})
}
