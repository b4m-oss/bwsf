# Contributing to bwenv

Thank you for your interest in contributing to bwenv! We welcome contributions from the community.

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming environment for everyone.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Docker** - Required for the development environment
- **Git** - For version control

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/bwenv.git
cd bwenv
```

3. Build the Docker image:

```bash
make build
```

4. Run the application:

```bash
make run
```

For more details, see the [Development section in README.md](./README.md#development).

## Branch Strategy (Git Flow)

We use Git Flow for branch management:

| Branch | Purpose |
|--------|---------|
| `main` | Production-ready releases |
| `develop` | Integration branch for development |
| `feature/*` | New features (e.g., `feature/add-sync-command`) |
| `fix/*` | Bug fixes (e.g., `fix/login-error`) |
| `hotfix/*` | Urgent production fixes |

### Creating a Branch

Always branch from `develop` for new features or fixes:

```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

| Prefix | Description |
|--------|-------------|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation changes |
| `test:` | Adding or updating tests |
| `refactor:` | Code refactoring (no functional changes) |
| `chore:` | Maintenance tasks (dependencies, configs, etc.) |

### Examples

```
feat: add sync command for bidirectional updates
fix: resolve authentication timeout issue
docs: update installation instructions
test: add unit tests for pull command
refactor: simplify config parsing logic
chore: update Go dependencies
```

## Pull Request Process

1. **Create a branch** from `develop`

2. **Make your changes** and commit following the commit convention

3. **Run tests** to ensure nothing is broken:

```bash
make test
```

4. **Run linting** to ensure code quality:

```bash
make lint
```

5. **Push your branch** to your fork:

```bash
git push origin feature/your-feature-name
```

6. **Open a Pull Request** targeting the `develop` branch

7. **Wait for review** - maintainers will review your PR and may request changes

### PR Guidelines

- Keep PRs focused on a single feature or fix
- Include a clear description of what the PR does
- Reference any related issues (e.g., "Closes #123")
- Ensure all tests pass before requesting review

## Reporting Issues

We use **Issues** for bug reports and **Discussions** for feature requests and questions.

### Bug Reports (Issues)

Found a bug? Please [open an Issue](https://github.com/b4m-oss/bwenv/issues/new/choose) using the bug report template.

The template will guide you through providing:

- Steps to reproduce the issue
- Expected vs actual behavior
- Your environment (OS, bwenv version, bw CLI version)
- Error logs if applicable

### Feature Requests (Discussions)

Have an idea for a new feature? Please use [GitHub Discussions](https://github.com/b4m-oss/bwenv/discussions/categories/ideas) instead of Issues.

This allows the community to discuss and refine ideas before implementation.

## Running Tests

We use the following Make commands for testing:

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-unit` | Run unit tests only |
| `make test-e2e` | Run E2E tests (mock-based) |
| `make test-coverage` | Generate coverage report |
| `make lint` | Run formatter and static analysis |

Always run `make test` and `make lint` before submitting a PR.

## Questions?

If you have any questions, please use [GitHub Discussions](https://github.com/b4m-oss/bwenv/discussions/categories/q-a) instead of opening an Issue.

Thank you for contributing!

