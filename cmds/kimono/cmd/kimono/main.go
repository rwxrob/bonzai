package main

import (
	"os"

	"github.com/rwxrob/bonzai/run"

	"github.com/rwxrob/bonzai/cmds/kimono"
)

func main() {
	if len(os.Getenv(`DEBUG`)) > 0 {
		run.AllowPanic = true
	}
	kimono.Cmd.Run()
}
