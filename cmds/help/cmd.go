package help

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/mark"
	"gopkg.in/yaml.v3"
)

//go:embed style.yaml
var Style []byte

var Cmd = &bonzai.Cmd{
	Name:  `help`,
	Alias: `h|-h|--help|--h|/?`,
	Vers:  `v0.8.0`,
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

		// load embedded yaml file and convert to json
		styleMap := map[string]any{}
		if err := yaml.Unmarshal(Style, &styleMap); err != nil {
			return err
		}
		jsonBytes, err := json.Marshal(styleMap)
		if err != nil {
			return err
		}

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithPreservedNewLines(),
			glamour.WithStylesFromJSONBytes(jsonBytes),
		)

		if err != nil {
			return fmt.Errorf("developer-error: %v", err)
		}

		rendered, err := renderer.Render(md)
		if err != nil {
			return fmt.Errorf("developer-error: %v", err)
		}

		fmt.Println("\u001b[2J\u001b[H" + rendered)

		return nil
	},
}
