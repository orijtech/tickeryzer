# tickeryzer

![Build status](https://github.com/orijtech/tickeryzer/workflows/Go/badge.svg?branch=main)

Package tickeryzer defines an [Analyzer](analyzer_link) that checks missing (*time.Ticker).Stop() call,
which can cause resources leak.

## Installation

With Go modules:

```sh
go get github.com/orijtech/tickeryzer/cmd/tickeryzer
```

Without Go modules:

```sh
$ cd $GOPATH/src/github.com/orijtech/tickeryzer
$ git checkout v0.0.1
$ go get
$ install ./cmd/tickeryzer
```

## Usage

You can run `tickeryzer` either on a Go package or Go files, the same way as
other Go tools work.

Example:

```sh
$ tickeryzer github.com/orijtech/tickeryzer/testdata/src/a
```

or:

```sh
$ tickeryzer ./testdata/src/a/a.go
```

Sample output:

```text
/go/src/github.com/orijtech/tickeryzer/testdata/a/a.go:9:2: missing ticker.Stop() call
```
 
## Development

Go 1.15+

### Running test

Add test case to `testdata/src/a` directory, then run:

```shell script
go test
```

## Contributing

All contributions are welcome, please report bug or open a pull request.

[analyzer_link]: https://pkg.go.dev/golang.org/x/tools/go/analysis#Analyzer
