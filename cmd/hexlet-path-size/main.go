package main

import (
	"fmt"
	"github.com/urfave/cli/v3"
	"context"
	"os"
	"log"
	"code"
)

func main() {
	cmd := &cli.Command{
		Name:  "hexlet-path-size",
        Usage: "print size of a file or directory",
        Action: func(ctx context.Context, cmd *cli.Command) error {
            if cmd.Args().Len() != 1 {
                return cli.Exit(fmt.Sprintf("incorrect usage: expected 1 argument, got %d", cmd.Args().Len()), 1)
            }

			path := cmd.Args().Get(0)
            fileSize, error := code.GetPathSize(path, cmd.Bool("recursive"), cmd.Bool("human"), cmd.Bool("all"))
            if error != nil {
                return  error
            }

            result := fmt.Sprintf("%s\t%s", fileSize, path)
            fmt.Println(result)
            return nil
        },
        Flags: []cli.Flag{
            &cli.BoolFlag{
                Name:  "recursive",
                Aliases: []string{"r"},
                Usage: "include hidden files and directories",
                DefaultText: "false",
            },
            &cli.BoolFlag{
                Name:  "human",
                Aliases: []string{"H"},
                Usage: "human-readable sizes (auto-select unit)",
                DefaultText: "false",
            },
            &cli.BoolFlag{
                Name:  "all",
                Aliases: []string{"a"},
                Usage: "include hidden files and directories",
                DefaultText: "false",
            },
        },
        ArgsUsage: "<path>",
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}