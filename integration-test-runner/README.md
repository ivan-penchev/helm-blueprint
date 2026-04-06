# integration-test-runner

Go binary replacement for ci/kind/run-local.sh.

## Structure

```text
integration-test-runner
├── main.go
├── cli
│   └── flags.go
├── cluster
│   └── kind.go
├── crd
│   └── install.go
├── logging
│   └── logger.go
├── testrun
│   └── execute.go
├── workspace
│   └── root.go
├── go.mod
└── README.md
```

## Usage

From repo root:

```bash
go run ./integration-test-runner
```

Or from inside integration-test-runner:

```bash
go run .
```

Supported flags:

- --keep-cluster
- --skip-tier1
- --skip-tier2
- --cluster-name <name>

Unknown flags/args are forwarded to ci/kind/test.sh.
