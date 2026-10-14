// SPDX-FileCopyrightText: 2026 Xquik contributors
//
// SPDX-License-Identifier: Apache-2.0

package autocomplete

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestGetCompletions_EmptyArgs(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "generate", Usage: "Generate SDK"},
			{Name: "test", Usage: "Run tests"},
			{Name: "build", Usage: "Build project"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 3)
	assert.Contains(t, result.Completions, ShellCompletion{Name: "generate", Usage: "Generate SDK"})
	assert.Contains(t, result.Completions, ShellCompletion{Name: "test", Usage: "Run tests"})
	assert.Contains(t, result.Completions, ShellCompletion{Name: "build", Usage: "Build project"})
}

func TestGetCompletions_SubcommandPrefix(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "generate", Usage: "Generate SDK"},
			{Name: "test", Usage: "Run tests"},
			{Name: "build", Usage: "Build project"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"ge"})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "generate", result.Completions[0].Name)
	assert.Equal(t, "Generate SDK", result.Completions[0].Usage)
}

func TestGetCompletions_HiddenCommand(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "visible", Usage: "Visible command"},
			{Name: "hidden", Usage: "Hidden command", Hidden: true},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{""})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "visible", result.Completions[0].Name)
}

func TestGetCompletions_NestedSubcommand(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "config",
				Usage: "Configuration commands",
				Commands: []*cli.Command{
					{Name: "get", Usage: "Get config value"},
					{Name: "set", Usage: "Set config value"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"config", "s"})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "set", result.Completions[0].Name)
	assert.Equal(t, "Set config value", result.Completions[0].Usage)
}

func TestGetCompletions_FlagCompletion(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output directory"},
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Verbose output"},
					&cli.StringFlag{Name: "format", Usage: "Output format"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "--o"})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "--output", result.Completions[0].Name)
	assert.Equal(t, "Output directory", result.Completions[0].Usage)
}

func TestGetCompletions_ShortFlagCompletion(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Output directory"},
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Verbose output"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "-v"})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "-v", result.Completions[0].Name)
}

func TestGetCompletions_BoolFlagDoesNotBlockCompletion(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}, Usage: "Verbose output"},
				},
				Commands: []*cli.Command{
					{Name: "typescript", Usage: "Generate TypeScript SDK"},
					{Name: "python", Usage: "Generate Python SDK"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "--verbose", "ty"})

	assert.Equal(t, ShellCompletionBehaviorDefault, result.Behavior)
	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "typescript", result.Completions[0].Name)
}

func TestGetCompletions_ColonCommands_NoColonTyped(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
			{Name: "config:list", Usage: "List config values"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"co"})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "config", result.Completions[0].Name)
	assert.Equal(t, "", result.Completions[0].Usage)
}

func TestGetCompletions_ColonCommands_ColonTyped_Bash(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
			{Name: "config:list", Usage: "List config values"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"config:"})

	assert.Len(t, result.Completions, 3)
	names := []string{result.Completions[0].Name, result.Completions[1].Name, result.Completions[2].Name}
	assert.Contains(t, names, "get")
	assert.Contains(t, names, "set")
	assert.Contains(t, names, "list")
}

func TestGetCompletions_ColonCommands_ColonTyped_Zsh(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
			{Name: "config:list", Usage: "List config values"},
		},
	}

	result := GetCompletions(CompletionStyleZsh, root, []string{"config:"})

	assert.Len(t, result.Completions, 3)
	names := []string{result.Completions[0].Name, result.Completions[1].Name, result.Completions[2].Name}
	assert.Contains(t, names, "config:get")
	assert.Contains(t, names, "config:set")
	assert.Contains(t, names, "config:list")
}

func TestGetCompletions_BashStyleColonCompletion(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"config:g"})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "get", result.Completions[0].Name)
	assert.Equal(t, "Get config value", result.Completions[0].Usage)
}

func TestGetCompletions_BashStyleColonCompletion_NoMatch(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"other:g"})

	assert.Len(t, result.Completions, 0)
}

func TestGetCompletions_ZshStyleColonCompletion(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
		},
	}

	result := GetCompletions(CompletionStyleZsh, root, []string{"config:g"})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "config:get", result.Completions[0].Name)
	assert.Equal(t, "Get config value", result.Completions[0].Usage)
}

func TestGetCompletions_MixedColonAndRegularCommands(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "generate", Usage: "Generate SDK"},
			{Name: "config:get", Usage: "Get config value"},
			{Name: "config:set", Usage: "Set config value"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{""})

	assert.Len(t, result.Completions, 2)
	names := []string{result.Completions[0].Name, result.Completions[1].Name}
	assert.Contains(t, names, "generate")
	assert.Contains(t, names, "config")
}

func TestGetCompletions_FlagWithBoolFlagSkipsValue(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}},
				},
				Commands: []*cli.Command{
					{Name: "typescript", Usage: "TypeScript SDK"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "-v", "ty"})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "typescript", result.Completions[0].Name)
}

func TestGetCompletions_MultipleFlagsBeforeSubcommand(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
				},
				Commands: []*cli.Command{
					{Name: "typescript", Usage: "TypeScript SDK"},
					{Name: "python", Usage: "Python SDK"},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "-c", "config.yml", "-v", "py"})

	assert.Len(t, result.Completions, 1)
	assert.Equal(t, "python", result.Completions[0].Name)
}

func TestGetCompletions_CommandAliases(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{Name: "generate", Aliases: []string{"gen", "g"}, Usage: "Generate SDK"},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"g"})

	assert.GreaterOrEqual(t, len(result.Completions), 2)
	names := []string{}
	for _, c := range result.Completions {
		names = append(names, c.Name)
	}
	assert.Contains(t, names, "generate")
	assert.Contains(t, names, "gen")
}

func TestGetCompletions_AllFlagsWhenNoPrefix(t *testing.T) {
	t.Parallel()

	root := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SDK",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}},
					&cli.BoolFlag{Name: "verbose", Aliases: []string{"v"}},
					&cli.StringFlag{Name: "format", Aliases: []string{"f"}},
				},
			},
		},
	}

	result := GetCompletions(CompletionStyleBash, root, []string{"generate", "-"})

	assert.GreaterOrEqual(t, len(result.Completions), 6)
}
