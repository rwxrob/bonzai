package choose

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var DefaultPrompt = `#? `

func readline() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

type Chooser[T any] interface {
	Choose() (T, error)
}

type Choices[T any] []T

// Choose prompts terminal user to pick from one of the choices and
// returns the index of the choice with the choice or an empty value and
// error.
func (c Choices[T]) Choose() (int, T, error) {
	var empty T
	width := strconv.Itoa(len(strconv.Itoa(len(c) + 1)))
	for i, v := range c {
		fmt.Printf("%"+width+"v. %v\n", i+1, v)
	}
	for {
		fmt.Print(DefaultPrompt)
		resp := readline()
		if resp == "q" {
			return -1, empty, nil
		}
		n, _ := strconv.Atoi(resp)
		if 0 < n && n < len(c)+1 {
			return n - 1, c[n-1], nil
		}
	}
}

// From prompts terminal user to pick from one of the choices using the
// DefaultChooser.
func From[T any](choices []T) (int, T, error) {
	return Choices[T](choices).Choose()
}
