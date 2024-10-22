package bonzai

// Vars is a map keyed to individual variable keys from Vars. See [Cmd].
type Vars map[string]string

var VarCmd = &Cmd{
	Name: `var`,
}
