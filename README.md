
# `installer`

Quickly install pre-compiled binaries from Github releases.

Installer is an HTTP server which returns shell scripts. The returned script will detect platform OS and architecture, choose from a selection of URLs, download the appropriate file, un(zip|tar|gzip) the file, find the binary (largest file) and optionally move it into your `PATH`. Useful for installing your favourite pre-compiled programs on hosts using only `curl`.

[![GoDev](https://img.shields.io/static/v1?label=godoc&message=reference&color=00add8)](https://pkg.go.dev/github.com/divyam234/installer)

## Usage

```sh
# install <user>/<repo> from github
curl instl.vercel.app/<user>/<repo>@<release> | bash
```

*Or you can use* `wget -qO- <url> | bash`

**Path API**

* `repo` Github repository belonging to `user` (**required**)
* `release` Github release name (defaults to the **latest** release)
* `move=0` When provided as query param, downloads binary directly into working directory  (defaults to `/usr/local/bin/`)
* If no matching release is found you can  use `include="search term"` query param to filter release by search term.
* Extract multiple binaries from archive by providing `as=binary1,binary2` query param.
```sh
curl "https://instl.vercel.app/BtbN/FFmpeg-Builds?include=gpl-7.1&as=ffmpeg,ffprobe" | bash
```
## Windows (Run in PowerShell or Cmd)
```powershell
powershell -c "irm https://instl.vercel.app/rclone/rclone?platform=windows|iex"
```

## Examples

* instl.vercel.app/yudai/gotty@v0.0.12
* instl.vercel.app/mholt/caddy
* instl.vercel.app/rclone/rclone

## Private repos

You'll have to set pass github token in `GITHUB_TOKEN` env var.
```sh
GITHUB_TOKEN=token curl -H "Authorization: Bearer $GITHUB_TOKEN" instl.vercel.app/private/private-repo
```
