package github_test

import (
	"fmt"

	"github.com/rwxrob/bonzai/github"
)

func ExampleClient_defaults() {
	gh := github.NewClient()
	fmt.Println(gh.Host())
	fmt.Println(gh.APIVersion())
	fmt.Println(gh.API(`rwxrob/web`))
	// Output:
	// github.com
	// v3
	// https://api.github.com/rwxrob/web
}

func ExampleClient_API() {
	gh := github.NewClient()
	fmt.Println(gh.API(`rwxrob/web`))
	gh.SetHost(`github.example.com`)
	fmt.Println(gh.API(`rwxrob/web`))
	// Output:
	// https://api.github.com/rwxrob/web
	// https://github.example.com/api/v3/rwxrob/web
}

func ExampleClient_Repo() {
	gh := github.NewClient()
	repo, err := gh.Repo(`rwxrob/web`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(repo["name"])
	// Output:
	// web
}

func ExampleRepo() {
	repo, err := github.Repo(`rwxrob/z`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(repo["name"])
	// Output:
	// z
}

func ExampleLatest() {
	ver, err := github.Latest(`docker/compose`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ver[0:2])
	// Output:
	// v2
}

/*
func ExampleLatest_private() {
	ver, err := github.Latest(`rwxrob/testprivate`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ver[0:2])
	// Output:
	// v0
}
*/
