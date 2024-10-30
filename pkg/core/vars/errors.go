package vars

import "fmt"

// ----------------------------- NotFound -----------------------------

type NotFound struct {
	Key string
}

func (e NotFound) Error() string {
	return fmt.Sprintf(`could not find key: %v`, e.Key)
}

// ---------------------------- MissingArg ----------------------------

type MissingArg struct {
	Name string
}

func (e MissingArg) Error() string {
	return fmt.Sprintf(`missing argument: %v`, e.Name)
}

// ---------------------------- EmptyKey ----------------------------

type EmptyKey struct{}

func (e EmptyKey) Error() string {
	return fmt.Sprintf(`key name is empty`)
}
