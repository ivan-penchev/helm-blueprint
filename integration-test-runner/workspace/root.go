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
		if isHelmBlueprintRoot(dir) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repository root (expected Chart.yaml and ci/)")
		}
		dir = parent
	}
}

func isHelmBlueprintRoot(dir string) bool {
	return fileExists(filepath.Join(dir, "Chart.yaml")) &&
		fileExists(filepath.Join(dir, "values.yaml")) &&
		dirExists(filepath.Join(dir, "ci"))
}

func fileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

func dirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
