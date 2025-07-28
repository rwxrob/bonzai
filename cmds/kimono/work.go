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

func WorkToggleRecursive(root, from, to string) error {
	return filepath.WalkDir(
		filepath.Dir(root),
		renameRecursive(from, to),
	)
}

func WorkToggleModule(from, to string) error {
	path, err := futil.HereOrAbove(`go.mod`)
	if err != nil {
		return fmt.Errorf(`not inside a module`)
	}
	path = filepath.Dir(path)
	if !futil.Exists(from) {
		return fmt.Errorf(
			`%s does not exist in current module`,
			filepath.Base(from),
		)
	}
	return os.Rename(from, to)
}

func renameRecursive(from, to string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".git" {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == from {
			err := os.Rename(
				path,
				filepath.Join(filepath.Dir(path), to),
			)
			if err != nil {
				fmt.Println(err)
			}
		}
		return nil
	}
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
