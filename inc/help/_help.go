package help

/*
var Help = &Z.Cmd{
	Name:    `help`,
	Aliases: []string{"h"},
	Call: func(x *Z.Cmd, args ...string) error {

		var buf string

		buf += util.Emph("**NAME**", 0, -1) + "\n       " + x.Title() + "\n\n"
		buf += util.Emph("**SYNOPSIS**", 0, -1) + "\n       " + x.Name + " " + x.Usage + "\n\n"

		if len(x.Commands.M) > 0 {
			buf += util.Emph("**COMMANDS**", 0, -1) + "\n" + x.Titles(7, 20) + "\n\n"
		}

		if len(x.Description) > 0 {
			buf +=
				util.Emph("**DESCRIPTION**", 0, -1) + "\n" +
					util.Emph(x.Description, 7, 65) + "\n\n"
		}

		if x.Source != "" || x.Issues != "" || x.Site != "" {

			buf += util.Emph("**LINKS**", 0, -1) + "\n"

			if x.Site != "" {
				buf += "       Site:   " + x.Site + "\n"
			}

			if x.Source != "" {
				buf += "       Source: " + x.Source + "\n"
			}

			if x.Issues != "" {
				buf += "       Issues: " + x.Issues + "\n"
			}

			buf += "\n"

		}

		if x.Copyright != "" {
			buf += util.Emph("**LEGAL**", 0, -1) + "\n" + util.Indent(x.Legal(), 7) + "\n\n"
		}

		return buf

	},
}
*/
