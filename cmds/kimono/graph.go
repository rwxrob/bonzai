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
	depNodes := graph.getDependents(modName)
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
	depNodes := graph.getDependencies(modName)
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
	filepath.WalkDir(
		filepath.Dir(root),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return nil
			}
			if d.Name() == ".git" || d.Name() == "vendor" {
				return filepath.SkipDir
			}
			if !futil.Exists(filepath.Join(path, "go.mod")) {
				return nil
			}
			if err := os.Chdir(path); err != nil {
				return err
			}
			modName := strings.TrimSpace(run.Out("go", "list", "-m"))
			graph.addNode(modName)
			dependencies := run.Out("go", "list", "-m", "all")
			for _, dep := range strings.Split(dependencies, "\n") {
				dep = strings.TrimSpace(dep)
				graph.addNode(dep)
				graph.addDependency(modName, dep)
			}
			return nil
		},
	)
	return graph, nil
}

func main() {
	graph, err := dependencyGraph()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Hello", graph)
	for modName, node := range graph.nodes {
		fmt.Printf("Module: %s\n", modName)

		// Print dependencies
		fmt.Print("  Dependencies: ")
		if len(node.dependencies) == 0 {
			fmt.Println("None")
		} else {
			for _, dep := range node.dependencies {
				fmt.Printf("%s ", dep.name)
			}
			fmt.Println()
		}

		// Print dependents
		fmt.Print("  Dependents: ")
		if len(node.dependents) == 0 {
			fmt.Println("None")
		} else {
			for _, dep := range node.dependents {
				fmt.Printf("%s ", dep.name)
			}
			fmt.Println()
		}

		fmt.Println()
	}
}
