# process

[![Build Status](https://travis-ci.org/mantyr/process.svg?branch=master)](https://travis-ci.org/mantyr/process)
[![GoDoc](https://godoc.org/github.com/mantyr/process?status.png)](http://godoc.org/github.com/mantyr/process)
[![Go Report Card](https://goreportcard.com/badge/github.com/mantyr/process?v=1)][goreport]
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)

This is not a stable version.

## Description

- `Environment` - Environment for controlled execution of commands
- `Process` - Controlled operating system process
    - [x] Start
    - [x] Stop
        - [ ] graceful shutdown
    - [x] Done
    - [x] Status

### Supports platforms

- [x] unix (linux, mac)
- [ ] windows
- [ ] other

## Installation

    $ go get github.com/mantyr/process

## Author

[Oleg Shevelev][mantyr]

[mantyr]: https://github.com/mantyr

[build_status]: https://travis-ci.org/mantyr/process
[godoc]:        http://godoc.org/github.com/mantyr/process
[goreport]:     https://goreportcard.com/report/github.com/mantyr/process
