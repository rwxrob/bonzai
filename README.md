# Go Bonzai™ Command Compositor

![logo](logo.png)

Meticulously manicured monoliths, on *any* device.

> Have a look at [rwxrob/z](https://github.com/rwxrob/z) for now to get a
sense of how it's coming along and how to use until the 1.0 release.

> Or, you can get started right away by cloning/forking [the sample
`foo` template repo](https://github.com/rwxrob/foo)

🚧 *under construction* 🚧

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

## Contributors/PRs Welcome

*... especially for "Completers", included popular commands, and Runtime
Detection.*

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

## "Why not just use Cobra?"

Just because something is popular (or first) doesn't mean it was well
designed. In fact, often inferior designs are rushed to market just to
gain adoption. Cobra seems to suffer from this. Discerning developers
and engineers have been not-so-quietly complaining about Cobra's
horrible design for years. It's time for something new.

Cobra requires wasteful and error-prone sourcing of thousands of lines
of shell code every time you create a new shell that needs to use a
Cobra command with shell tab completion (`kubectl` requires 12637). It
is not uncommon for operations people to be sourcing 100s of thousands
of lines of shell code just to enable basic completion that could have
been enabled easily with `complete -C` instead. Bonzai manages all
completion in Go instead of shell and therefore allows the modular
addition of any number of Completers including the standard file
completion as well as calculators, dates, and anything anyone can
conceive of completing. Completion is not dependent on any underlying
operating system. Any Bonzai command can provide its own completion
algorithm or use one of the many already provided. Cobra can never do
this.

Corba is also not designed to be a command compositor at all, which is
really unfortunate because they missed a golden opportunity there.
Bonzai branches can be imported and composed into other branches and
monoliths with just a few lines of Go. Registries of Bonzai commands can
be easily inferred from dependencies on the `bonzai` package and
creators are free to compose their monoliths from a rich eco-system of
Bonzai branches and commands. Bonzai allows creation of Go multicall
binary monoliths (like BusyBox) to be made easily, and from a diverse,
modular, importable, composable sources. Such is simply not possible
with Cobra and never will be.

Cobra buys into the broken boomer "getopt" view of the world requiring
people to remember all sorts of ungodly different combinations of dashes
and equals signs hoping things will just work. Bonzai takes a no-dashes
approach promoting much cleaner command lines with context and promotion
of domain specific languages (created with PEGN, scan.X, or others) that
easily translate directly to chat and other command-line interfaces for
most humans without much need to look up the documentation, which, by
the way, is embedded in the Bonzai command tree.

Cobra provides minimal, unappealing command documentation that is
virtually unreadable in source form. And Cobra provides no means of markup
or use of color and doesn't even promote the same look and feel of
manual page documentation. Bonzai has its own subset of Markdown,
BonzaiMark, respects the well established readability of manual pages,
and allows for the creation of elegant documentation that can be viewed
from the command line or easily from a local browser on the same
computer running the command. And the source containing the
documentation is as easy to read as the documentation itself.

In short, Cobra got us a long way, but has proved to be so laden with
crushing technical debt from failed base design decisions that it simply
is not sustainable given current modern expectations for good user
interfaces and documentation. Bonzai is a fresh, extensible,
sustainable, human-friendly command compositor to take us into the
future of command line interfaces, for everyone.

## What People Are Saying

> "It's like a modular, multicall BusyBox builder for Go with built in
> completion and embedded documentation support."

> "The utility here is that Bonzai lets you maintain your own personal
> 'toolbox' with built in auto-complete that you can assemble from
> various Go modules. Individual commands are isolated and unaware of
> each other and possibly maintained by other people." (tadasv123)

## Example GitHub Template

<https://github.com/rwxrob/foo>

## Design Considerations

* **Promote high-level package library API calls over Cmd bloat.**

  Code belongs in package libraries, not Cmds.

  While Bonzai allows for rapid applications development by putting
  everything initially in Cmd Call first-class function, Cmds are most
  designed for documentation and completion, not heavy Call
  implementations (even though many will organically start there from
  people's personal Bonzai trees).

  Eventually, most Call implementations should be moved into their own
  package libraries, perhaps even in the same Go module. Cmds should
  *never* communicate with each other directly (other than calling one
  another on occasion). While the idea of adding a Channel attribute was
  intriguing, it quickly became clear that doing so would promote 
  undesirable tight coupling --- even with channels --- between
  specific commands.

* **Cmds should be very light.** 

  Most Cmds should assign their first-class Call function to one that
  lightly wraps a similar function signature in a callable, high-level
  library that works entirely independently from the bonzai package.
  It's best to promote strong support for sustainable API packages.

* **Only bash completion and shell.Cmd planned.**

  If it doesn't work with `complete -C` or equivalent then just run the
  Bonzai command tree monolith as a temporary shell (shell.Cmd) and use
  its cross-platform support for tab completion.

  Zsh, Powershell, and Fish have no equivalent to `complete -C` (which
  allows any executable to provide its own completion). This forces
  inferior dependencies on overly verbose external "completer" scripts
  written in only those languages for those specific shells. This
  dependency completely negates any possibility of providing modular
  completers and composable commands that carry their own completion
  logic. This one objective fact alone should give everyone pause before
  opting to use one of these inferior shells for their command line
  interactions. 

  Bonzai commands leverage this best-of-breed completion functionality
  of bash to provide an unlimited number of completion methods and
  combinations. The equivalent implementations, perhaps as an export
  collected from all composed commands providing their shell equivalent
  of completion scripts, would be preposterously large just for its
  basic completion tree). Instead, Bonzai uses Go itself to manage
  that completion --- drawing on a rich collection of completers
  included in the standard Bonzai module --- and provides documented
  shortcut aliases when completion is not available (h|help, for
  example).

* **Bonzai commands may default to `shell.Cmd` or `help.Cmd`** 

  These provide help information and optional interactive assistance
  including tab completion in runtime environments that do not have
  `complete -C foo foo` enabled. 

  *shell.Cmd is still under development and likely will be for a while*

* **One major use case is to replace shell scripts in "dot files"
  collections.**

  By creating a `cmd` subdirectory of a dot files repo a multi-call
  Bonzai command named `cmd` can be easily maintained and added to just
  as quickly as any shell script. This has the added bonus of allowing
  others to quickly add one of your little commands with just a simple
  import (for example, `import github.com/rwxrob/dot/cmd` and then
  `cmd.Isosec`) from their own `cmd` monoliths. This also enables the
  fastest possible prototyping of code that would otherwise require
  significant, problematic mocks. Developers can work out the details of
  a thing just as fast as with shell scripting --- but with the power of
  all the Go standard library --- and then factor out their favorites as
  they grow into their own Bonzai command repos. This approach keeps "Go
  on the brain" (instead of having to port a bunch of bash later) and
  promotes the massive benefits of rapid applications development the
  fullest extent.

* **Use either `foo.Cmd` or `cmd.Foo` convention.**

  People may decide to put all their Bonzai commands into a single `cmd`
  package or to put each command into its own package. Both are
  perfectly acceptable and allow the developer making the import to
  alias the packages as needed using Go's excellent package import
  design.

## Style Guidelines

* Everything through `go fmt` or equiv, no exceptions
* In Vim `set textwidth=72` (not 80 to line numbers fit)
* Use `/* */` for package documentation comment, `//` elsewhere
* Smallest possible names for given scope while still clear
* Favor additional packages (possibly in `internal`) over long names
* Package globals that will be used a lot can be single capital
* Must be good reason to use more than 4 character pkg name
* Avoid unnecessary comments

## Acknowledgements

The <https://twitch.tv/rwxrob> community has been constantly involved
with the development of this project, making suggestions about
everything from my use of init, to the name "bonzai". While all their
contributions are too numerous to name specifically, they 
more than deserve a huge thank you here.

* <https://github.com/alessio/shellescape> ([shell.go](shell.go))

## Legal

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
SPDX-License-Identifier: Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the Bonzai™ project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
