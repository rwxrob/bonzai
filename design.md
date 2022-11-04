# Style and Design

* Use "deciduous tree" emoji 🌳 to mark Bonzai stuff
* Everything through `go fmt` or equiv, no exceptions
* In Vim `set textwidth=72` (not 80 to line numbers fit)
* Use `/* */` for package documentation comment, `//` elsewhere
* Smallest possible names for given scope while still clear
* Favor additional packages (possibly in `internal`) over long names
* Package globals that will be used a lot can be single capital
* Must be good reason to use more than 4 character pkg name
* Avoid unnecessary comments

## Design Considerations

Here is a summary of the thinking behind important design decisions
when creating Bonzai.

### Promote high-level package library API calls over Cmd bloat

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

### Cmds should be very light

Most Cmds should assign their first-class Call function to one that
lightly wraps a similar function signature in a callable, high-level
library that works entirely independently from the bonzai package.
It's best to promote strong support for sustainable API packages.

### Only bash completion and shell.Cmd planned

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

### Bonzai commands may default to `shell.Cmd` or `help.Cmd`

These provide help information and optional interactive assistance
including tab completion in runtime environments that do not have
`complete -C foo foo` enabled.

### Target users replacing shell scripts in "dot files" collections

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

### Use either `foo.Cmd` or `cmd.Foo` convention

People may decide to put all their Bonzai commands into a single `cmd`
package or to put each command into its own package. Both are
perfectly acceptable and allow the developer making the import to
alias the packages as needed using Go's excellent package import
design.

### Default capital "Z" import name

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

### Promote lower "z" personal Bonzai trees

Just as "dotfiles" has become a convention, use of the simple "z"
GitHub repo should be promoted as a common, "ez" way to find people's
personal monoliths. Obviously, people can use what they want. This is
also consistent with the capital "Z" import name of the `bonzai`
package.

### Use simple regular expressions for usage

Bonzai takes a fundamentally different approach to the command line
and usage documentation. Any command line is a minimal domain specific
language. As such, usage notation simple is not enough. Regular
expressions allow the quick understanding of what is allowed and
should become mandatory learning in a world of minimal domain specific
grammars. Only the most basic regular expressions are required to
produce rich usage strings for any Bonzai command. However, those
wanting traditional usage notation can easily override the
`DefaultUsageFunc` with their own.

### Prioritize output to VT100-compatible terminals

Every modern terminal supports VT100 compatibility. If it doesn't,
people should just not use it. This means emphasis and formatting are
dependent on the [`rwxrob/term`](https://github.com/rwxrob/term)
package for the main output. Bonzai tree developers will likely want
terminal output helpers more than anything (yes even web rendering).

### Secondary output to local web browser

Bonzai will eventually provide the option to enable use of the local
web browser instead of terminal output for help and other
documentation output. This will allow graphic application executables
to simple be double-clicked from graphic desktops. Bonzai will ship
with an embedded web template for such applications, but will also
allow users to granularly customize their own modifications to the
default theme. Bonzai branch creators are encouraged to provide
downloadable themes in this regard. In this way, Bonzai will provide
the web shell to encapsulate other web applications.

### Prefer local web apps over other GUI platforms

While it is obviously possible to create any graphic application with
Bonzai, the creation of localized web applications should be the
preferred convention of the Bonzai GUI applications community. This
simplifies application development and takes aim at bloated
alternatives that embed full GUI web clients (such as Electron) while
still providing a rich terminal interface for those who prefer it.
Every Bonzai binary is its own local web server which Go's rich
standard library to draw on.

### []Section instead of map[string]string for Other

This allows composition notation and ensures the author has control of
how the Other sections appear.

### Move `help.Cmd` into its own package

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

### Leave hidden default command params

When the default command is invoked any of it's params will be
automatically passed as if the command specified them. But they are
not included in the tab completion. This is because there will
inevitably be conflicts between default command params and other
potential completions at that level for other commands. Rather than
disable this, or add tab completion, it was decided to keep them as a
useful shortcut side-effect without calling direct attention to them.
When and if dependencies on them become an issue it can be addressed
then.

### Not exporting Dynamic FuncMap builtins

Since those builtins will land in `mark` subpackage eventually, don't
want to build any dependencies on them now that will break. The
builtins themselves can obviously be used immediately and has a much
smaller chance of changing in the future.

### Custom errors as package globals

When it became clear that there are a number of canned error
conditions that we want to allow Cmd developers to use quickly and
easily, and that not all of them directly relate to a particular Cmd
(`x`) is was decided to move all errors into a central `errors.go`
file as package global custom error structures that implement `error`
interface (`Error() string`) and have concrete, public fields so that
composition notation can be used (instead of forcing an unnecessary
function call). The benefits of custom errors that can be interrogated
with the `errors` functions are obvious and this also consolidates
most English language usage for later when we add full
language-by-locale support.

Note that this was a breaking change (added in `v0.10.0`) with
high-impact because so many have hard-coded `x.UsageError` (which is
why the full `v1.0` is not expected until December 25th 2022, some
months after things should feel settled and finished).

### Move package global Aliases/Shortcuts into Cmd

In Cmd they can be documented. They also tend to be long and eat up
too much of the command line when using them with completion. Better
to just have them in the help docs to lookup when the full path is
wanted.

### Allow dashes, but discourage

Being able to put a command anywhere in the commands tree path of
arguments will never be allowed since it is fundamentally against the
entire raison d'etre of Bonzai. We want people to memorize specific
paths and create the occasional alias and shortcut, not just through
whatever onto the command line hoping it will work. This is why `getopt`
things are so problematic, they all deal with the "end of options and
switch" and the "beginning of arguments" differently. Not Bonzai, where
the arguments are *always* consistent as well as the command paths.
Because of this, dashes should never be used by any Bonzai command or
alias. In fact, the only time a dash will exist is when passing on
arguments to the underlying command shell. For example, when delegating
to `gh` which requires it (although every effort should be made to
remake the entire UI for such things to be Bonzai friendly).

However, the no-getopt approach is so foreign to those training in
decades of this insanity that we need to support dashed aliases (and
other options) that will be intuitive to these people for a very long
time. Of particular note are things like `-h`, `--help`, `-help`, and
`/?` all of which are indelibly marked into most of our muscle memory.
will be intuitive to such people. The addition of support for aliases
that begin with dashes --- but that do not appear in the help
documentation --- is specifically to address this issue. Like all
aliases, the bash shell will dynamically change the `-h` into `help`
when the user taps tab prompting them with the correct word, but, also
like all aliases, if they execute the command with `-h` it will just
work.

Allowing dashed aliases does carry the risk that some will permanently
use them in their scripts. But not allowing them in most help
documentation should be enough to discourage it. Theoretically, it is
entirely possible to create a Bonzai command that does parse its
arguments and parameters using getopt notation. Having this freedom is a
part of FOSS and for those who insist is supported. It is simply
strongly discouraged by the Bonzai project itself which will make
specific design decisions based on the assumptions that dashes are never
used in command names.

### Abandon plans for Bonzai "shell" and double-dash pipes

At one point the thought that allowing a PipeToken to be set would allow
an internal unix-like pipeline construct in the Shortcuts as well as
from the command line, which would easily be expandable into its own
scripting language, but these plans have been abandoned for the sake of
simplicity and anti-bloatware. We want to encourage creation of other
commands that combine the input and output in their own ways and have
added `Input io.Reader` instead. Commands can check this in addition to
checking for `args` to determine what to do with such input.
