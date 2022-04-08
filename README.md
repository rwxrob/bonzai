# 🌳 Go Bonzai™ Command Compositor

[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai?status.svg)](https://godoc.org/github.com/rwxrob/bonzai)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)
[![Go Report
Card](https://goreportcard.com/badge/github.com/rwxrob/bonzai)](https://goreportcard.com/report/github.com/rwxrob/bonzai)

Meticulously manicured monolith and multicall binaries, built from
imported composite commands, on *any* device, with recursive,
light-weight tab completion, and colorful embedded documentation from
terminal or local web browser. Replace messy collections of shell
scripts ported to clean Go code and compiled into a single, portable `z`
command.

## 🤚 TL;DR; 🛑

Just want to see one? Check or clone the sample `foo` template:

👆 <https://github.com/rwxrob/foo>

Or, have a look at the `z` monolith/multicall that started it all:

👆 <https://github.com/rwxrob/z>

![logo](logo.png)

## Installation

🎉 ***Bonzai shamelessly requires Go 1.18+*** 💋

1. Install Go 1.18 and the tooling your require for it
1. `go install github.com/rwxrob/bonzai@latest` 
1. `import Z "github.com/rwxrob/bonzai"`
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
well-manicured tree (since it is for creating trees of composite
commands).

The misunderstood word "banzai" is 'a traditional Japanese idiom
meaning "ten thousand years" of long life,' a cheer used in
celebrations like "Hurrah" or "Viva".' So combining the notion of a
happy, well-manicured, beautiful tree and "ten thousand years of
long life" works out just fine for us.

And, yes, Buckaroo Banzai was always a favorite. We like to think he
would use Bonzai today to make amazing things and last for a long time
to defeat evil aliens and save the world.

It turns out that the "call to war" associated with Bonzai is not
entirely without merit as well. Bonzai is excellent for unorthodox,
rapid applications development (instead of writing scripts) and makes
short work of creating offensive and defensive tool kits all wrapping
into one nice Go multicall binary, popular for building single-binary
Linux container distros like BusyBox and Alpine, as well as root kits,
and other security tools

## "Why not just use Cobra?"

Just because something is popular (or first) doesn't mean it was well
designed. In fact, often inferior designs are rushed to market just to
gain adoption. Cobra seems to suffer from this. Discerning developers
and engineers have been not-so-quietly complaining about Cobra's
horrible design for years. It's time for something new. Read on if you
want the specific reasons.

* **Cobra tab completion is wasteful and error-prone.**

  Cobra often requires sourcing thousands of lines of shell code every
  time you run a new shell that needs to use a Cobra command with shell
  tab completion (`kubectl` requires 12637). It is not uncommon for
  operations people to be sourcing 100s of thousands of lines of shell
  code just to enable basic completion that could have been enabled
  easily with `complete -C` instead. Bonzai manages all completion in Go
  instead of shell and therefore allows the modular addition of any
  number of Completers including the standard file completion as well as
  calculators, dates, and anything anyone can conceive of completing.
  Completion is not dependent on any underlying operating system. Any
  Bonzai command can provide its own completion algorithm or use one of
  the many already provided. Cobra can never do this.

* **Corba is not designed to be a command compositor at all.**

  This is really unfortunate because the designers missed a golden
  opportunity. Bonzai branches can be imported and composed into other
  branches and monoliths with just a few lines of Go. Registries of
  Bonzai commands can be easily inferred from dependencies on the
  `bonzai` package and creators are free to compose their monoliths or
  multicall binaries from a rich eco-system of Bonzai branches and
  commands. Bonzai allows creation of Go multicall binary monoliths
  (like BusyBox) to be made easily, and from a diverse, modular,
  importable, composable sources. Such is simply not possible with Cobra
  and never will be.

* **Cobra suffers from broken boomer "getopt" design.**

  The world if finally realizing how bad dashed arguments and options
  have always been for *good* human-computer interactions from the
  command line, perhaps because more people are using chat interfaces as
  their command line. People simply cannot remember all sorts of ungodly
  combinations of dashes and equals signs hoping things will
  just work. Bonzai takes a no-dashes approach with aliases promoting
  cleaner, understandable command lines with context and promotion of
  domain specific languages (created with PEGN, scan.X, or others) that
  easily translate directly to chat and other command-line interfaces
  that most humans can use without even looking up the documentation,
  which, by the way, is embedded in any Bonzai command tree.

* **Cobra provides bad, brittle, command documentation.**

  Cobra documentation is virtually unreadable in source form. And Cobra
  provides no means of markup or use of color and doesn't even promote
  the same look and feel of manual page documentation. In contrast,
  Bonzai has its own subset of Markdown, BonzaiMark, respects the well
  established readability of manual pages, and allows for the creation
  of elegant documentation that can be viewed from the command line or
  easily from a local browser on the same computer running the command.
  Bonzai command documentation is as easy to read in source form as the
  documentation itself.

* **Cobra suffers from crushing technical debt.**

  The problems listed (and more) are never going to come out of Cobra.
  Because it is filled with bad design decisions and was rushed to
  market without serious consideration for its API, it is now doomed to
  never lose its warts (kinda like JavaScript). There is no possible way
  it can ever upgrade to address the very reasonable modern expectations
  for good command line user experiences. No wonder you never see people
  using Cobra for their replace-my-shell-scripts utilities. Cobra is
  simply horrible for this. Thankfully, Bonzai is a fresh, extensible,
  sustainable, human-friendly command compositor to take us into the
  future of command line interfaces.

## What People Are Saying

> "It's like a modular, multicall BusyBox builder for Go with built in
> completion and embedded documentation support."

> "The utility here is that Bonzai lets you maintain your own personal
> 'toolbox' with built in auto-complete that you can assemble from
> various Go modules. Individual commands are isolated and unaware of
> each other and possibly maintained by other people." (tadasv123)

## Examples

* <https://github.com/rwxrob/foo> - clone-able GitHub template
* <https://github.com/rwxrob/z> - first Bonzai command tree ever made

## Design Considerations

* **Promote high-level package library API calls over Cmd bloat**

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

* **Cmds should be very light**

  Most Cmds should assign their first-class Call function to one that
  lightly wraps a similar function signature in a callable, high-level
  library that works entirely independently from the bonzai package.
  It's best to promote strong support for sustainable API packages.

* **Only bash completion and shell.Cmd planned**

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

* **Target users replacing shell scripts in "dot files" collections**

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

* **Use either `foo.Cmd` or `cmd.Foo` convention**

  People may decide to put all their Bonzai commands into a single `cmd`
  package or to put each command into its own package. Both are
  perfectly acceptable and allow the developer making the import to
  alias the packages as needed using Go's excellent package import
  design.

* **Default capital "Z" import name**

  People can easily change the default "Z" import name to whatever else
  if they don't like it (or worse, if has conflicts). But, we actually
  want naming conflicts even though this seems counter-intuitive.
  Developers should be putting most of their code into their own `pkg`
  libraries and calling into them from their wrapping Bonzai trees and
  branches that import Z. If someone is importing Z into their package
  library they are likely doing something they shouldn't. Bonzai should
  *only* be imported into the composable branch or standalone command
  (`main`). This is a reminder as well to Bonzai developers not to stuff
  things into the Z package that would never be used outside a Bonzai
  command, tree or `main.go` standalone.

* **Promote lower "z" personal Bonzai trees**

  Just as "dotfiles" has become a convention, use of the simple "z"
  GitHub repo should be promoted as a common, "ez" way to find people's
  personal monoliths. Obviously, people can use what they want. This is
  also consistent with the capital "Z" import name of the `bonzai`
  package.

* **Use simple regular expressions for usage**

  Bonzai takes a fundamentally different approach to the command line
  and usage documentation. Any command line is a minimal domain specific
  language. As such, usage notation simple is not enough. Regular
  expressions allow the quick understanding of what is allowed and
  should become mandatory learning in a world of minimal domain specific
  grammars. Only the most basic regular expressions are required to
  produce rich usage strings for any Bonzai command. However, those
  wanting traditional usage notation can easily override the
  `DefaultUsageFunc` with their own.

* **Prioritize output to VT100-compatible terminals**

  Every modern terminal supports VT100 compatibility. If it doesn't,
  people should just not use it. This means emphasis and formatting are
  dependent on the [`rwxrob/term`](https://github.com/rwxrob/term)
  package for the main output. Bonzai tree developers will likely want
  terminal output helpers more than anything (yes even web rendering).

* **Secondary output to local web browser**

  Bonzai will eventually provide the option to enable use of the local
  web browser instead of terminal output for help and other
  documentation output. This will allow graphic application executables
  to simple be double-clicked from graphic desktops. Bonzai will ship
  with an embedded web template for such applications, but will also
  allow users to granularly customize their own modifications to the
  default theme. Bonzai branch creators are encouraged to provide
  downloadable themes in this regard. In this way, Bonzai will provide
  the web shell to encapsulate other web applications.

* **Prefer local web apps over other GUI platforms**

  While it is obviously possible to create any graphic application with
  Bonzai, the creation of localized web applications should be the
  preferred convention of the Bonzai GUI applications community. This
  simplifies application development and takes aim at bloated
  alternatives that embed full GUI web clients (such as Electron) while
  still providing a rich terminal interface for those who prefer it.
  Every Bonzai binary is its own local web server which Go's rich
  standard library to draw on.

* **[]Section instead of map[string]string for Other**

  This allows composition notation and ensures the author has control of
  how the Other sections appear.

* **Move `help.Cmd` into its own package**

  Although it breaks backward compatibility for many applications
  updating between `v0.1` and `v0.2` the decision to put `help.Cmd` into
  it's own Bonzai branch git repo was the right one. It is now on equal
  footing with every other potential Bonzai branch and can keep its
  `Version` in sync with Git version (as all Bonzai branches should). It
  is highly likely that GUI/web/hybrid help commands will be preferred
  by some and including one over another --- when not used --- ends up
  just being unnecessary bloat. This also serves to clarify that the
  legal information is related to that specific `help.Cmd` Bonzai branch
  and not Bonzai itself. It's conceivable that another `help.Cmd`
  creator may wish another legal agreement.

* **Leave hidden default command params**

  When the default command is invoked any of it's params will be
  automatically passed as if the command specified them. But they are
  not included in the tab completion. This is because there will
  inevitably be conflicts between default command params and other
  potential completions at that level for other commands. Rather than
  disable this, or add tab completion, it was decided to keep them as a
  useful shortcut side-effect without calling direct attention to them.
  When and if dependencies on them become an issue it can be addressed
  then.

## Style Guidelines

* Everything through `go fmt` or equiv, no exceptions
* In Vim `set textwidth=72` (not 80 to line numbers fit)
* Use `/* */` for package documentation comment, `//` elsewhere
* Smallest possible names for given scope while still clear
* Favor additional packages (possibly in `internal`) over long names
* Package globals that will be used a lot can be single capital
* Must be good reason to use more than 4 character pkg name
* Avoid unnecessary comments
* Use "deciduous tree" emoji 🌳 to mark Bonzai branches and commands

## Printing, Formatting, and Emphasis

* `Z.Lines(a string) []string`
* `Z.Blocks(a string) []string`

• `Z.Emph(a string) string`   - just emphasis
• `Z.Wrap(a string) string`   - wraps to Z.Columns
• `Z.Indent(a string) string` - indents Z.IndentBy
• `Z.InWrap(a string) string` - indents Z.IndentBy, wraps to Z.Columns
• `Z.Mark(a string) string`   - block aware, indents, wraps all but verbatim

* `Z.Emphf(a string, f ...any) string`       - Emph with Sprintf first
* `Z.Indentf(a string, f ...any) string`       - Indent with Sprintf first
* `Z.Wrapf(a string, f ...any) string`       - Wrap with Sprintf first
* `Z.InWrapf(a string, f ...any) string`     - InWrap with Sprintf first
* `Z.Markf(a string, f ...any) string`  - Mark with Sprintf first
`
• `Z.PrintEmph(a string)`     - shorthand for fmt.Print(Z.Emph(a string))
• `Z.PrintWrap(a string)`     - shorthand for fmt.Print(Z.Wrap(a string))
• `Z.PrintIndent(a string)`     - shorthand for fmt.Print(Z.Indent(a string))
• `Z.PrintInWrap(a string)`     - shorthand for fmt.Print(Z.InWrap(a string))
• `Z.PrintMark(a string)`     - shorthand for fmt.Print(Z.Mark(a string))
`
• `Z.PrintfEmph(a string, f ...any)` - fmt.Print(Z.Emphf(a string, f ...any))
• `Z.PrintfWrap(a string, f ...any)` - fmt.Print(Z.Wrapf(a string, f ...any))
• `Z.PrintfIndent(a string, f ...any)` - fmt.Print(Z.Indentf(a string, f ...any))
• `Z.PrintfInWrap(a string, f ...any)` - fmt.Print(Z.InWrapf(a string, f ...any))
• `Z.PrintfMark(a string, f ...any)` - fmt.Print(Z.Markf(a string, f ...any))

## Acknowledgements

The <https://twitch.tv/rwxrob> community has been constantly involved
with the development of this project, making suggestions about
everything from my use of init, to the name "bonzai". While all their
contributions are too numerous to name specifically, they 
more than deserve a huge thank you here.

## Legal

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
SPDX-License-Identifier: Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the Bonzai™ project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.

## Shwag

Looking to get with
[https://www.trendhunter.com/trends/amuseable-bonsai-tree] to make a
plushy mascot.

## TODO

Here is a list of major things that are still needed before a v1.0
release:

* Complete default `help.Cmd`
* Complete `bonzai` helper command
  * Initialize new branch including optional GitHub repo creation
  * Vet a branch repo for orphan commands, etc.
  * Cache/search list of package depending on `bonzai` as registry
* Complete final art
  * Large logo for landing page
  * Small logo for lists
  * Animated emote suitable for Twitch
* Create a Bonzai™ merch outlet
