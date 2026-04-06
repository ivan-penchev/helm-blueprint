package cluster

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func CheckPrereqs(commands ...string) error {
	for _, c := range commands {
		if _, err := exec.LookPath(c); err != nil {
			return fmt.Errorf("%s is required but not installed", c)
		}
	}
	return nil
}

func Exists(clusterName string) (bool, error) {
	cmd := exec.Command("kind", "get", "clusters")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("failed to list kind clusters: %w", err)
	}

	for _, line := range strings.Split(out.String(), "\n") {
		if strings.TrimSpace(line) == clusterName {
			return true, nil
		}
	}
	return false, nil
}

func Create(clusterName string, wait time.Duration) error {
	cmd := exec.Command("kind", "create", "cluster", "--name", clusterName, "--wait", wait.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Delete(clusterName string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
