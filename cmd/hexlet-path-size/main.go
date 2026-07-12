package main

import (
	"code/internal/cli"
	"context"
	"fmt"
	"os"
)

func main() {
	cmd := cli.NewCommand()
	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintln(cmd.ErrWriter, err.Error())
		os.Exit(1)
	}
}
