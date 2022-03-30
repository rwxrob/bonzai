package bonzai_test

import (
	"fmt"

	"github.com/rwxrob/bonzai"
)

func ExampleEsc() {
	fmt.Println(bonzai.Esc("|&;()<>![]"))
	fmt.Printf("%q", bonzai.Esc(" \n\r"))
	// Output:
	// \|\&\;\(\)\<\>\!\[\]
	// "\\ \\\n\\\r"
}

func ExampleEscAll() {
	list := []string{"so!me", "<here>", "other&"}
	fmt.Println(bonzai.EscAll(list))
	// Output:
	// [so\!me \<here\> other\&]
}
