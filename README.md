# bitbrew

[![CircleCI](https://circleci.com/gh/micnncim/bitbrew.svg?style=svg)](https://circleci.com/gh/micnncim/bitbrew)
[![Go Report Card](https://goreportcard.com/badge/github.com/micnncim/bitbrew)](https://goreportcard.com/report/github.com/micnncim/bitbrew)
[![codecov](https://codecov.io/gh/micnncim/bitbrew/branch/master/graph/badge.svg)](https://codecov.io/gh/micnncim/bitbrew)
[![Maintainability](https://api.codeclimate.com/v1/badges/6481fea60b20eefb9af9/maintainability)](https://codeclimate.com/github/micnncim/bitbrew/maintainability)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/1b68067c1d53421e96eee157d8fc349f)](https://www.codacy.com/app/micnncim/bitbrew?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=micnncim/bitbrew&amp;utm_campaign=Badge_Grade)
[![CodeFactor](https://www.codefactor.io/repository/github/micnncim/bitbrew/badge)](https://www.codefactor.io/repository/github/micnncim/bitbrew)
[![codebeat badge](https://codebeat.co/badges/9b906c1d-c209-4a9d-a560-5f866d296378)](https://codebeat.co/projects/github-com-micnncim-bitbrew-master)


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

## LICENSE

[MIT](./LICENSE)
