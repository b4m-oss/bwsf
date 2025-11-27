# bwenv

bwenv is a CLI tool that uses [Bitwarden](https://bitwarden.com/) to manage .env files.

## Overview

bwenv commands supports your dotenv files are manged in your Bitwarden.

Simple usage below:

| command | |
|----|----|
| bwenv push | .env files push to your Bitwarden host |
| bwenv pull | .env files pull from your Bitwarden host |
| bwenv list | Show list stored .env files at your Bitwarden host |

## Motivation

We use Bitwarden as our password manager long time ago.
Also, we store .env files our Bitwarden host, manage them as shell scripts.
This project migrates our hand-maded shell scripts to modern CLI command with Go.

## Requirements

**`bw`** command is needed to be installed your machine.

[To install bw command, please read this document.](https://bitwarden.com/help/cli/#download-and-install)

### Machine OS

- [Is planning] macOS
- [Is planning] Linux
- [Is planning] Windows

## Installation

[!Note]
This is under planning.

- **macOS & Linux**: Homebrew
- **Windows**: Chocolaty

## Development

### Requirement

**Docker** is needed to be installed your development machine.

### Start up to dev

```
git clone https://github.com/b4m-oss/bwenv.git
cd bwenv
make run
```

## License

[MIT License](./LICENSE)