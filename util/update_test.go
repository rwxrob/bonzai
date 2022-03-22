package util_test

import (
	"fmt"
	"net/http"
	ht "net/http/httptest"

	"github.com/rwxrob/bonzai/util"
)

func ExampleCheckUpdated() {

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `20220322080540`)
		})
	older := ht.NewServer(handler)
	defer older.Close()

	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `20220322080546`)
		})
	newer := ht.NewServer(handler)
	defer newer.Close()

	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `20220322080542`)
		})
	same := ht.NewServer(handler)
	defer same.Close()

	fmt.Println(util.NeedsUpdate(20220322080542, older.URL))
	fmt.Println(util.NeedsUpdate(20220322080542, newer.URL))
	fmt.Println(util.NeedsUpdate(20220322080542, same.URL))
	fmt.Println(util.NeedsUpdate(20220322080542, "foobar"))

	// Output:
	// 0
	// 1
	// 0
	// -1
}
