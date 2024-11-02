package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/rwxrob/ssh"
	"gopkg.in/yaml.v3"
)

func main() {

	if !(len(os.Args) > 1) || os.Args[1] == `-h` {
		log.SetFlags(0)
		log.Fatal(`usage: runonany ARG ... < STDIN`)
	}

	// read the command and standard input to the binary
	cmdline := strings.Join(os.Args[1:], ` `)
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// create and configure a controller with three clients

	confpath := `/etc/runonany.yaml`
	if tmp := os.Getenv(`RUNONANYCONF`); len(tmp) > 0 {
		confpath = tmp
	}
	ctl := new(ssh.Controller)
	byt, err := os.ReadFile(confpath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(byt, ctl)
	if err != nil {
		log.Fatal(err)
	}

	// override central ssh host

	if tmp := os.Getenv(`RUNONANY_TARGET`); len(tmp) > 0 {
		for _, c := range ctl.Clients {
			c.Host.Addr = tmp
		}
	}

	// connect and report on client connection status

	ctl.Connect()
	ctl.LogStatus()

	// run it and capture the output

	stdout, stderr, err := ctl.RunOnAny(cmdline, stdin)
	fmt.Println(`STDOUT --------------`)
	fmt.Println(stdout)
	fmt.Println(`STDERR --------------`)
	fmt.Println(stderr)
	fmt.Println(`ERR --------------`)
	fmt.Println(err)
}
