# Contributing to portkey-mcp-server

We welcome contributions to this project! Here's how to get started:

## Local Development

Execute the following to install git hooks in your local repo, which will ensure that mocks are regenerated and committed before pushing:
```shell
# From repo root
make install-hooks
```

If you are seeing stale linter errors coming from the result of `make lint` (part of those installed git hooks), you could try clearing your linter cache with `make lint-clear-cache`.

## Submitting a Pull Request

1. Fork this repo
2. Create a new branch (`git checkout -b my-feature`)
3. Make your changes and commit them with clear, descriptive commit messages
4. Run all tests locally using the commands below
5. Push to your fork (`git push origin my-feature`)
6. Open a Pull Request (PR) back to `main` in this repo

We will review your PR and work together to get it merged!

### Running Tests

Before submitting a PR, please consider running all tests locally to ensure your changes don't introduce any issues:

```shell
# From repo root

# Run all tests
make test

# Run benchmarks
make benchmark

# Generate mocks
make mocks

# Run linter
make lint

# Run security checks
make security
```

> Note: The above are run automatically in our GitHub Actions workflows. Running them locally before pushing reduces noise in your PR and helps catch issues early.

### CI/CD Requirements

All GitHub Actions jobs must pass before a PR can be merged. This includes:
- Build verification
- Unit tests
- Linting
- Security checks
- etc.

If any job fails, you can find detailed error output in the GitHub Actions artifacts:
1. Go to the "Actions" tab in the repository
2. Click on the failed workflow run
3. Click on the failed job
4. Look for the "Artifacts" section at the bottom of the job page
5. Download and review the relevant artifacts (e.g., `lint-report.json` for linter errors)

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
