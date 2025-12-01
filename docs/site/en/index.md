---
layout: home

hero:
  name: bwsf
  text: Secure .env Management
  tagline: Manage your .env files with Bitwarden CLI
  actions:
    - theme: brand
      text: Get Started
      link: /en/guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/b4m-oss/bwsf

features:
  - icon: 🔐
    title: Secure Storage
    details: Store your .env files securely in your Bitwarden vault. No more plain text secrets in shared drives.
  - icon: 🔄
    title: Easy Sync
    details: Push and pull .env files between your local machine and Bitwarden with simple commands.
  - icon: 📋
    title: Multi-Environment
    details: Manage multiple environment files (.env, .env.staging, .env.production) in a single project.
  - icon: 🖥️
    title: Cross-Platform
    details: Works on macOS and Linux. Windows support is planned.
---

## Quick Start

```bash
# Install via Homebrew
brew tap b4m-oss/tap && brew install bwsf

# Initial setup
bwsf setup

# Pull .env from Bitwarden
cd /path/to/your_project
bwsf pull

# Push .env to Bitwarden
bwsf push
```

## How It Works

bwsf uses the official Bitwarden CLI (`bw`) to securely store and retrieve your `.env` files. Your environment variables are stored as **Note items** in a dedicated `dotenvs` folder within your Bitwarden vault.

Each project's `.env` files are identified by the directory name, making it easy to organize and manage multiple projects.
