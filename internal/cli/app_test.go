package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	urfave "github.com/urfave/cli/v3"
)

func TestNewCommand(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "default",
			args: []string{"testdata"},
			want: "2906B\ttestdata",
		},
		{
			name: "help",
			args: []string{"--help"},
			want: getHelp(),
		},
		{
			name: "recursive",
			args: []string{"testdata", "-r"},
			want: "5320B\ttestdata",
		},
		{
			name: "with hidden",
			args: []string{filepath.Join("testdata", ".hidden"), "-a"},
			want: "101B\ttestdata/.hidden",
		},
		{
			name: "human",
			args: []string{filepath.Join("testdata", "file1.txt"), "-H"},
			want: "2.8KB\ttestdata/file1.txt",
		},
		{
			name: "recursive with hidden",
			args: []string{"testdata", "-r", "-a"},
			want: "5523B\ttestdata",
		},
		{
			name: "recursive with hidden and human",
			args: []string{"testdata", "-r", "-a", "-H"},
			want: "5.4KB\ttestdata",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := NewCommand()
			app.ExitErrHandler = func(_ context.Context, cmd *urfave.Command, err error) {
				if err != nil {
					_, err := fmt.Fprintln(cmd.ErrWriter, err.Error())
					require.NoError(t, err)
				}
			}

			var stdout, stderr bytes.Buffer
			app.Writer = &stdout
			app.ErrWriter = &stderr

			runArgs := append([]string{"hexlet-path-size"}, tc.args...)

			err := app.Run(context.Background(), runArgs)
			assert.NoError(t, err)
			assert.Empty(t, stderr.String())
			assert.Equal(t, tc.want, strings.TrimSpace(stdout.String()))
		})
	}
}

func TestNewApp_Run_FSErrors(t *testing.T) {
	tempDir := t.TempDir()
	err := os.Mkdir(filepath.Join(tempDir, "unreadable"), 0o000)
	require.NoError(t, err)

	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0o750)
	require.NoError(t, err)

	notFoundFile := filepath.Join(tempDir, "not-found")

	testCases := []struct {
		name   string
		args   []string
		want   string
		stdOut string
	}{
		{
			name: "file not found",
			args: []string{notFoundFile},
			want: "path does not exist: " + notFoundFile,
		},
		{
			name: "unreadable directory",
			args: []string{filepath.Join(tempDir, "unreadable")},
			want: "permission denied: " + filepath.Join(tempDir, "unreadable"),
		},
		{
			name: "unreadable file recursive",
			args: []string{filepath.Join(tempDir, "unreadable"), "-r"},
			want: "permission denied: " + filepath.Join(tempDir, "unreadable"),
		},
		{
			name:   "without arguments",
			args:   []string{},
			want:   "incorrect usage: expected 1 argument, got 0",
			stdOut: getHelp(),
		},
		{
			name:   "to many arguments",
			args:   []string{"testdata", "testdata2"},
			want:   "incorrect usage: expected 1 argument, got 2",
			stdOut: getHelp(),
		},
		{
			name:   "incorrect flag usage",
			args:   []string{"-H=g"},
			want:   "incorrect usage: invalid value \"g\" for flag -H: parse error",
			stdOut: getHelp(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := NewCommand()
			app.ExitErrHandler = func(_ context.Context, cmd *urfave.Command, err error) {
				if err != nil {
					_, writeErr := fmt.Fprintln(cmd.ErrWriter, err.Error())
					require.NoError(t, writeErr)
				}
			}

			var stdout, stderr bytes.Buffer
			app.Writer = &stdout
			app.ErrWriter = &stderr

			runArgs := append([]string{"hexlet-path-size"}, tc.args...)

			err := app.Run(context.Background(), runArgs)
			require.Error(t, err)
			assert.EqualError(t, err, tc.want)

			if tc.stdOut == "" {
				assert.Empty(t, stdout.String())
			} else {
				assert.Equal(t, tc.stdOut, strings.TrimSpace(stdout.String()))
			}
		})
	}
}

func getHelp() string {
	return `NAME:
   hexlet-path-size - print size of a file or directory; supports -r (recursive), -H (human-readable), -a (include hidden)

USAGE:
   hexlet-path-size [global options] <path>

DESCRIPTION:
   Calculate and print the size of a file or directory.

   CALCULATION RULES:
   - Regular files: evaluated by their actual size.
   - Symlinks: evaluated by the size of the link itself, not the target file.
   - Hidden files/dirs nested during traversal: ignored (returns 0B) unless the --all (-a) flag is provided. A hidden path passed directly as an argument is always evaluated.
   - Special files (sockets, devices, pipes): ignored (returns 0B).
   - Hard links: evaluated as regular files. No deduplication is performed during recursive directory traversal.

GLOBAL OPTIONS:
   --recursive, -r  recursive size of directories (default: false)
   --human, -H      human-readable sizes (auto-select unit) (default: false)
   --all, -a        include hidden files and directories (default: false)
   --help, -h       show help`
}
