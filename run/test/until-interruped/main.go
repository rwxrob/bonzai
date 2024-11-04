package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/rwxrob/bonzai/run"
)

func main() {
	run.Trap(
		func() { fmt.Println("interrupted") },
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	for {
		time.Sleep(1 * time.Millisecond)
	}
}
