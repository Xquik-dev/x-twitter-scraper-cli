package autocomplete

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func runCompletionScript(t *testing.T, shell string, writer io.Writer) error {
	t.Helper()
	command := &cli.Command{
		Name:      "xquik",
		Writer:    writer,
		ErrWriter: io.Discard,
		ExitErrHandler: func(context.Context, *cli.Command, error) {
		},
		Commands: []*cli.Command{
			{Name: "completion", Action: OutputCompletionScript},
		},
	}
	args := []string{"xquik", "completion"}
	if shell != "" {
		args = append(args, shell)
	}
	return command.Run(context.Background(), args)
}

func TestOutputCompletionScript(t *testing.T) {
	require.Error(t, runCompletionScript(t, "", io.Discard))
	require.Error(t, runCompletionScript(t, "unknown", io.Discard))

	for _, shell := range []string{"bash", "fish", "pwsh", "zsh"} {
		var output bytes.Buffer
		require.NoError(t, runCompletionScript(t, shell, &output))
		assert.Contains(t, output.String(), "xquik")
	}

	require.Error(t, runCompletionScript(t, "bash", failingWriter{}))

	original := shellCompletions["bash"]
	shellCompletions["bash"] = func(*cli.Command, string) (string, error) {
		return "", errors.New("render failed")
	}
	t.Cleanup(func() {
		shellCompletions["bash"] = original
	})
	require.Error(t, runCompletionScript(t, "bash", io.Discard))
}

func runDynamicCompletion(t *testing.T, style string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("COMPLETION_STYLE", style)
	var output bytes.Buffer
	command := &cli.Command{
		Name:      "xquik",
		Writer:    &output,
		ErrWriter: io.Discard,
		ExitErrHandler: func(context.Context, *cli.Command, error) {
		},
		Commands: []*cli.Command{
			{
				Name:    "tweets:search",
				Usage:   "Search tweets",
				Aliases: []string{"search"},
			},
			{
				Name:   "__complete",
				Action: ExecuteShellCompletion,
			},
		},
	}
	runArgs := append([]string{"xquik", "__complete"}, args...)
	err := command.Run(context.Background(), runArgs)
	return output.String(), err
}

func TestExecuteShellCompletion(t *testing.T) {
	t.Setenv("COMPLETION_STYLE", "")
	_, err := runDynamicCompletion(t, "", "tw")
	require.Error(t, err)

	_, err = runDynamicCompletion(t, "invalid", "tw")
	require.Error(t, err)

	for _, style := range []string{"bash", "zsh", "pwsh", "fish"} {
		output, err := runDynamicCompletion(t, style, "tweets", ":", "s")
		require.Error(t, err)
		assert.Contains(t, output, "search")
	}
}

func TestRebuildColonSeparatedArgs(t *testing.T) {
	t.Parallel()

	assert.Nil(t, rebuildColonSeparatedArgs(nil))
	assert.Equal(t, []string{"a", "b:c", "d"}, rebuildColonSeparatedArgs([]string{"a", "b", ":", "c", "d"}))
	assert.Equal(t, []string{"a::b"}, rebuildColonSeparatedArgs([]string{"a", ":", ":", "b"}))
	assert.Equal(t, []string{"a:", "b"}, rebuildColonSeparatedArgs([]string{"a:", "b"}))
}

func TestFindVisibilityAndMissingFlags(t *testing.T) {
	t.Parallel()

	hidden := &cli.StringFlag{Name: "hidden", Hidden: true}
	visible := &cli.StringFlag{Name: "visible"}
	root := &cli.Command{
		Flags: []cli.Flag{hidden, visible},
		Commands: []*cli.Command{
			{Name: "shown"},
			{Name: "hidden", Hidden: true},
		},
	}
	assert.Nil(t, findFlag(root, "--hidden"))
	assert.Nil(t, findFlag(root, "--missing"))
	assert.NotNil(t, findFlag(root, "--visible"))
	assert.Nil(t, findChild(root, "hidden"))
	assert.Nil(t, findChild(root, "missing"))
	assert.NotNil(t, findChild(root, "shown"))
}
