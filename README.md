# bwenv

bwenv is a CLI tool that uses [Bitwarden](https://bitwarden.com/) to manage .env files.

## 🚨🚨 BREAKING CHANGE 🚨🚨

From v0.9.0, bwenv stores multiple enviroment .env files, like `.env | .env.staging | .env.production`.
Cause with this, stored data at Bitwarden Note item structure is changed.
Stored data before v0.8.0 is no compatiblity after v0.9.0.
We will not provide migration system.

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

- macOS
- [Is planning] Linux
- [Is planning] Windows

## Installation

| OS | command |
|----|----|
| macOS | brew tap b4m-oss/tap && brew install bwenv |

## Confirm installation

```shell
bwenv -v
# bwenv version 0.5.5
```

## Usage

### Initial setup

```shell
bwenv setup
```

Set up your Bitwarden host and your account information.

### Pull .env file from Bitwarden host

```shell
cd /path/to/your_project
bwenv pull
```

bwenv searchs your .env data in Bitwarden host with the current directory's name.
If it exists, pull data as .env file at current directory.
If already .env files current directory, bwenv asks overwrite them or not.
The data is stored as Bitwarden's Note item.

### Push .env file to Bitwarden host

bwenv pushs your .env data at the current directory to your Bitwarden host.
If it exists same name Bitwarden's Note item in dotenvs folder, bwenv asks overwrite it or not.

### List up .env datas in Bitwarden host

```shell
bwenv list
```

List up your .env datas from Bitwarden host.
They will showed up project names list on stdout.

## Uninstall

```shell
brew uninstall bwenv
```

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