package util_test

import (
	"fmt"
	"net/http"
	ht "net/http/httptest"

	"github.com/rwxrob/bonzai/util"
)

func ExampleCompareUpdated() {

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

	fmt.Println(util.CompareUpdated(20220322080542, older.URL))
	fmt.Println(util.CompareUpdated(20220322080542, newer.URL))
	fmt.Println(util.CompareUpdated(20220322080542, same.URL))
	fmt.Println(util.CompareUpdated(20220322080542, "foobar"))

	// Output:
	// -1
	// 1
	// 0
	// -2
}

func ExampleCompareVersions() {

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `"v0.0.1"`)
		})
	older := ht.NewServer(handler)
	defer older.Close()

	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `"v0.1.0"`)
		})
	newer := ht.NewServer(handler)
	defer newer.Close()

	handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `"v0.0.2"`)
		})
	same := ht.NewServer(handler)
	defer same.Close()

	fmt.Println(util.CompareVersions(`v0.0.2`, older.URL))
	fmt.Println(util.CompareVersions(`v0.0.2`, newer.URL))
	fmt.Println(util.CompareVersions(`v0.0.2`, same.URL))
	fmt.Println(util.CompareVersions(`v0.0.2`, "foobar"))

	// Output:
	// 1
	// -1
	// 0
	// -2
}
