# envexpand

Small Go CLI to expand variables inside `.env` files for deployment workflows (e.g. k3s + kustomize secrets), without using shell `envsubst` directly.

## Why this exists

This tool is intended for this monorepo where env files rely on variable expansion (for example `MAIL_FROM_NAME="${APP_NAME}"`) and you want deterministic expansion behavior when preparing production secrets.

## How it works

### Supported variable formats

- `${VAR}`
- `$VAR`

### Resolution order

For each variable reference, value is resolved in this order:

1. Variables already parsed earlier in the same input file
2. OS environment variables

If a variable is still not found, the placeholder is kept unchanged.

### Quoted values

Quoted values are expanded, and quotes are preserved.

Example:

```env
APP_NAME="Order Online"
MAIL_FROM_NAME="${APP_NAME}"
```

becomes:

```env
APP_NAME="Order Online"
MAIL_FROM_NAME="Order Online"
```

### Compact mode (`-c`)

When `-c` is enabled:

- strips full-line comments
- strips blank/whitespace-only lines
- outputs compact `KEY=VALUE` lines

## Usage (prebuilt binary, no Go required)

```bash
./bin/envexpand_<os>_<arch> -i <input-file> [-o <output-file>] [-c]
```

Where:

- `<os>` is one of: `linux`, `darwin`, `windows`
- `<arch>` is one of: `amd64`, `arm64`

On Windows, use `.exe` binaries:

```powershell
.\bin\envexpand_windows_amd64.exe -i C:\path\to\your\env\.env.example
```

Examples of Linux/macOS binaries:

- `./bin/envexpand_linux_amd64`
- `./bin/envexpand_darwin_arm64`

### Flags

- `-i` (required): input `.env` file path
- `-o` (optional): output file path
- `-c` (optional): compact output (strip comments/blank lines)
- `-h`: show help

### Output path behavior

If `-o` is not provided, output is written in the current working directory using:

`<base>_out<ext>`

Example:

- input: `api/.env.example`
- default output: `./.env_out.example`

## Examples

From this directory (`tools/envsubst`):

```bash
# Linux amd64 binary
./bin/envexpand_linux_amd64 -i /path/to/your/env/.env.example

# macOS arm64 binary
./bin/envexpand_darwin_arm64 -i /path/to/your/env/.env.example

# Default output path
./bin/envexpand_linux_amd64 -i /path/to/your/env/.env.example

# Explicit output path
./bin/envexpand_linux_amd64 -i /path/to/your/env/.env.example -o /path/to/your/output/.env.processed

# Compact output
./bin/envexpand_linux_amd64 -i /path/to/your/env/.env.example -c -o /path/to/your/output/.env.example.compact
```

### Dev-only (with Go)

If you are developing the tool itself:

```bash
go run . -i <input-file> [-o <output-file>] [-c]
```

## Build binaries

Use the provided script:

```bash
./build.sh
```

It builds minimal-size binaries for Linux, macOS, and Windows:

- `bin/envexpand_linux_amd64`
- `bin/envexpand_linux_arm64`
- `bin/envexpand_darwin_amd64`
- `bin/envexpand_darwin_arm64`
- `bin/envexpand_windows_amd64.exe`
- `bin/envexpand_windows_arm64.exe`

using:

- `CGO_ENABLED=0`
- `-trimpath`
- `-ldflags "-s -w -buildid="`

## Automated releases (GitHub)

This repository is configured to publish a GitHub Release automatically when you push a tag prefixed with `v`.

Examples:

- `v1.0.0`
- `v1.2.1`

To create a release:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

The release workflow will build all platform binaries and upload them as release assets.

## Notes

- Recursive expansion is supported.
- Circular references are guarded to avoid infinite recursion.
- Parsing is line-oriented and intended for `.env` style files.
