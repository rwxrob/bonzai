# Design decisions

## Not showing aliases under `Usage`

Aliases are up to the developer to document since they are generally just to help with shortcuts when completion is not allowed, can be hidden, can be optionally completed, and often are just old copies of deprecated ways to refer to the newly named command. One possibility is to add an `# Aliases` section to Cmd.Long.

## Preference for nouns first with leaf commands

When given the following options, the traditional object-oriented idea of putting the noun first should be considered (but not mandated).

```
kimono deps list
kimono list deps
```

This keeps command trees the most organized since many, many things could have a `list` command.
