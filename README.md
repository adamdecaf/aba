# aba

[![Build Status](https://github.com/adamdecaf/aba/workflows/Go/badge.svg)](https://github.com/adamdecaf/aba/actions)
[![Coverage Status](https://codecov.io/gh/adamdecaf/aba/branch/master/graph/badge.svg)](https://codecov.io/gh/adamdecaf/aba)
[![Go Report Card](https://goreportcard.com/badge/github.com/adamdecaf/aba)](https://goreportcard.com/report/github.com/adamdecaf/aba)
[![Apache 2 License](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/adamdecaf/aba/master/LICENSE)

aba is a CLI tool for looking up (and expanding) ABA routing numbers.

## Install

You can install from source:
```
go install github.com/adamdecaf/aba@latest
```

## Usage

Signup for a [Moov Account](https://dashboard.moov.io) and [setup the environment variables](https://github.com/moovfinancial/moov-go?tab=readme-ov-file#moov---go-client) for Moov's Go SDK.

```
$ aba -ach 27397636
Routing Number  Customer Name          Phone Number  Address
273976369       VERIDIAN CREDIT UNION  3192878332    1827 ANSBOROUGH AVE WATERLOO IA 50701
```

```
$ aba -rtp 071919133
Routing Number  Customer Name     Receive Payments  Receive Request for Payment
071919133       Fifth Third Bank  true              false
```

```
$ aba -wire 021201383
Routing Number  Telegraphic Name  Customer Name         Fund Transfers  Settlement Only  Book Entry Transfers  Address
021201383       VALLEY PASSAIC    VALLEY NATIONAL BANK  Y               N                Y
```

### Name Search

Searching by a Financial Institution's name is also supported.

```
$ aba -limit 3 Veridian
ACH:
Routing Number  Customer Name          Phone Number  Address
273975700       VERIDIAN CREDIT UNION  3192747598    1827 ANSBOROUGH AVE WATERLOO IA 50701
273976369       VERIDIAN CREDIT UNION  3192878332    1827 ANSBOROUGH AVE WATERLOO IA 50701
273976479       VERIDIAN CREDIT UNION  3192878332    1827 ANSBOROUGH AVE WATERLOO IA 50701

RTP:
Routing Number  Customer Name          Receive Payments  Receive Request for Payment
273976369       Veridian Credit Union  true              false

Wire:
Routing Number  Customer Name          Fund Transfers  Settlement Only  Book Entry Transfers  Address
273975700       VERIDIAN CREDIT UNION  true            false            false                 OELWEIN IA
273976369       VERIDIAN CREDIT UNION  true            false            false                 WATERLOO IA
031918828       MERIDIAN BANK          true            false            true                  MALVERN PA
```

## Supported and tested platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## License

Apache License 2.0 - See [LICENSE](LICENSE) for details.
