package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

func FindRepoRoot() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current directory: %w", err)
	}

	dir := start
	for {
		target := filepath.Join(dir, "ci", "kind", "test.sh")
		if st, err := os.Stat(target); err == nil && !st.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root (expected ci/kind/test.sh)")
		}
		dir = parent
	}
}
