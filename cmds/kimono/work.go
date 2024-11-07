package kimono

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/futil"
)

func WorkOn() error {
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return err
	}
	return filepath.WalkDir(filepath.Dir(root), workOnWalkDirFn)
}

func workOnWalkDirFn(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return nil
	}
	if d.Name() == ".git" || d.Name() == "vendor" {
		return filepath.SkipDir
	}
	if !futil.Exists(filepath.Join(path, "go.work")) {
		return nil
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	_ = os.Rename("go.work.off", "go.work")
	return nil
}

func WorkOff() error {
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return err
	}
	return filepath.WalkDir(filepath.Dir(root), workOffWalkDirFn)
}

func workOffWalkDirFn(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return nil
	}
	if d.Name() == ".git" || d.Name() == "vendor" {
		return filepath.SkipDir
	}
	if !futil.Exists(filepath.Join(path, "go.work")) {
		return nil
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	_ = os.Rename("go.work", "go.work.off")
	return nil
}
