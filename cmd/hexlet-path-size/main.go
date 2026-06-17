package main

import (
	"code"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

type ArgumentsCountError struct {
	Expected int
	Got      int
}

func (e *ArgumentsCountError) Error() string {
	return fmt.Sprintf("incorrect usage: expected %d argument, got %d", e.Expected, e.Got)
}

type FlagError struct {
	Err error
}

func (e *FlagError) Error() string {
	return fmt.Sprintf("incorrect usage: %s", e.Err)
}

func isCommandUsageError(err error) bool {
	var argumentsCountError *ArgumentsCountError
	var flagError *FlagError

	return errors.As(err, &argumentsCountError) || errors.As(err, &flagError)
}

func main() {
	cmd := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)",
		Action: func(_ context.Context, cmd *cli.Command) error {
			const expectedArgsCount = 1
			if cmd.Args().Len() != expectedArgsCount {
				return &ArgumentsCountError{
					Expected: expectedArgsCount,
					Got:      cmd.Args().Len(),
				}
			}

			path := cmd.Args().Get(0)
			fileSize, err := code.GetPathSize(
				path,
				cmd.Bool("recursive"),
				cmd.Bool("human"),
				cmd.Bool("all"),
			)
			if err != nil {
				return err
			}

			result := fmt.Sprintf("%s\t%s", fileSize, path)
			fmt.Println(result)
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
- Hidden files/dirs: ignored (returns 0B) unless the --all (-a) flag is provided. This rule applies even if a hidden file is passed directly as an argument.
- Special files (sockets, devices, pipes): ignored (returns 0B).
- Hard links: evaluated as regular files. No deduplication is performed during recursive directory traversal.`,
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			return &FlagError{Err: err}
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		isUsageError := isCommandUsageError(err)
		if isUsageError {
			showHelpError := cli.ShowAppHelp(cmd)
			if showHelpError != nil {
				fmt.Println("Run 'make help' for usage instructions.")
			}
		}

		var userErr *code.UserError
		if errors.As(err, &userErr) {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		} else {
			if isUsageError {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "system error: %v\n", err)
			}
		}

		os.Exit(1)
	}
}
