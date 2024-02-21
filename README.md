# aba

[![GoDoc](https://godoc.org/github.com/adamdecaf/aba?status.svg)](https://pkg.go.dev/github.com/adamdecaf/aba/pkg/aba)
[![Build Status](https://github.com/adamdecaf/aba/workflows/Go/badge.svg)](https://github.com/adamdecaf/aba/actions)
[![Coverage Status](https://codecov.io/gh/adamdecaf/aba/branch/master/graph/badge.svg)](https://codecov.io/gh/adamdecaf/aba)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamdecaf/aba)](https://goreportcard.com/report/github.com/adamdecaf/aba)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/adamdecaf/aba/master/LICENSE)

aba is a CLI tool for looking up (and expanding) ABA routing numbers.

## Usage

Signup for a [https://dashboard.moov.io](Moov Account) and [setup the environment variables](https://github.com/moovfinancial/moov-go?tab=readme-ov-file#moov---go-client) for Moov's Go SDK.

```
$ aba -ach 27397636
Routing Number  Customer Name          Phone Number  Address
273976369       VERIDIAN CREDIT UNION  3192878332    1827 ANSBOROUGH AVE WATERLOO IA 50701
```

```
$ aba -wire 021201383
Routing Number  Telegraphic Name  Customer Name         Fund Transfers  Settlement Only  Book Entry Transfers  Address
021201383       VALLEY PASSAIC    VALLEY NATIONAL BANK  Y               N                Y
```

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
