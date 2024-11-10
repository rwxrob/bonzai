# Design decisions

Bonzai takes issue with several traditional software design practices and ideas about human computer interaction --- especially about command line interfaces. As such, several design decisions need some explanation to help people understand why we do the things the way we do.

## Dashes it command-line arguments are the devil

Sometimes we joke about being racist against dashes, or "dashists" (without hopefully trivialising the problem of real racism in the world). We really are not. We actually prefer "kebab-case" when it comes to compound names in configuration files. But on the command line, the "getopts" approach is ancient tech developed with entirely different requirements and, frankly, is an anti-pattern today.  Just look at some of the most modern commands and you will see those that support dashed options are either providing non-dash alternatives or giving them up all together in favor of stateful commands and command tree monoliths often in a multicall binary form (think BusyBox).

## Not showing aliases under `Usage`

Aliases are up to the developer to document since they are generally just to help with shortcuts when completion is not allowed, can be hidden, can be optionally completed, and often are just old copies of deprecated ways to refer to the newly named command. One possibility is to add an `# Aliases` section to Cmd.Long.

## Preference for nouns first with leaf commands

When given the following options, the traditional object-oriented idea of putting the noun first should be considered (but not mandated).

```
kimono deps list
kimono list deps
```

This keeps command trees the most organized since many, many things could have a `list` command.
