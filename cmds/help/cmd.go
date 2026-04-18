package help

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"github.com/rwxrob/bonzai/term"
	"github.com/rwxrob/bonzai/to"
)

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Alias: `h|-h|--help|--h|/?`,
	Vers:  `v0.9.0`,
	Short: `display command help`,
	Long: `
		The {{code .Name}} command displays the help information for the
		immediate previous command unless it is passed arguments, in which
		case it resolves the arguments as if they were passed to the
		previous command and the help for the leaf command is displayed
		instead.`,

	Do: func(x *bonzai.Cmd, args ...string) (err error) {

		if len(args) > 0 {
			x, args, err = x.Caller().SeekInit(args...)
		} else {
			x = x.Caller()
		}

		md, err := mark.Bonzai(x)
		if err != nil {
			return err
		}

		term.WinSizeUpdate()
		width := int(term.WinSize.Col)
		if width < 10 {
			width = 80
		}

		rendered := renderMarkdown(md, width)
		fmt.Print("\033[2J\033[H" + rendered)

		return nil
	},
}

// renderMarkdown performs minimal markdown rendering for terminal
// output: section headers become bold uppercase, code blocks are
// preserved as-is, and paragraphs are word-wrapped to the given width.
func renderMarkdown(md string, width int) string {
	lines := strings.Split(md, "\n")
	var out strings.Builder
	var inCode bool

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// fenced code blocks (~~~ or ```)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") || strings.HasPrefix(trimmed, "```") {
			if inCode {
				inCode = false
				out.WriteString("\n")
				continue
			}
			inCode = true
			out.WriteString("\n")
			continue
		}

		if inCode {
			out.WriteString("    " + line + "\n")
			continue
		}

		// indented code blocks (4+ spaces)
		if len(line) > 0 && strings.HasPrefix(line, "    ") {
			out.WriteString(line + "\n")
			continue
		}

		// section headers
		if strings.HasPrefix(trimmed, "#") {
			title := strings.TrimLeft(trimmed, "# ")
			out.WriteString(term.Bold + strings.ToUpper(title) + term.Reset + "\n\n")
			continue
		}

		// blank lines
		if trimmed == "" {
			out.WriteString("\n")
			continue
		}

		// regular paragraph: collect contiguous non-blank lines
		var para strings.Builder
		para.WriteString(line)
		for i+1 < len(lines) {
			next := lines[i+1]
			nextTrimmed := strings.TrimSpace(next)
			if nextTrimmed == "" ||
				strings.HasPrefix(nextTrimmed, "#") ||
				strings.HasPrefix(nextTrimmed, "~~~") ||
				strings.HasPrefix(nextTrimmed, "```") ||
				strings.HasPrefix(next, "    ") {
				break
			}
			i++
			para.WriteString(" " + next)
		}

		wrapped, _ := to.Wrapped(width, para.String())
		out.WriteString(wrapped + "\n")
	}

	return out.String()
}
