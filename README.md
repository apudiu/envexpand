# envexpand

[![PR Checks](https://github.com/apudiu/envexpand/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/apudiu/envexpand/actions/workflows/pr-checks.yml)
[![Release](https://github.com/apudiu/envexpand/actions/workflows/release.yml/badge.svg)](https://github.com/apudiu/envexpand/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/v/release/apudiu/envexpand?sort=semver)](https://github.com/apudiu/envexpand/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/apudiu/envexpand)](https://github.com/apudiu/envexpand/blob/main/go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](./CONTRIBUTING.md)

A small Go CLI that expands variables in `.env`-style files with deterministic behavior, without relying on shell `envsubst`.

## Origin

This tool started as an internal utility to process env files in one of our projects. It solved a recurring need around predictable variable expansion during deployment preparation. We later open-sourced it in this repository because teams with similar workflows may find it useful too.

## Highlights

- ✅ PR-based contribution workflow with automated CI checks
- ✅ Automated versioned releases on tags like `v1.2.1`
- ✅ Cross-platform binary builds for Linux/macOS/Windows (amd64 + arm64)
- ✅ Deterministic `.env` expansion with circular-reference protection

## Quick start

Use a prebuilt binary:

```bash
./bin/envexpand_<os>_<arch> -i <input-file> [-o <output-file>] [-c]
```

Where:

- `<os>`: `linux`, `darwin`, `windows`
- `<arch>`: `amd64`, `arm64`

Or run from source during development:

```bash
go run . -i <input-file> [-o <output-file>] [-c]
```

## How it works

### Supported variable formats

- `${VAR}`
- `$VAR`

### Resolution order

For each variable reference, values are resolved in this order:

1. Variables parsed earlier in the same input file
2. OS environment variables

If still unresolved, the placeholder is preserved as-is.

### Quoted values

Quoted values are expanded and quotes are preserved.

Example:

```env
APP_NAME="Order Online"
MAIL_FROM_NAME="${APP_NAME}"
```

Becomes:

```env
APP_NAME="Order Online"
MAIL_FROM_NAME="Order Online"
```

### Compact mode (`-c`)

When compact mode is enabled:

- full-line comments are removed
- blank/whitespace-only lines are removed
- output is normalized to compact `KEY=VALUE` lines

## Flags

- `-i` (required): input `.env` file path
- `-o` (optional): output file path
- `-c` (optional): compact output
- `-h`: show help

## Output path behavior

If `-o` is not provided, output is written in the current working directory as:

`<base>_out<ext>`

Example:

- input: `api/.env.example`
- default output: `./.env_out.example`

## Examples

```bash
# Linux amd64 binary
./bin/envexpand_linux_amd64 -i /path/to/.env.example

# macOS arm64 binary
./bin/envexpand_darwin_arm64 -i /path/to/.env.example

# Explicit output path
./bin/envexpand_linux_amd64 -i /path/to/.env.example -o /path/to/.env.processed

# Compact output
./bin/envexpand_linux_amd64 -i /path/to/.env.example -c -o /path/to/.env.compact
```

## Build binaries

```bash
./build.sh
```

Build output:

- `bin/envexpand_linux_amd64`
- `bin/envexpand_linux_arm64`
- `bin/envexpand_darwin_amd64`
- `bin/envexpand_darwin_arm64`
- `bin/envexpand_windows_amd64.exe`
- `bin/envexpand_windows_arm64.exe`

Build flags:

- `CGO_ENABLED=0`
- `-trimpath`
- `-ldflags "-s -w -buildid="`

## Quality signals

- PR CI validates formatting, vet, tests, and build
- Release CI also runs quality checks before publishing assets
- Unit tests cover core expansion behavior and edge cases

## Automated releases (GitHub)

A release is published automatically when you push a tag prefixed with `v`.

Examples:

- `v1.0.0`
- `v1.2.1`

Create a release:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## Contributing

Contributions are welcome via Pull Requests.

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for contribution flow, branch naming, validation commands, and PR checklist.

## Notes

- Recursive expansion is supported.
- Circular references are guarded to avoid infinite recursion.
- Parsing is line-oriented and intended for `.env`-style files.

## License

MIT — see [LICENSE](./LICENSE).
