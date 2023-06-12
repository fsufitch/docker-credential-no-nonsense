# A no-nonsense Docker credential helper

Are you tired of `docker-credential-pass` or other Docker credential sources doing weird stuff? I was. Then I made this.

`docker-credential-no-nonsense` is a credential helper that implements the [official interface](https://github.com/docker/docker-credential-helpers#development) for credential helpers. It has no system-specific or external dependencies (unlike `pass` requiring a GPG key, etc). Instead, it **encrypts your password using AES-256** and stores it in a user-scoped file ready for later use.

Isn't that neat? Easy and quick, no nonsense.

## Installation

To install the credential helper, check the release downloads for your platform. Download that, and put it somewhere on your PATH.

## Build from source

**If you have Go and Make on your system**, you can build the binary by cloning the repository and running:

    make

The Makefile can also receive `DISTS` and `GOARGS` variables in order to build for a different system (or systems). For example, this will be very verbose and build the code for Linux AMD64 and Windows ARM64:

    make DISTS='linux/amd64 windows/arm64' GOARGS='-x -v'

**If you don't have Go,** but do have Docker, you can try this instead:

    docker build -q --target make . | xargs docker run -v "$(pwd):/app"

**If you don't have Docker**, I'm not sure why this project is of interest. 

## Usage

This package is executed by the official Docker credential helper entrypoint. Thus, its usage is:

    Usage: docker-credential-no-nonsense <store|get|erase|list|version>

To make Docker use it, edit your Docker config (`~/.docker/config.json`) and set `"credsStore": "no-nonsense"`.

The credential helper can be configured via environment variables:

**`NO_NONSENSE_ENC_KEY`** (required!). This is the encryption key proper, for encrypting/decrypting the secrets in the JSON storage. It is required and must be supplied via the environment variable. This is because it is the only way that Docker can pass the value through to the point where it internally uses the helper. **You must set this variable before using any `docker` commands.

You can set it for your current (Unix-ish) shell and subprocesses by using:

    export NO_NONSENSE_ENC_KEY=$(read -srp"Key: " && echo $REPLY)

Or for one time usage: 

    NO_NONSENSE_ENC_KEY=$(read -srp"Key: " && echo $REPLY) docker ...
    
**`NO_NONSENSE_CREDFILE`**. This is the JSON file that the helper uses to store data. If unset, it will use the default, as defined by XDG for your system. You can check what path it is using by providing the `--where` flag. 

    $ docker-credential-no-nonsense --where
    /home/developer/.local/share/dkr-no-nonsense-credfile.json

