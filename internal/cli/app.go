package cli

import (
	"code"
	"code/internal/formatter"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

// NewCommand creates a new command.
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		Action: func(_ context.Context, cmd *cli.Command) error {
			const expectedArgsCount = 1
			if cmd.Args().Len() != expectedArgsCount {
				return newArgumentsCountError(cmd, expectedArgsCount, cmd.Args().Len())
			}

			path := cmd.Args().Get(0)
			fileSize, err := code.GetPathSize(
				path,
				cmd.Bool("recursive"),
				cmd.Bool("human"),
				cmd.Bool("all"),
			)
			if err != nil {
				return newRuntimeError(err)
			}

			_, err = fmt.Fprintln(cmd.Writer, formatter.FormatOutput(fileSize, path))
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to write output: %v", err), 1)
			}

			return nil
		},
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			cli.ShowAppHelp(cmd)

			return cli.Exit(fmt.Sprintf("incorrect usage: %v", err), 1)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "recursive size of directories",
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "human",
				Aliases:     []string{"H"},
				Usage:       "human-readable sizes (auto-select unit)",
				DefaultText: "false",
			},
			&cli.BoolFlag{
				Name:        "all",
				Aliases:     []string{"a"},
				Usage:       "include hidden files and directories",
				DefaultText: "false",
			},
		},
		ArgsUsage: "<path>",
		Description: `Calculate and print the size of a file or directory.

CALCULATION RULES:
- Regular files: evaluated by their actual size.
- Symlinks: evaluated by the size of the link itself, not the target file.
- Hidden files/dirs nested during traversal: ignored (returns 0B) unless the --all (-a) flag is provided. A hidden path passed directly as an argument is always evaluated.
- Special files (sockets, devices, pipes): ignored (returns 0B).
- Hard links: evaluated as regular files. No deduplication is performed during recursive directory traversal.`,
	}
}

func newArgumentsCountError(cmd *cli.Command, expected, got int) error {
	cli.ShowAppHelp(cmd)

	return cli.Exit(
		fmt.Sprintf("incorrect usage: expected %d argument, got %d", expected, got),
		1,
	)
}

func newRuntimeError(err error) error {
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		if errors.Is(pathErr.Err, os.ErrNotExist) {
			return cli.Exit("path does not exist: "+pathErr.Path, 1)
		}

		if errors.Is(pathErr.Err, os.ErrPermission) {
			return cli.Exit("permission denied: "+pathErr.Path, 1)
		}
	}

	return fmt.Errorf("system error: %w", err)
}
