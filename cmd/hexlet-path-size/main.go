package main

import (
	"code"
	"code/internal/fs"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func newArgumentsCountError(cmd *cli.Command, expected, got int) error {
	if showHelpError := cli.ShowAppHelp(cmd); showHelpError != nil {
		fmt.Println("Run 'make help' for usage instructions.")
	}

	return cli.Exit(
		fmt.Sprintf("incorrect usage: expected %d argument, got %d", expected, got),
		1,
	)
}

func newRuntimeError(err error) error {
	switch {
	case errors.Is(err, fs.ErrPathNotExist), errors.Is(err, fs.ErrPermissionDenied):
		return cli.Exit(err, 1)
	default:
		return cli.Exit(fmt.Sprintf("system error: %v", err), 1)
	}
}

func main() {
	cmd := &cli.Command{
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

			fmt.Printf("%s\t%s\n", fileSize, path)
			return nil
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
- Hidden files/dirs: ignored (returns 0B) unless the --all (-a) flag is provided. This rule applies even if a hidden path is passed directly as an argument.
- Special files (sockets, devices, pipes): ignored (returns 0B).
- Hard links: evaluated as regular files. No deduplication is performed during recursive directory traversal.`,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
