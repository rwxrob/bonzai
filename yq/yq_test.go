package yq_test

import yq "github.com/rwxrob/bonzai/yq"

func ExampleEvaluate() {
	yq.Evaluate(`.[] | keys`, `testdata/sample.yaml`)
	// Output:
	// # here is a comment
	// - bar
	// - other
}
