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
	testdataPath := filepath.Join("testdata", "fs")
	testCases := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "default",
			args: []string{testdataPath},
			want: "2906B\t" + testdataPath,
		},
		{
			name: "help",
			args: []string{"--help"},
			want: getHelp(t),
		},
		{
			name: "recursive",
			args: []string{testdataPath, "-r"},
			want: "5320B\t" + testdataPath,
		},
		{
			name: "with hidden",
			args: []string{filepath.Join(testdataPath, ".hidden"), "-a"},
			want: "101B\t" + filepath.Join(testdataPath, ".hidden"),
		},
		{
			name: "human",
			args: []string{filepath.Join(testdataPath, "file1.txt"), "-H"},
			want: "2.8KB\t" + filepath.Join(testdataPath, "file1.txt"),
		},
		{
			name: "recursive with hidden",
			args: []string{testdataPath, "-r", "-a"},
			want: "5523B\t" + testdataPath,
		},
		{
			name: "recursive with hidden and human",
			args: []string{testdataPath, "-r", "-a", "-H"},
			want: "5.4KB\t" + testdataPath,
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
			stdOut: getHelp(t),
		},
		{
			name:   "to many arguments",
			args:   []string{"testdata", "testdata2"},
			want:   "incorrect usage: expected 1 argument, got 2",
			stdOut: getHelp(t),
		},
		{
			name:   "incorrect flag usage",
			args:   []string{"-H=g"},
			want:   "incorrect usage: invalid value \"g\" for flag -H: parse error",
			stdOut: getHelp(t),
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

func getHelp(t *testing.T) string {
	t.Helper()
	helpPath := filepath.Join("testdata", "fixtures", "help.txt")
	helpContent, err := os.ReadFile(filepath.Clean(helpPath))
	require.NoError(t, err)

	return strings.TrimSpace(string(helpContent))
}
