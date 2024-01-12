# Magnet

This project is a CLI tool to get magnet links from Torrent Galaxy based on a search query and optional season and episode parameters.

## Usage

```sh
magnet "For All Mankind S04E10"
magnet "For All Mankind" --season 4 --episode 10
magnet "For All Mankind" -s 4 -e 10
```

## Description

The tool searches for magnet links on Torrent Galaxy based on the provided search query and optional season and episode parameters.

## Installation

To use this tool, you can clone the repository and build the project using Go.

```sh
# clone the repository
git clone https://github.com/alissonsleal/magnet
cd magnet

# download dependencies
go mod download

# build the project
go build -o magnet
# or install the project
go install

# run the project (if built)
./magnet "For All Mankind S04E10"
# run the project (if installed)
magnet "For All Mankind S04E10"
```

## Dependencies

Developed with Go 1.21.6, but should work with older versions

We're using go modules to manage dependencies. The following dependencies are used in this project:

- github.com/andybalholm/cascadia
- github.com/atotto/clipboard
- github.com/pterm/pterm
- github.com/spf13/cobra
- golang.org/x/net/html

## Usage

To use the tool, simply run the executable and provide the search query along with optional season and episode parameters.

```sh
magnet "For All Mankind S04E10"
magnet "For All Mankind" --season 4 --episode 10
magnet "For All Mankind" -s 4 -e 10
```

## Optional Flags

The following flags are available:

- `--season` or `-s`: The season number
- `--episode` or `-e`: The episode number

## License

This project is licensed under the MIT License.

## TODO

- [ ] Add support for other torrent sites
- [ ] Release binaries for Linux
- [ ] Release binaries for macOS
- [ ] Release binaries for Windows
