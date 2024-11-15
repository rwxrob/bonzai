package kimono

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

type node struct {
	name         string
	dependencies []*node
	dependents   []*node
}

type graph struct {
	nodes map[string]*node
}

func newGraph() *graph {
	return &graph{nodes: make(map[string]*node)}
}

func (g *graph) addNode(name string) {
	if _, ok := g.nodes[name]; !ok {
		g.nodes[name] = &node{name: name}
	}
}

func (g *graph) addDependency(from, to string) {
	fromNode := g.nodes[from]
	toNode := g.nodes[to]
	fromNode.dependencies = append(fromNode.dependencies, toNode)
	toNode.dependents = append(toNode.dependents, fromNode)
}

func (g *graph) getDependencies(name string) []*node {
	return g.nodes[name].dependencies
}

func (g *graph) getDependents(name string) []*node {
	return g.nodes[name].dependents
}

func getDependents(g *graph, name string) []*node {
	// Strip version information and search for dependents
	strippedName := stripVersion(name)
	var dependents []*node
	for _, node := range g.nodes {
		if stripVersion(node.name) == strippedName {
			dependents = append(
				dependents,
				node.dependents...)
		}
	}
	return dependents
}

func stripVersion(name string) string {
	parts := strings.Split(name, "@")
	return parts[0]
}

func (g *graph) hasDependencies(name string) bool {
	return len(g.nodes[name].dependencies) > 0
}

// ListDependents returns a list of all Go modules that depend on the
// current module.
func ListDependents() ([]string, error) {
	modName := strings.TrimSpace(run.Out("go", "list", "-m"))
	graph, err := dependencyGraph()
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s@%s", modName, latestTag())
	depNodes := getDependents(graph, name)
	dependents := make([]string, 0)
	for _, dep := range depNodes {
		if dep.name == modName {
			continue
		}
		dependents = append(dependents, dep.name)
	}
	return dependents, nil
}

// ListDependencies returns a list of all Go modules that the current
// module depends on.
func ListDependencies() ([]string, error) {
	modName := strings.TrimSpace(run.Out("go", "list", "-m"))
	graph, err := dependencyGraph()
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%s@%s", modName, latestTag())
	depNodes := graph.getDependencies(name)
	dependencies := make([]string, 0)
	for _, dep := range depNodes {
		dependencies = append(dependencies, dep.name)
	}
	return dependencies, nil
}

func dependencyGraph() (*graph, error) {
	graph := newGraph()
	root, err := futil.HereOrAbove(".git")
	if err != nil {
		return nil, err
	}
	ppwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filepath.WalkDir(
		filepath.Dir(root),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			if !futil.Exists(
				filepath.Join(path, "go.mod"),
			) {
				return nil
			}
			if err := os.Chdir(path); err != nil {
				return err
			}
			err = os.Chdir(path)
			if err != nil {
				return nil
			}
			modName := strings.TrimSpace(
				run.Out("go", "list", "-m"),
			)
			modName = fmt.Sprint(modName, "@", latestTag())
			graph.addNode(modName)
			dependencies := to.Lines(run.Out("go", "list", "-m", "all"))
			if len(dependencies) < 2 {
				return nil
			}
			dependencies = dependencies[1:]
			for _, dep := range dependencies {
				dep = strings.TrimSpace(dep)
				parts := strings.Split(dep, " ")
				if len(parts) < 2 {
					continue
				}
				name := parts[0]
				ver := parts[1]
				dep = fmt.Sprint(name, "@", ver)
				graph.addNode(dep)
				graph.addDependency(modName, dep)
			}
			return nil
		},
	)
	err = os.Chdir(ppwd)
	if err != nil {
		return nil, err
	}
	return graph, nil
}
