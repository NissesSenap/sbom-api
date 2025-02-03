# Contributing

## Linting

In the project [golangci-lint](https://golangci-lint.run/) is used, and you can run it locally before creating a PR.
To see which version of golangci-lint to use, see [pull.yaml](.github/workflows/pull.yaml)

## e2e

```shell
# create kind cluster
kind create cluster --config e2e/kind.yaml
# create postgresql instance
k apply -f e2e/postgres.yaml
```

Don't know why, but I need to do port-forward...

```shell
kubectl port-forward svc/postgres 5432:5432
```

Using psql

```shell
psql -h localhost -d sbom -U sbom
```

## Database managment

After some back and forth I went with [sqlc](https://github.com/sqlc-dev/sqlc) to generate the queries in a type safe way.
I use pgx for connection pool maangment.

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

In the long run I will setup a [CI/CD](https://docs.sqlc.dev/en/stable/howto/ci-cd.html) solution to verify that it's up to date.
