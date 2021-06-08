# kenv

A command line tool for extracting environment variables from K8S's `ConfigMap` and `Secret` configuration. This configuration is created by calling `kustomize` tool.

## Setup

```sh
# Install dependencies
go mod tidy

# Build the tool
go build -ldflags '-s -w' -o kenv ./main/

# Run the tests
ginkgo -tags unitTests -randomizeAllSpecs -failFast -r .

# Run the tool
./kenv prepare -i waas-config/environments/be-gcw1/pp -k $KEY -o dotenv
```

## Contributing

If you have suggestions for how `kenv` could be improved, or want to report a bug, open an issue! We'd love all and any contributions.

