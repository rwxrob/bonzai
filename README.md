# Bonzai! CLI Command Trees for Discerning Gophers

[![Go Version](https://img.shields.io/github/go-mod/go-version/rwxrob/bonzai)](https://tip.golang.org/doc/go1.18)
[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai?status.svg)](https://godoc.org/github.com/rwxrob/bonzai)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/bonzai)](https://goreportcard.com/report/github.com/rwxrob/bonzai)

Bonzai command trees are unlike anything you've probably encountered so
far, no getopt dashes (we kind of hate them), no ugly commander
interface to learn, no 12637 lines of shell tab completion bloat to
source before your command will complete, just well manicured
nested-tab-complete-with-magical-aliases-enabled commands organized into
rooted node trees for your command-line enjoyment. Your right-pinky will
be particularly grateful.

## Installation

🎉 ***Bonzai shamelessly requires Go 1.18+*** 💋

1. Install Go 1.18 and the tooling your require for it
1. `go install github.com/rwxrob/bonzai@latest` 
1. `import "github.com/rwxrob/bonzai"`
1. Consider using the [template][] to get started

[template]: <https://github.com/rwxrob/bonzai-template>

😎 *Yes, we use the wonderful new generics within the [`filter`](filter)
subpackage.*

## Embedded Text or Web Docs FTW!

And, all the documentation for your command tree goes *inside* the
binary (if you want). That's right. Portable, text or web-enabled "man"
pages without the man page. You can use one of the composable
interactive-color-terminal-sensing help documentation commands like
`cmd.Help` that will easily marshal into JSON, or text, or well-styled
HTML locally while enabling you to write your embedded docs in
simplified CommonMark. Think "readthedocs" but without the Internet
dependency. And if you don't want `cmd.Help` you don't need it. You can
even write your own composable Bonzai command
or pick from a rich ecosystem of embeddable Bonzai command trees
available from anywhere on the Internet or maintained by the Bonzai
project. No registries to worry about. Just use good 'ol Go module imports are all that are need to share your creations.

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

On a lighter note, it just so happens that "banzai" means 'a traditional
Japanese idiom meaning "ten thousand years" of long life,' a cheer used
in celebrations. So combining the notion of a happy, little,
well-manicured, beautiful tree and "ten thousand years of long life"
works out just fine for us.

Then, of course, there's the image of Buckaroo invoked every time we say
the name. In fact, I think we have our new theme song.

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

## Acknowledgements

The <https://twitch.tv/rwxrob> community has been constantly involved
with the development of this project, making suggestions about
everything from my use of init, to the name "bonzai". While all their
contributions are too numerous to name specifically, they deserve a
more than deserve a huge thank you here.

## Legal 

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
Licensed under Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the cmdbox project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
