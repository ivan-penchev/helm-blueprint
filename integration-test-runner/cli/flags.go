package cli

import (
	"fmt"
	"strconv"
)

type Options struct {
	ClusterName string
	NamespacePrefix string
	MaxParallel int
	KeepCluster bool
	ExtraArgs   []string
}

func Parse(args []string) (Options, error) {
	opts := Options{
		ClusterName: "helm-blueprint-test",
		NamespacePrefix: "hb-it",
		MaxParallel: 4,
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--keep-cluster":
			opts.KeepCluster = true
		case "--cluster-name":
			if i+1 >= len(args) {
				return Options{}, fmt.Errorf("missing value for --cluster-name")
			}
			i++
			opts.ClusterName = args[i]
		case "--namespace-prefix":
			if i+1 >= len(args) {
				return Options{}, fmt.Errorf("missing value for --namespace-prefix")
			}
			i++
			opts.NamespacePrefix = args[i]
		case "--max-parallel":
			if i+1 >= len(args) {
				return Options{}, fmt.Errorf("missing value for --max-parallel")
			}
			i++
			parsed, err := strconv.Atoi(args[i])
			if err != nil || parsed < 1 {
				return Options{}, fmt.Errorf("invalid value for --max-parallel: %s", args[i])
			}
			opts.MaxParallel = parsed
		default:
			opts.ExtraArgs = append(opts.ExtraArgs, args[i])
		}
	}

	return opts, nil
}
