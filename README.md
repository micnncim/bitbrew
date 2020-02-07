# bitbrew

[![actions-workflow-Test][actions-workflow-Test-badge]][actions-workflow-Test]
[![codecov][codecov-badge]][codecov]
[![codefactor][codefactor-badge]][codefactor]
[![release][release-badge]][release]
[![license][license-badge]][license] 

![bitbrew](./_doc/bitbrew.svg)

## Description

[BitBar](https://github.com/matryer/bitbar) is empowered by your and third party's plugins. However it does not have a plugin manager. This is inconvenient for us (specially developers).

`Bitbrew` is a small but powerful package manager. This supports installation, uninstallation and sync for [bitbar-plugins](https://github.com/matryer/bitbar-plugins).

## Installation

### go get

```
$ go get github.com/micnncim/bitbrew/cmd/bitbrew
```

### Homebrew

```
$ brew install micnncim/tap/bitbrew
```

### GitHub Releases

Check out [Releases](https://github.com/micnncim/bitbrew/releases).

## Requirement

- [BitBar](https://github.com/matryer/bitbar)
- [GitHub Access Token](https://github.com/settings/tokens)

## Usage

`Bitbrew` synchronizes your plugins to `formula.yaml`. Every time you install or uninstall a plugin, your formula gets updated. And you can configure a plugin folder and formula file path via `config.yaml`. Check out [example](./_example).

You **NEED** specify the plugin folder as `BitBar Plugin Folder`.  

You can also find the usage with `bitbrew help`.

### Available commands

- `install, i`: bitbrew install `<FILENAME>`
- `uninstall, u`: bitbrew uninstall `<FILENAME>`
- `sync`: bitbrew sync
- `list, l`: bitbrew list
- `search, s`: bitbrew search `<TEXT>`
- `browse, b`: bitbrew browse `<FILENAME>`
- `config, c`: bitbrew config

### Search plugin

Search published plugins by any word.

![search](./_doc/search.svg)

### Install plugin

Installs a plugin by specifying its filename. You can find a plugin's filename by `bitbrew search`.

![install](./_doc/install.svg)

### Uninstall plugin

Uninstalls a plugin by specifying its filename. You can find a plugin's filename by `bitbrew list`.

![uninstall](./_doc/uninstall.svg)

### Sync plugins

Synchronizes your plugins to your local plugins. This concurrently installs and uninstalls plugins so it's fast.
It helps setup when you change your computer.

![sync](./_doc/sync.svg)

### Other

- `bitbrew list`: Lists your plugins
- `bitbrew config`: Edits your `config.yaml`
- `bitbrew browse`: Browses a plugin on https://getbitbar.com.

<!-- badge links -->

[actions-workflow-Test]: https://github.com/micnncim/bitbrew/actions?query=workflow%3ATest
[codecov]: https://codecov.io/gh/micnncim/bitbrew
[codefactor]: https://www.codefactor.io/repository/github/micnncim/bitbrew
[release]: https://github.com/micnncim/bitbrew/releases
[license]: LICENSE

[actions-workflow-Test-badge]: https://img.shields.io/github/workflow/status/micnncim/bitbrew/Test?label=Test&style=for-the-badge&logo=github
[codecov-badge]: https://img.shields.io/codecov/c/github/micnncim/bitbrew?style=for-the-badge&logo=codecov
[codefactor-badge]: https://img.shields.io/codefactor/grade/github/micnncim/bitbrew?logo=codefactor&style=for-the-badge
[release-badge]: https://img.shields.io/github/v/release/micnncim/bitbrew?style=for-the-badge&logo=github
[license-badge]: https://img.shields.io/github/license/micnncim/bitbrew?style=for-the-badge
