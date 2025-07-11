
# `installer`

Quickly install pre-compiled binaries from GitHub, Codeberg, or Forgejo releases with a single command.

Installer is an HTTP server which returns shell scripts. The returned script will:
1. Detect platform OS and architecture
2. Choose the appropriate binary from available URLs
3. Download and extract the file (supports zip, tar, gz, bz2)
4. Find and install the binary (optionally into your `PATH`)

Perfect for installing pre-compiled programs on any host with just `curl` or `wget`.

[![GoDev](https://img.shields.io/static/v1?label=godoc&message=reference&color=00add8)](https://pkg.go.dev/github.com/aljabri00056/installer)

## Quick Start

```sh
# Install latest release
curl aj-get.vercel.app/user/repo | bash

# Install specific release
curl aj-get.vercel.app/user/repo@v1.2.3 | bash

# Using wget
wget -qO- aj-get.vercel.app/user/repo | bash
```

## Supported Providers

```sh
# GitHub (default)
curl aj-get.vercel.app/github/user/repo | bash

# Codeberg
curl aj-get.vercel.app/codeberg/user/repo | bash

# Forgejo/Gitea
curl aj-get.vercel.app/forgejo/user/repo | bash
```

## Features

### Installation Location
- Installs to `/usr/local/bin` by default
- Use `move=0` to download to current directory:
  ```sh
  curl "aj-get.vercel.app/user/repo?move=0" | bash
  ```

### Multiple Binaries
Extract multiple binaries from an archive:
```sh
# Install both ffmpeg and ffprobe
curl "aj-get.vercel.app/BtbN/FFmpeg-Builds?as=ffmpeg,ffprobe" | bash
```

### Release Filtering
Filter releases by name:
```sh
# Only consider releases containing "gpl"
curl "aj-get.vercel.app/user/repo?include=gpl" | bash
```

### Platform Selection
Force specific platform:
```sh
# Force Windows binary
curl "aj-get.vercel.app/user/repo?platform=windows" | bash
```

### Architecture Selection
Force specific architecture:
```sh
# Force arm64 binary
curl "aj-get.vercel.app/user/repo?arch=arm64" | bash
```

## Windows Support
Run in PowerShell:
```powershell
# Standard installation
powershell -c "irm aj-get.vercel.app/user/repo?platform=windows | iex"

# With options
powershell -c "irm 'aj-get.vercel.app/user/repo?platform=windows&move=0' | iex"
```

## Private Repositories
Access private repos by setting a GitHub token:
```sh
# Via environment variable
export GITHUB_TOKEN="your-token"
curl aj-get.vercel.app/user/private-repo | bash

# Or via Authorization header
curl -H "Authorization: Bearer your-token" aj-get.vercel.app/user/private-repo | bash
```

## Popular Examples

### Command Line Tools
```sh
# Install Micro editor
curl aj-get.vercel.app/zyedidia/micro | bash

# Install Rclone
curl aj-get.vercel.app/rclone/rclone | bash

# Install Hugo
curl aj-get.vercel.app/gohugoio/hugo | bash

# Install gotty
curl aj-get.vercel.app/yudai/gotty@v0.0.12 | bash
```
