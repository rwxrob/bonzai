package vars

import "fmt"

// ----------------------------- NotFound -----------------------------

type ErrNotFound struct {
	Key string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf(`could not find key: %v`, e.Key)
}

// ---------------------------- MissingArg ----------------------------

type ErrMissingArg struct {
	Name string
}

func (e ErrMissingArg) Error() string {
	return fmt.Sprintf(`missing argument: %v`, e.Name)
}

// ---------------------------- EmptyKey ----------------------------

type ErrEmptyKey struct{}

func (e ErrEmptyKey) Error() string {
	return fmt.Sprintf(`key name is empty`)
}
