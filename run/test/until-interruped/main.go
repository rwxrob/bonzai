package main

import (
	"time"

	"github.com/rwxrob/bonzai/run"
)

func main() {
	run.UntilInterrupted()
	for {
		time.Sleep(1 * time.Millisecond)
	}
}
