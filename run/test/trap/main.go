package main

import (
	"syscall"
	"time"

	"github.com/rwxrob/bonzai/run"
)

func main() {
	run.Trap(nil, syscall.SIGINT, syscall.SIGTERM)
	for {
		time.Sleep(1 * time.Millisecond)
	}
}
