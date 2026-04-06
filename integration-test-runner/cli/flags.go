package cli

import "fmt"

type Options struct {
	ClusterName string
	KeepCluster bool
	SkipTier1   bool
	SkipTier2   bool
	ExtraArgs   []string
}

func Parse(args []string) (Options, error) {
	opts := Options{
		ClusterName: "helm-blueprint-test",
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--keep-cluster":
			opts.KeepCluster = true
		case "--skip-tier1":
			opts.SkipTier1 = true
		case "--skip-tier2":
			opts.SkipTier2 = true
		case "--cluster-name":
			if i+1 >= len(args) {
				return Options{}, fmt.Errorf("missing value for --cluster-name")
			}
			i++
			opts.ClusterName = args[i]
		default:
			opts.ExtraArgs = append(opts.ExtraArgs, args[i])
		}
	}

	return opts, nil
}
