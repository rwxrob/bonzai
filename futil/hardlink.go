//go:build !windows
// +build !windows

package futil

import (
	"fmt"
	"os"
	"syscall"
)

// IsHardLink attempts to determine if the file at the end of path is
// a unix/linux hard link by counting its number of links. On Windows
// always returns false.
func IsHardLink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false, fmt.Errorf("unable to retrieve file information")
	}
	return stat.Nlink > 1, nil
}
