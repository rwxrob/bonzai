# 🌳 Go Bonzai™ Command Compositor

[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai?status.svg)](https://godoc.org/github.com/rwxrob/bonzai)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/bonzai)](https://goreportcard.com/report/github.com/rwxrob/bonzai)

> "It's like a modular, multicall BusyBox builder for Go with built in completion and embedded documentation support."

> "The utility here is that Bonzai lets you maintain your own personal 'toolbox' with built in auto-complete that you can assemble from  various Go modules. Individual commands are isolated and unaware of  each other and possibly maintained by other people."

![logo](logo.png)

Bonzai was born from a very real need to replace messy collections of shell scripts, wasteful completion sourcing, and OS-specific documentation with a single, stateful, binary composed of commands organized as rooted node trees with a clean, modular, portable, statically-compiled, and dynamically self-documenting design.

There's no better language than Go for such things.

Bonzai gets its name from the fact that Bonzai users are fond of meticulously manicuring their own stateful command trees, built from imported composite commands that they can easily copy and run on on any device, anywhere.

Bonzai users can easily share their own commands with others just like they would any other Go module and since any Bonzai command also doubles as a high-level library package, even non-Bonzai users benefit.

Most realize Bonzai really distinguishes itself from anything else out there the first time they turn any command branch into a fully-documented, tab-completing, stand-alone binary simply by wrapping it in five lines of code. Such is the beauty of stateful command tree design. There's simply nothing else like it, in Go or any other language.

## Getting Started

Read the book:

* https://github.com/rwxrob/book-bonzai

Copy or clone the example template:

* https://github.com/rwxrob/bonzai-example

Get ideas for your own by looking at others

* https://github.com/rwxrob/z
* https://github.com/rwxrob/pomo

## Legal

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
SPDX-License-Identifier: Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the Bonzai™ project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
