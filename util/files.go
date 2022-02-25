package util

import (
	"os"

	"github.com/rwxrob/bonzai/filter"
)

func Files(dir string) []string {
	if dir == "" {
		dir = "."
	}
	files := []string{}
	finfo, _ := os.ReadDir(dir)
	for _, f := range finfo {
		files = append(files, f.Name())
	}
	return files
}

func FilesWith(dir, pre string) []string {
	return filter.HasPrefix(Files(dir), pre)
}
