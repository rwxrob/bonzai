package scan_test

import (
	"log"
	"os"

	"github.com/rwxrob/bonzai/scan"
)

func ExampleCur() {
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.Flags())
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	m := new(scan.Cur)
	m.Print()
	m.NewLine()
	m.Print()
	m.Log()
	//Output:
	// U+0000 '\x00' 0,0-0 (0-1)
	// U+0000 '\x00' 1,1-1 (0-1)
	// U+0000 '\x00' 1,1-1 (0-1)
}
