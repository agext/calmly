# A Go package for handling of runtime panics

Package calmly implements convenient runtime panic recovery and handling in [Go](http://golang.org).

## Maturity

[![Build Status](https://travis-ci.org/agext/calmly.svg?branch=master)](https://travis-ci.org/agext/calmly)

v1.0 Stable: Guaranteed no breaking changes to the API in future v1.x releases. No known bugs or performance issues. Probably safe to use in production, though provided on "AS IS" basis.

## Overview

[![GoDoc](https://godoc.org/github.com/agext/calmly?status.png)](https://godoc.org/github.com/agext/calmly)

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
