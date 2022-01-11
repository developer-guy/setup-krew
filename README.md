# setup-krew GitHub Action

[![Test GitHub Action](https://github.com/developer-guy/setup-krew/actions/workflows/testaction.yml/badge.svg?event=push)](https://github.com/developer-guy/setup-krew/actions/workflows/testaction.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/developer-guy/setup-krew)


This action enables you to download and install [kubernetes-sigs/krew](https://github.com/kubernetes-sigs/krew) binary.

`setup-krew` verifies the integrity of the `krew` release during installation by verifying its SHA256 against its SHA256 file release along the way with the binary itself.

## Usage

This action currently supports GitHub-provided Linux, macOS and Windows runners (self-hosted runners may not work).

Add the following entry to your Github workflow YAML file:

```yaml
uses: developer-guy/setup-krew@main
with:
 krew-version: "v0.4.2" # optional
```

Example using a pinned version:

```yaml
jobs:
  test_cosign_action:
    runs-on: ubuntu-latest
    name: Install Krew and test presence in path
    steps:
      - name: Install krew
        uses: developer-guy/setup-krew@main
        with:
          krew-version: 'v0.4.2'
      - name: Check install!
        run: krew version
```
