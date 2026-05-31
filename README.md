# Hexlet Path Size

[![Actions Status](https://github.com/vyacheslavkor/go-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/vyacheslavkor/go-project-242/actions)

A fast, standard-compliant command-line utility written in Go for calculating the size of files and directories. 

## Demo
[![asciicast](https://asciinema.org/a/1077409.svg)](https://asciinema.org/a/1077409)

## Installation

You can clone the repository to any convenient local directory on your machine. There are no strict web-server path requirements.

```bash
# Clone the repository
git clone [https://github.com/vyacheslavkor/go-project-242.git](https://github.com/vyacheslavkor/go-project-242.git)

# Navigate into the project directory
cd go-project-242

# Build the executable (assuming make is installed)
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

### Output Format

The standard output strictly follows the contract: `<size>\t<path>` (separated by a tab character).

### Examples

**1. Basic file size calculation (raw bytes):**
```bash
$ ./bin/hexlet-path-size ./testdata/file1.txt
2906B	file.txt
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
