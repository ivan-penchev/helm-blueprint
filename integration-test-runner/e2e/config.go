package e2e

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ClusterName     string
	NamespacePrefix string
	MaxParallel     int
	KeepCluster     bool
	RepoRoot        string
}

func LoadConfig() Config {
	return Config{
		ClusterName:     getEnvOrDefault("ITR_CLUSTER_NAME", "helm-blueprint-test"),
		NamespacePrefix: getEnvOrDefault("ITR_NAMESPACE_PREFIX", "hb-it"),
		MaxParallel:     envInt("ITR_MAX_PARALLEL", 4),
		KeepCluster:     envBool("ITR_KEEP_CLUSTER", false),
		RepoRoot:        os.Getenv("ITR_REPO_ROOT"),
	}
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.ClusterName) == "" {
		return fmt.Errorf("cluster name cannot be empty")
	}
	if strings.TrimSpace(c.NamespacePrefix) == "" {
		return fmt.Errorf("namespace prefix cannot be empty")
	}
	if c.MaxParallel < 1 {
		return fmt.Errorf("max parallel must be >= 1")
	}
	return nil
}

func envInt(name string, defaultValue int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return defaultValue
	}
	return parsed
}

func envBool(name string, defaultValue bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(name)))
	if value == "" {
		return defaultValue
	}

	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return defaultValue
	}
}

func getEnvOrDefault(name string, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return defaultValue
	}
	return value
}
