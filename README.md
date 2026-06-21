# Hexlet Path Size

[![Actions Status](https://github.com/vyacheslavkor/go-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/vyacheslavkor/go-project-242/actions) [![CI Status](https://github.com/vyacheslavkor/go-project-242/actions/workflows/ci.yml/badge.svg)](https://github.com/vyacheslavkor/go-project-242/actions/workflows/ci.yml)

A fast, standard-compliant command-line utility written in Go for calculating the size of files and directories.

## Demo
[![asciicast](https://asciinema.org/a/1188002.svg)](https://asciinema.org/a/1188002)

## Installation

You can clone the repository to any convenient local directory on your machine. There are no strict web-server path requirements.

1. Clone the repository
```bash
git clone https://github.com/vyacheslavkor/go-project-242.git
```

2. Navigate into the project directory
```bash
cd go-project-242
```

3. Build the executable (assuming make is installed)
```bash
make build
```

## Usage

The utility expects exactly one argument: the path to a file or directory. By default, it outputs the size in bytes.

```bash
./bin/hexlet-path-size [global options] <path>
```

### Options / Flags

| Flag | Alias | Default | Description |
| :--- | :--- | :--- | :--- |
| `--recursive` | `-r` | `false` | Calculate the recursive size of directories (including all subdirectories). |
| `--human` | `-H` | `false` | Display sizes in human-readable format (e.g., KB, MB, GB). Auto-selects the most appropriate unit. |
| `--all` | `-a` | `false` | Include hidden files and directories (those starting with a dot `.`) in the size calculation. |
| `--help` | `-h` | `false` | Show the help screen and usage instructions. |

### Calculation Rules

- **Regular files** are evaluated by their actual size.
- **Symlinks** are evaluated by the size of the link itself, not the target file.
- **Hidden files and directories** are ignored and contribute `0B` unless `--all` (`-a`) is provided. This also applies when the hidden path itself is passed directly as the argument.
- **Special files** (sockets, devices, pipes) are ignored and contribute `0B`.
- **Hard links** are evaluated as regular files. No deduplication is performed during recursive directory traversal.

### Output Format

The standard output strictly follows the contract: `<size>\t<path>` (separated by a tab character). The path is printed exactly as it was passed on the command line.

### Examples

**1. Basic file size calculation (raw bytes):**
```bash
$ ./bin/hexlet-path-size ./testdata/file1.txt
2906B	./testdata/file1.txt
```

**2. Calculating directory size (human-readable):**
```bash
$ ./bin/hexlet-path-size -H ./testdata/in
2.4KB	./testdata/in
```

**3. Full recursive calculation including hidden files:**
```bash
$ ./bin/hexlet-path-size -r -H -a ./testdata
5.3KB	./testdata
```
