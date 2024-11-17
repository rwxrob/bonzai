# 🌳 Go Bonzai™ CLI Framework and Library

***Bonzai is undergoing a major upgrade and a [v1.0 should be here by 2025](https://github.com/rwxrob/bonzai/issues/226). We recommend waiting until then before migrating any existing stuff or embarking on RRIB shell script collections in your Dot files. Come join us on https://linktr.ee/rwxrob to help us finish up the new version.***

[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai?status.svg)](https://godoc.org/github.com/rwxrob/bonzai)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/bonzai)](https://goreportcard.com/report/github.com/rwxrob/bonzai)

> "It's like a modular, multicall BusyBox builder for Go with built in completion and embedded documentation support."

> "The utility here is that Bonzai lets you maintain your own personal 'toolbox' with built in auto-complete that you can assemble from  various Go modules. Individual commands are isolated and unaware of  each other and possibly maintained by other people."

> "I used Bonzai for a pentest ... was able to inject ebpf.io kernel modules that were part of bin lol! Bonzai is Genius!"

![logo](assets/logo.png)

Bonzai™ was born from a very real need to replace messy collections of shell scripts, wasteful completion sourcing, and OS-specific documentation with a single, stateful, multicall binary composed of commands organized as rooted node trees with a clean, modular, portable, statically-compiled, and dynamically self-documenting design. There's simply nothing else like it in Go or any other language, and there's no better language than Go for such things. Crafting homekit/rootkit binaries with embedded resources that don't bloat RAM consumption is a breeze. No other language can do it.

Bonzai gets its name from the fact that Bonzai users are fond of meticulously manicuring their own stateful command trees, built from imported composite commands that they can easily copy and run on on any device, anywhere.

Bonzai users can easily share their own commands with others just like they would any other Go module and since most Bonzai commands also double as a high-level library package, even non-Bonzai users benefit. In fact, this monorepo is full of other Go modules containing cookbook recipe code for many of the things an avid techy with Go skills would want, the missing "batteries" that make Go really light up a well-crafted multicall command tree.

# Getting started

Take a look at the following commands to get an idea of what can be done:

- [`kimono` - Go monorepo utility](./cmds/kimono)
- [`help` - importable help command](./cmds/help)
- [`var` - persistent variable tool](./vars/cmd/var)
- [`sunrise` - fun terminal performance tester](./cmds/sunrise)

We have worked hard keep things as simple as possible so they are intuitive and to document this package as succinctly as possible so it is very usable from any decent tool that allows looking up documentation while writing the code.

## Legal

Copyright 2024 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
SPDX-License-Identifier: Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the Bonzai™ project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
