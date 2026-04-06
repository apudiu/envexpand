# Contributing to envexpand

Thanks for your interest in contributing! 🎉

## Contribution flow (Pull Requests)

This project accepts contributions through GitHub Pull Requests.

1. Fork the repository
2. Create a feature/fix branch from `main`
3. Make your changes
4. Run checks locally
5. Open a Pull Request to `main`

## Branch naming (recommended)

- `feat/<short-description>`
- `fix/<short-description>`
- `docs/<short-description>`
- `chore/<short-description>`

Examples:

- `feat/add-default-output-flag`
- `fix/handle-empty-env-line`

## Local checks before opening a PR

Run these commands before pushing:

```bash
gofmt -w .
go vet ./...
go test ./...
go build ./...
```

## Commit messages (recommended)

Use clear commit messages, for example:

- `feat: add support for ...`
- `fix: handle ...`
- `docs: update usage examples`

## Pull Request checklist

- [ ] Code is formatted (`gofmt`)
- [ ] `go vet` passes
- [ ] `go test` passes
- [ ] `go build` passes
- [ ] README/docs updated if behavior changed

## CI on PRs

Automated PR checks run on every Pull Request to `main`.
Your PR should pass all checks before merge.
