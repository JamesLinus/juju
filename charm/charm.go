package charm

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// The Charm interface is implemented by any type that
// may be handled as a charm.
type Charm interface {
	Meta() *Meta
	Config() *Config
	Revision() int
}

// Read reads a Charm from path, which can point to either a charm bundle or a
// charm directory.
func Read(path string) (Charm, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return ReadDir(path)
	}
	return ReadBundle(path)
}

// InferRepository returns a charm repository inferred from
// the provided URL. Local URLs will use the provided path.
func InferRepository(curl *URL, localRepoPath string) (repo Repository, err error) {
	switch curl.Schema {
	case "cs":
		repo = Store()
	case "local":
		if localRepoPath == "" {
			return nil, errors.New("path to local repository not specified")
		}
		repo = &LocalRepository{localRepoPath}
	default:
		return nil, fmt.Errorf("unknown schema for charm URL %q", curl)
	}
	return
}

// validatePath is used when reading to ensure that charms contain no unwanted
// files.
func validatePath(p string) error {
	if strings.HasPrefix(path.Clean(p), "hooks/juju-") {
		return fmt.Errorf("reserved hook name: %q", p)
	}
	return nil
}
