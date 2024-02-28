
# `installer`

Quickly install pre-compiled binaries from Github releases.

Installer is an HTTP server which returns shell scripts. The returned script will detect platform OS and architecture, choose from a selection of URLs, download the appropriate file, un(zip|tar|gzip) the file, find the binary (largest file) and optionally move it into your `PATH`. Useful for installing your favourite pre-compiled programs on hosts using only `curl`.

[![GoDev](https://img.shields.io/static/v1?label=godoc&message=reference&color=00add8)](https://pkg.go.dev/github.com/divyam234/installer)

## Usage

```sh
# install <user>/<repo> from github
curl https://sh-install.vercel.app/<user>/<repo>@<release>! | bash
```

```sh
# search web for github repo <query>
curl https://sh-install.vercel.app/<query>! | bash
```

*Or you can use* `wget -qO- <url> | bash`

**Path API**

* `user` Github user (defaults to @jpillora, customisable if you [host your own](#host-your-own), searches the web to pick most relevant `user` when `repo` not found)
* `repo` Github repository belonging to `user` (**required**)
* `release` Github release name (defaults to the **latest** release)
* `!` When provided, downloads binary directly into `/usr/local/bin/` (defaults to working directory)

**Query Params**

* `?type=` Force the return type to be one of: `script` or `homebrew`
    * `type` is normally detected via `User-Agent` header
    * `type=homebrew` is **not** working at the moment – see [Homebrew](#homebrew)
* `?insecure=1` Force `curl`/`wget` to skip certificate checks
* `?as=` Force the binary to be named as this parameter value

## Examples

* https://sh-install.vercel.app/yudai/gotty@v0.0.12
* https://sh-install.vercel.app/mholt/caddy
* https://sh-install.vercel.app/rclone/rclone

    ```sh
    $ curl -s sh-install.vercel.app/mholt/caddy! | bash
    Downloading mholt/caddy v0.8.2 (https://github.com/mholt/caddy/releases/download/v0.8.2/caddy_darwin_amd64.zip)
    ######################################################################## 100.0%
    Downloaded to /usr/local/bin/caddy
    $ caddy --version
    Caddy 0.8.2
    ```

## Private repos

You'll have to set `GITHUB_TOKEN` on both your server (instance of `installer`) and client (before you run `curl https://sh-install.vercel.app/foobar?private=1 | bash`)

## Force a particular `user/repo`

In some cases, people want an installer server for a single tool

```sh
export FORCE_USER=zyedidia
export FORCE_REPO=micro
./installer
```

Then calls to `curl localhost:8080` will return the install script for `zyedidia/micro`
