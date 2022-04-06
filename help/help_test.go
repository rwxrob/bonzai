package help_test

import (
	"github.com/rwxrob/bonzai/help"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/term"
)

func ExampleCmd_name() {
	term.AttrOff()
	x := &Z.Cmd{
		Name: "foo",
	}
	help.ForTerminal(x, "name")
	// Output:
	// foo
}

func ExampleCmd_summary() {
	term.AttrOff()
	x := &Z.Cmd{
		Summary: `foo all the things`,
	}
	help.ForTerminal(x, "summary")
	// Output:
	// foo all the things
}

func ExampleCmd_other() {
	term.AttrOff()
	x := &Z.Cmd{
		Other: []Z.Section{{"foo", `some foo text`}},
	}
	help.ForTerminal(x, "foo")
	// Output:
	// some foo text
}

func ExampleCmd_all_Error_No_Call_No_Commands() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:    `foo`,
		Version: "v0.1.0",
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        foo
	//
	// SYNOPSIS
	//        {ERROR: neither Call nor Commands defined}

}

func ExampleCmd_all_Error_Params_No_Call() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:   `foo`,
		Params: []string{"p1", "p2"},
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        foo
	//
	// SYNOPSIS
	//        {ERROR: Params without Call: p1, p2}
}

func ExampleCmd_all_Call_No_Params() {
	term.AttrOff()
	x := &Z.Cmd{
		Name: `foo`,
		Call: func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        foo
	//
	// SYNOPSIS
	//        foo
}

func ExampleCmd_all_Call_With_Optional_Params() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:   `foo`,
		Params: []string{"p1", "p2"},
		Call:   func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        foo
	//
	// SYNOPSIS
	//        foo (p1|p2)?
}

func ExampleCmd_all_Call_With_Optional_Params_and_Commands() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:   `cmd`,
		Params: []string{"p1", "p2"},
		Commands: []*Z.Cmd{
			&Z.Cmd{
				Name:    "foo",
				Aliases: []string{"f"},
				Summary: "foo the things",
			},
			&Z.Cmd{
				Name:    "bar",
				Summary: "bar the things",
			},
			&Z.Cmd{
				Name: "nosum",
			},
		},
		Call: func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        cmd
	//
	// SYNOPSIS
	//        cmd (COMMAND|(p1|p2)?)
	//
	// COMMANDS
	//        f|foo - foo the things
	//        bar   - bar the things
	//        nosum
}

func ExampleCmd_all_Legal_Copyright_Only() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:      `cmd`,
		Copyright: "Copyright @2022 Rob",
		Call:      func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        cmd
	//
	// SYNOPSIS
	//        cmd
	//
	// LEGAL
	//        cmd Copyright @2022 Rob
}

func ExampleCmd_all_Legal_Copyright_and_Version() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:      `cmd`,
		Copyright: "Copyright @2022 Rob",
		Version:   "v0.0.1",
		Call:      func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        cmd
	//
	// SYNOPSIS
	//        cmd
	//
	// LEGAL
	//        cmd (v0.0.1) Copyright @2022 Rob
}

func ExampleCmd_all_Legal_Copyright_and_Version_and_License() {
	term.AttrOff()
	x := &Z.Cmd{
		Name:      `cmd`,
		Copyright: "Copyright @2022 Rob",
		Version:   "v0.0.1",
		License:   "Apache 2.0",
		Call:      func(_ *Z.Cmd, _ ...string) error { return nil },
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        cmd
	//
	// SYNOPSIS
	//        cmd
	//
	// LEGAL
	//        cmd (v0.0.1) Copyright @2022 Rob
	//        License Apache 2.0
}

func ExampleCmd_all_Other() {
	term.AttrOff()
	x := &Z.Cmd{
		Name: `cmd`,
		Call: func(_ *Z.Cmd, _ ...string) error { return nil },
		Other: []Z.Section{
			{"FOO", "A whole section dedicated to foo."},
			{"bar", "WTF is a bar anyway"},
		},
	}
	help.ForTerminal(x, "all")
	// Output:
	// NAME
	//        cmd
	//
	// SYNOPSIS
	//        cmd
	//
	// FOO
	//        A whole section dedicated to foo.
	//
	// BAR
	//        WTF is a bar anyway
}
