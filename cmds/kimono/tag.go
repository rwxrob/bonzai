package kimono

type VerPart int

const (
	Major VerPart = 0
	Minor VerPart = 1
	Patch VerPart = 2
)

// TagBump identifies the current module path, identifies the latest
// version tag and tags the repo with the bumped version.
func TagBump(part VerPart) error {
	return nil
}

// TagList returns the list of tags for the current module.
func TagList() ([]string, error) {
	// out := run.Out(`git`, `tag`, `-l`)
	return nil, nil
}
