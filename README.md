# Bonzai! Dominate the Command Line

[![Go Version](https://img.shields.io/github/go-mod/go-version/rwxrob/bonzai)](https://tip.golang.org/doc/go1.18)
[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai?status.svg)](https://godoc.org/github.com/rwxrob/bonzai)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/bonzai)](https://goreportcard.com/report/github.com/rwxrob/bonzai)

## Installation

🎉 ***Bonzai shamelessly requires Go 1.18+*** 💋

1. Install Go 1.18 and the tooling your require for it
1. `go install github.com/rwxrob/bonzai@latest` 
1. `import "github.com/rwxrob/bonzai"`
1. Consider using the [template][] to get started

[template]: <https://github.com/rwxrob/bonzai-template>

😎 *Yes, we use the wonderful new generics all [over](fn).* 👍

## Welcome to Bonzai

Yes, "banzai" is something people yell going into battle. But isn't
that what making command line utilities in Go (instead of your favorite
shell script) actually is? 

And yes, "bonsai" trees are well-manicured, meticulously crafted,
miniature trees that rival their larger cousins, just like Bonzai
command and data node trees. They are unlike anything you've probably
encountered so far, no getopt dashes (we kind of hate them), no ugly
commander interface to learn, no 12637 lines of shell tab completion
bloat to source before your command will complete, just well manicured
nested-tab-complete-with-magical-aliases-enabled commands organized into
rooted node trees for your command-line enjoyment. Your right-pinky will
be particularly grateful.

But wait, there's more! What about all those other tasks you need to do
to make a command line application honorable in anyone's eyes? Tools are
needed. 

## Create Powerful Parsers Easily with Pseudocode

How about the world's most approachable text scanner API with automatic
verbose, specific error messages that will make your users bow to your
attentiveness. Build complex parsers that are easily read and generated
from any number of meta-languages including PEGN, PEG, EBNF, and ABNF.
Spend your time working with the beautiful and standardized abstract
syntax trees instead of grinding out parser code. Domain language
specifications are a breeze, which is good, because you can customize
any completer for any nested command node or your application enabling
powerful, intelligent, and intuitive grammars you enable for your users
directly from the command line, hell, why not directly from your
favorite chat integration. After all, Bonzai commands are perfectly
doable directly from a chat client.

## A JSON Package Without Unnecessary Limitations

The standard Go `json` package is not what most people want. It uses
reflection which unnecessarily slows everything down, and produces
marshaled JSON output that escapes *every single Unicode character*
making the output impossible to read and virtual useless for any use
with languages that depend heavily on Unicode. It's really unfortunate
since the JSON standard fully supports Unicode characters as they are.
Bonzai JSON overcomes these limitations and is used to marshal
everything in the module and fulfill all `fmt.Stringer` interfaces as
2-space indented JSON instead of Go's virtually unusable default string
marshaling format.

## Single Module to Keep Legal Happy

Why have we put all of this stuff into one module? It's simple, if you
understand how enterprise software legal teams work. They are required
to track the legal pedigree of every single package that is included
into any software project. This task is laborious enough and we have
done everything to simplify it by providing most of the parts and tools
you would need to create almost any command-line utility, from simple
bash script port to monoliths with 100s of composed commands, even our
own BusyBox-like container Linux distro.
## Embedded Text or Web Docs FTW!

And, all the documentation for your command tree goes *inside* the
binary (if you want). That's right. Portable, text or web-enabled "man"
pages without the man page. You can use one of the composable
interactive-color-terminal-sensing help documentation commands like
`help.Cmd` that will easily marshal into JSON, or text, or well-styled
HTML locally while enabling you to write your embedded docs in
simplified CommonMark. Think "readthedocs" but without the Internet
dependency. And if you don't want `help.Cmd` you don't need it. You can
even write your own composable Bonzai command or pick from a rich
ecosystem of embeddable Bonzai command trees available from anywhere on
the Internet or maintained by the Bonzai project. No registries to worry
about. Just use good 'ol Go module imports are all that are need to
share your creations.

## Contributors/PRs Welcome

*... especially for "Completers" and Runtime Detection.*

Speaking of sharing, why not send a PR for your addition to the ever
growing collection of `comp` subpackage `Completers` for everything from
days of the week, to tab-driven inline math calculations, to a list of
all the listening ports running on your current system.

[CONTRIBUTING](CONTRIBUTING)

## "It's spelled bonsai/banzai."

We know. The domains were taken. Plus, this makes it more unique and
easier to find once you know the strange spelling we chose to use. Sorry
if that triggers your OCD.

If you must know, the primary motivator was the similarity to a
well-manicured tree (since it is for creating trees of commands). And
Buckaroo Banzai was always a favorite. We like to think he would use
Bonzai today to make amazing things.

On a lighter note, it just so happens that "banzai" means 'a traditional
Japanese idiom meaning "ten thousand years" of long life,' a cheer used
in celebrations. So combining the notion of a happy, little,
well-manicured, beautiful tree and "ten thousand years of long life"
works out just fine for us.

It turns out that the "call to war" associated with Bonzai is not
entirely without merit as well. Bonzai makes short work of creating
offensive and defensive tool kits all wrapping into one nice Go binary,
popular for building single-binary Linux container distros (like BusyBox
and Alpine, watch for Bonzai Linux soon), as well as root kits, and
other security tools

## What People Are Saying

> "It's like a modular, multicall BusyBox builder for Go with built in
> completion and embedded documentation support."

> "The utility here is that Bonzai lets you maintain your own personal
> 'toolbox' with built in auto-complete that you can assemble from
> various Go modules. Individual commands are isolated and unaware of
> each other and possibly maintained by other people." (tadasv123)

## Example GitHub Template

<https://github.com/rwxrob/bonzai-template>

## Design Considerations

* **Promote high-level package library API calls over Cmd bloat.** Code
  belongs in package libraries, not Cmds. While Bonzai allows for rapid
  applications development by putting everything initially in Cmd
  Methods, Cmds are most designed for documentation and completion, not
  heavy Method implementations. Eventually, most Method implementations
  should be moved into their own package libraries, perhaps even in the
  same Go module. Cmds should *never* be communicating with each other
  directly. While the idea of adding a Channel attribute was intriguing,
  it quickly became clear that doing so would promote too much code and
  tight coupling --- even with channels --- between specific commands.
  Cmds should be *very* light. In fact, most Cmds should assign their
  Method directly from one matching the Method function signature in a
  callable, high-level library.

* **Only bash completion planned.** Zsh, Powershell, and Fish have no
  equivalent to `complete -C` (which allows any executable to provide
  its own completion) This forces inferior dependencies on overly verbose
  external "completer" scripts written in only those languages for
  those specific shells. This dependency completely negates any
  possibility of providing modular completers and composable commands
  that carry their own completion logic. This one objective fact alone
  should give every pause before opting to use one of these inferior
  shells for their command line interactions. Bonzai commands leverage
  this best-of-breed completion functionality of bash to provide an
  unlimited number of completion methods and combinations. The
  equivalent implementations, perhaps as an export collected from all
  composed commands providing their shell equivalent of completion
  scripts, would be preposterously large (even though `kubectl` requires
  12637 lines of bash just for its basic completion). Bonzai uses Go
  itself to manage that completion --- drawing on a rich collection of
  completers included in the standard Bonzai module --- and provides
  documented shortcut aliases when completion is not available (h|help,
  for example). 

* **Bonzai commands may default to `shell.Cmd` or `help.Cmd`** These
  provide help information and optional interactive assistance
  including tab completion in runtime environments that do not have
  `complete -C foo foo` enabled. 

* **Scanner implementation as a struct with no interface.** In order to
  better tilt the scale toward Scanner performance the decision was made
  to keep it as a struct with tight coupling to the scanner.Cur and
  scanner.Pos structs as well.

## Style Guidelines

* Everything through `go fmt` or equiv, no exceptions
* In Vim `set textwidth=72` (not 80 to line numbers fit)
* Use `/* */` for package documentation comment, `//` elsewhere
* Smallest possible names for given scope while still clear
* Package globals that will be used a lot can be single capital
* Must be good reason to use more than 4 character pkg name
* Avoid unnecessary comments

## Acknowledgements

The <https://twitch.tv/rwxrob> community has been constantly involved
with the development of this project, making suggestions about
everything from my use of init, to the name "bonzai". While all their
contributions are too numerous to name specifically, they 
more than deserve a huge thank you here.

* Quint
* Greg 

## Legal 

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
SPDX-License-Identifier: Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the Bonzai project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
