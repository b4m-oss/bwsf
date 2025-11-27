# Getting Started

## What is bwenv?

**bwenv** is a CLI tool that uses [Bitwarden](https://bitwarden.com/) to manage `.env` files securely.

Instead of sharing `.env` files through insecure channels like email or Slack, bwenv lets you store them in your Bitwarden vault and sync them across your team.

## Prerequisites

Before using bwenv, make sure you have:

1. **Bitwarden Account** - Either [Bitwarden Cloud](https://bitwarden.com/) or a self-hosted Bitwarden server
2. **Bitwarden CLI (`bw`)** - The official Bitwarden command-line tool

### Installing Bitwarden CLI

Follow the [official Bitwarden CLI installation guide](https://bitwarden.com/help/cli/#download-and-install) to install the `bw` command on your machine.

Verify the installation:

```bash
bw --version
```

## How bwenv Works

bwenv stores your `.env` files in a special folder called `dotenvs` within your Bitwarden vault. Here's how the structure looks:

```
Bitwarden Vault
└── dotenvs/                    # Reserved folder for bwenv
    ├── my-web-app/             # Project name (directory name)
    │   ├── .env
    │   ├── .env.staging
    │   └── .env.production
    └── another-project/
        └── .env
```

::: info
The folder name `dotenvs` is reserved by bwenv. Don't use it for other purposes in your Bitwarden vault.
:::

## Initial Setup

After installing bwenv, run the setup command:

```bash
bwenv setup
```

This will configure:
- Your Bitwarden server URL (for self-hosted instances)
- Your Bitwarden account credentials

## Next Steps

- [Install bwenv](/en/guide/installation) - Installation instructions for your platform
- [Commands](/en/guide/commands) - Learn all available commands
