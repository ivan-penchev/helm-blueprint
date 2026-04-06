# integration-test-runner

Go binary and Testify e2e suite for local/CI integration tests.

## Structure

```text
integration-test-runner
├── e2e
│   ├── cases.go
│   ├── config.go
│   ├── helpers.go
│   └── suite_test.go
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
cd integration-test-runner
go run .
```

Or from inside integration-test-runner:

```bash
go run .
```

Supported flags:

- --keep-cluster
- --cluster-name <name>
- --namespace-prefix <name>
- --max-parallel <n>

Unknown flags/args are forwarded to go test for [integration-test-runner/e2e](integration-test-runner/e2e).

Each testcase uses its own namespace (`<namespace-prefix>-<test>-<suffix>`), which enables parallel-safe isolation.
Cases run in parallel (bounded by `--max-parallel`, default `4`).

If your machine does not have `gcc`, run direct `go test` commands with `CGO_ENABLED=0`.
