# bwsf

bwsf (Bitwarden Secured Files) is a CLI tool that uses [Bitwarden](https://bitwarden.com/) to manage .env files.

[日本語版はこちら](./README_ja.md)

## 🚨🚨 BREAKING CHANGE 🚨🚨

### Changed CLI Name

From v0.11.0, `bwenv` is re-named as `bwsf`. This is cause some bwenv commands already existed. We decieded to change our CLI name to avoid confusing.

#### MIGRATE

Rename youre setting directory.

```bash
mv ~/.config/bwenv ~/.config/bwsf
```

Uninstall your current version, and re-install latest version.

```bash
brew uninstall bwenv
brew install bwsf
```

### Multiple `.env.enviroment` files

From v0.9.0, bwsf stores multiple enviroment .env files, like `.env | .env.staging | .env.production`.

Cause with this, stored data at Bitwarden Note item structure is changed.

Stored data before v0.8.0 is no compatiblity after v0.9.0.

We will not provide migration system.

## Overview

bwsf commands supports your dotenv files are manged in your Bitwarden.

Simple usage below:

| command | |
|----|----|
| bwsf push | .env files push to your Bitwarden host |
| bwsf pull | .env files pull from your Bitwarden host |
| bwsf list | Show list stored .env files at your Bitwarden host |

## Motivation

We use Bitwarden as our password manager long time ago.
Also, we store .env files our Bitwarden host, manage them as shell scripts.
This project migrates our hand-maded shell scripts to modern CLI command with Go.

## Requirements

**`bw`** command is needed to be installed your machine.

[To install bw command, please read this document.](https://bitwarden.com/help/cli/#download-and-install)

**Homebrew**: Need to install bwenv.

### Machine OS

- macOS
- Linux
- [Is planning] Windows

## Installation

| OS | command |
|----|----|
| macOS | brew tap b4m-oss/tap && brew install bwsf |
| Linux | brew tap b4m-oss/tap && brew install bwsf |

> Note: Linux requires [Homebrew on Linux](https://docs.brew.sh/Homebrew-on-Linux) to be installed first.

## Verify installation

```shell
bwenv -v
# bwenv version 0.10.0
```

## Usage

### Initial setup

```shell
bwsf setup
```

Set up your Bitwarden host and your account information.

### Pull .env file from Bitwarden host

```shell
cd /path/to/your_project
bwsf pull
```

bwsf searchs your .env data in Bitwarden host with the current directory's name.
If it exists, pull data as .env file at current directory.
If already .env files current directory, bwsf asks overwrite them or not.
The data is stored as Bitwarden's Note item.

### Push .env file to Bitwarden host

bwsf pushs your .env data at the current directory to your Bitwarden host.
If it exists same name Bitwarden's Note item in dotenvs folder, bwsf asks overwrite it or not.

### List up .env datas in Bitwarden host

```shell
bwsf list
```

List up your .env datas from Bitwarden host.
They will showed up project names list on stdout.

## Uninstall

```shell
brew uninstall bwsf
```

## FAQ

<details>
<summary>Q. I don't have Bitwarden account.</summary>

To use bwenv, you need a Bitwarden account.

You can access to [Bitwarden Cloud](https://bitwarden.com/), sign up a account.

No fee, No credit card.

</details>

<details>
<summary>Q. I'm Bitwarden self hosted user.</summary>

Ofcourse, bwenv is available for Bitwarden self hosted users.

You can input your self hosted URL when initial setup.

</details>

<details>
<summary>Q. How does my .env file store at Bitwarden host?</summary>

Your .env files are converted to JSON syntax. bwenv creates Bitwarden Note item, put into Note section to JSON.

</details>

<details>
<summary>Q. Where are my Bitwarden account info</summary>

bwenv stores your config data at `~/.config/bwenv/`.

But, secure information (ex. password) is never stored.

</details>

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