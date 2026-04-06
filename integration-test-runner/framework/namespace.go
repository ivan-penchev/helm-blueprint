package framework

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var nonDNS1123 = regexp.MustCompile(`[^a-z0-9-]`)

func NewNamespaceName(prefix string, testName string) string {
	basePrefix := sanitizeDNS1123(prefix)
	if basePrefix == "" {
		basePrefix = "hb-it"
	}

	testPart := sanitizeDNS1123(testName)
	if testPart == "" {
		testPart = "case"
	}

	suffix := fmt.Sprintf("%x", time.Now().UnixNano())
	if len(suffix) > 8 {
		suffix = suffix[len(suffix)-8:]
	}

	name := fmt.Sprintf("%s-%s-%s", basePrefix, testPart, suffix)
	if len(name) > 63 {
		name = name[:63]
	}
	name = strings.Trim(name, "-")
	if name == "" {
		return "hb-it-case"
	}
	return name
}

func CreateNamespace(ctx context.Context, repoRoot string, namespace string) error {
	_, err := RunCommand(ctx, repoRoot, nil, "kubectl", "create", "namespace", namespace)
	return err
}

func DeleteNamespace(ctx context.Context, repoRoot string, namespace string) error {
	_, err := RunCommand(ctx, repoRoot, nil, "kubectl", "delete", "namespace", namespace, "--ignore-not-found=true", "--wait=true", "--timeout=120s")
	return err
}

func sanitizeDNS1123(value string) string {
	v := strings.ToLower(strings.TrimSpace(value))
	v = strings.ReplaceAll(v, "_", "-")
	v = strings.ReplaceAll(v, "/", "-")
	v = strings.ReplaceAll(v, ".", "-")
	v = nonDNS1123.ReplaceAllString(v, "-")
	v = strings.Trim(v, "-")
	for strings.Contains(v, "--") {
		v = strings.ReplaceAll(v, "--", "-")
	}
	return v
}
