# A Go package for handling of runtime panics

[![Release](https://img.shields.io/github/release/agext/calmly.svg?style=flat)](https://github.com/agext/calmly/releases/latest)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/agext/calmly) 
[![Build Status](https://travis-ci.org/agext/calmly.svg?branch=master&style=flat)](https://travis-ci.org/agext/calmly)
[![Coverage Status](https://coveralls.io/repos/github/agext/calmly/badge.svg?style=flat)](https://coveralls.io/github/agext/calmly)
[![Go Report Card](https://goreportcard.com/badge/github.com/agext/calmly?style=flat)](https://goreportcard.com/report/github.com/agext/calmly)


Package calmly implements convenient runtime panic recovery and handling in [Go](http://golang.org).

## Project Status

v1.0.1 Stable: Guaranteed no breaking changes to the API in future v1.x releases. Probably safe to use in production, though provided on "AS IS" basis.

This package is being actively maintained. If you encounter any problems or have any suggestions for improvement, please [open an issue](https://github.com/agext/calmly/issues). Pull requests are welcome.

## Overview

When a panic condition needs to be handled by the program (rather than crashing it), wrap the code that can trigger such condition in a `Try`, which allows you to `Catch` the panic for further processing.

The `Outcome` of a `Try`ed code also offers convenience methods to:
- `KeepCalm` downgrading a panic to an error condition;
- `Escalate` upgrading a panic to a fatal error;
- `Log` the error, panic or fatal condition, using the appropriate logger method - presumably triggering a new panic or exiting the program.

## Installation

```
go get github.com/agext/calmly
```

## License

Package calmly is released under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
