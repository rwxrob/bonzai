# Bonzai, Well-Tended CLI Command Trees for Gophers

Bonzai commands are unlike anything you've probably encountered so far,
no getopt dashes, no ugly commander interface to learn, no 12637 lines
of shell tab completion bloat to source before it will work, just well
manicured commands organized into rooted node trees for your
command-line, nested-tab-completing enjoyment with built-in (optional)
help documentation that can easily be marshalled into JSON and presented
either as text or well-styled HTML. Think "readthedocs" but served
locally from your Bonzai command-line apps themselves.

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

## Example GitHub Template

<https://github.com/rwxrob/template-bonzai>

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

## Legal 

Copyright 2022 Robert S. Muhlestein (<mailto:rob@rwx.gg>)  
Licensed under Apache-2.0

"Bonzai" and "bonzai" are legal trademarks of Robert S. Muhlestein but
can be used freely to refer to the cmdbox project
<https://github.com/rwxrob/bonzai> without limitation. To avoid
potential developer confusion, intentionally using these trademarks to
refer to other projects --- free or proprietary --- is prohibited.
