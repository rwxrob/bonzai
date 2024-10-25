package vars

import "fmt"

type NotFound struct {
	Key string
}

func (e NotFound) Error() string {
	return fmt.Sprintf(`could not find key: %v`, e)
}
