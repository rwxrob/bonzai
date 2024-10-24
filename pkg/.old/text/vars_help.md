cache variables in {{ execachedir "vars"}}

The {{aka}} command provides a cross-platform, persistent alternative to environment/system variables. The subcommands are designed to be safe and convenient.

Variables are stored as key=val (property) pairs, one to a line, in the 
{{ execachedir "vars" }} file.

Key names are automatically prefixed with the Cmd.Path ({{ .Path }} in this case) which changes depending on where this Bonzai branch is composed into your command tree.

Keys must not include an equal sign (=) which is the only line delimiter.

Carriage returns (\r) and line returns (\n) are escaped and each line is terminated with a line return (\n).
