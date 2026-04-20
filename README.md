<p align="center">
  <h2 align="center">BACKLOG UTILS</h2>
  <p align="center">Backlog API utilities</p>
  <p align="center">
    <a href="https://github.com/nekrassov01/backlog-utils/actions/workflows/ci.yml"><img src="https://github.com/nekrassov01/backlog-utils/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI" /></a>
    <a href="https://pkg.go.dev/github.com/nekrassov01/backlog-utils"><img src="https://pkg.go.dev/badge/github.com/nekrassov01/backlog-utils.svg" alt="Go Reference" /></a>
    <a href="https://goreportcard.com/report/github.com/nekrassov01/backlog-utils"><img src="https://goreportcard.com/badge/github.com/nekrassov01/backlog-utils" alt="Go Report Card" /></a>
    <img src="https://img.shields.io/github/license/nekrassov01/backlog-utils" alt="LICENSE" />
    <a href="https://deepwiki.com/nekrassov01/backlog-utils"><img src="https://deepwiki.com/badge.svg" alt="Ask DeepWiki" /></a>
  </p>
</p>

## Overview

Backlog Utils is a client tool for the Backlog API.

## Features

At this time we only support Wiki operations.

- List wiki pages with optional pattern.
- Rename wiki page
- Replace strings in the content of wiki page
- List wiki pages and rename them with optional pattern
- List wiki pages and replace strings in the content with optional pattern.

## Commands

```text
NAME:
   bkl - Backlog utilities

USAGE:
   bkl [global options] [command [command options]]

VERSION:
   0.0.1 (revision: 950fbaa)

DESCRIPTION:
   A cli application for Backlog utilities.

COMMANDS:
   wiki  Backlog wiki utilities

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Wiki subcommands

```text
NAME:
   bkl wiki - Backlog wiki utilities

USAGE:
   bkl wiki [command [command options]] 

COMMANDS:
   list         List wiki pages with optional pattern
   rename       Rename wiki page
   replace      Replace strings in the content of wiki page
   rename-all   List wiki pages and rename them with optional pattern
   replace-all  List wiki pages and replace strings in the content with optional pattern

OPTIONS:
   --help, -h  show help
```

#### List

```text
NAME:
   bkl wiki list - List wiki pages with optional pattern

USAGE:
   bkl wiki list

OPTIONS:
   --log-level string    set log level (default: "INFO") [$BACKLOG_LOG_LEVEL]
   --base-url string     set backlog base url [$BACKLOG_URL]
   --api-key string      set backlog api key [$BACKLOG_API_KEY]
   --project-key string  set backlog project key
   --pattern string      set pattern to search for wiki pages
   --help, -h            show help
```

#### Rename

```text
NAME:
   bkl wiki rename - Rename wiki page

USAGE:
   bkl wiki rename [command [command options]]

OPTIONS:
   --log-level string  set log level (default: "INFO") [$BACKLOG_LOG_LEVEL]
   --base-url string   set backlog base url [$BACKLOG_URL]
   --api-key string    set backlog api key [$BACKLOG_API_KEY]
   --wiki-id int       set backlog wiki id
   --old string        set string to be replaced in wiki page
   --new string        set new string after replacement in wiki page
   --help, -h          show help
```

#### Replace

```text
NAME:
   bkl wiki replace - Replace strings in the content of wiki page

USAGE:
   bkl wiki replace [command [command options]]

OPTIONS:
   --log-level string                 set log level (default: "INFO") [$BACKLOG_LOG_LEVEL]
   --base-url string                  set backlog base url [$BACKLOG_URL]
   --api-key string                   set backlog api key [$BACKLOG_API_KEY]
   --wiki-id int                      set backlog wiki id
   --pairs string [ --pairs string ]  set pairs of old and new repalacements for wiki page
   --help, -h                         show help
```

#### Rename All

```text
NAME:
   bkl wiki rename-all - List wiki pages and rename them with optional pattern

USAGE:
   bkl wiki rename-all [command [command options]]

OPTIONS:
   --log-level string    set log level (default: "INFO") [$BACKLOG_LOG_LEVEL]
   --base-url string     set backlog base url [$BACKLOG_URL]
   --api-key string      set backlog api key [$BACKLOG_API_KEY]
   --project-key string  set backlog project key
   --pattern string      set pattern to search for wiki pages
   --old string          set string to be replaced in wiki page
   --new string          set new string after replacement in wiki page
   --help, -h            show help
```

#### Replace All

```text
NAME:
   bkl wiki replace-all - List wiki pages and replace strings in the content with optional pattern

USAGE:
   bkl wiki replace-all [command [command options]]

OPTIONS:
   --log-level string                 set log level (default: "INFO") [$BACKLOG_LOG_LEVEL]
   --base-url string                  set backlog base url [$BACKLOG_URL]
   --api-key string                   set backlog api key [$BACKLOG_API_KEY]
   --project-key string               set backlog project key
   --pattern string                   set pattern to search for wiki pages
   --pairs string [ --pairs string ]  set pairs of old and new repalacements for wiki page
   --help, -h                         show help
```

## Installation

Install with homebrew

```sh
brew install nekrassov01/tap/backlog-utils
```

Install with go

```sh
go install github.com/nekrassov01/backlog-utils/cmd/bkl@latest
```

Or download binary from [releases](https://github.com/nekrassov01/backlog-utils/releases)

## Prerequisites

Set the following environment variables.

```sh
export BACKLOG_URL=https://your-space.backlog.jp
export BACKLOG_API_KEY=****
```

## Completion

Shell completion support if bash, fish, pwsh, and zsh.

```sh
source <(bkl completion bash)
```

## Todo

- [ ] Non-Wiki implementations
- [ ] Colorized logger implementations

## Author

[nekrassov01](https://github.com/nekrassov01)

## License

[MIT](https://github.com/nekrassov01/backlog-utils/blob/main/LICENSE)
