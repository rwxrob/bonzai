package kimono

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
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
	if d.Name() == ".git" || d.Name() == "vendor" {
		return filepath.SkipDir
	}
	if d.IsDir() {
		return nil
	}
	if d.Name() == "go.work.off" {
		_ = os.Rename(
			path,
			filepath.Join(filepath.Dir(path), "go.work"),
		)
	}
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

func WorkInit(args ...string) error {
	if err := run.Exec("go", "work", "init"); err != nil {
		return err
	}
	if err := run.Exec("go", "work", "use", "."); err != nil {
		return err
	}
	for _, arg := range args {
		fmt.Println("using", arg)
		if err := run.Exec("go", "work", "use", arg); err != nil {
			return err
		}
	}
	return nil
}

func WorkGenerate() error {
	deps, error := dependencyGraph()
	if error != nil {
		return error
	}
	modName := strings.TrimSpace(run.Out("go", "list", "-m"))
	name := fmt.Sprintf("%s@%s", modName, latestTag())
	if err := run.Exec("go", "work", "init"); err != nil {
		return err
	}
	if err := run.Exec("go", "work", "use", "."); err != nil {
		return err
	}
	for _, dep := range deps.getDependencies(name) {
		name := stripVersion(dep.name)
		rel, err := getRelPath(modName, name)
		if err != nil {
			continue
		}
		fmt.Println("using", rel)
		if err := run.Exec("go", "work", "use", rel); err != nil {
			return err
		}
	}
	return run.Exec("go", "work", "sync")
}

func getRelPath(modName, depName string) (string, error) {
	modRepo := strings.Split(modName, "/")[:2]
	depRepo := strings.Split(depName, "/")[:2]
	if modRepo[0] != depRepo[0] || modRepo[1] != depRepo[1] {
		return "", fmt.Errorf("dependencies in different repos")
	}
	modPathLen := len(strings.Split(modName, "/"))
	depPathLen := len(strings.Split(depName, "/"))
	rel, err := filepath.Rel(modName, depName)
	if err != nil {
		return "", err
	}
	relPathLen := len(strings.Split(rel, "/"))
	if relPathLen == modPathLen+depPathLen {
		return "", fmt.Errorf("dependencies outside monorepo")
	}
	return "./" + rel, nil
}
