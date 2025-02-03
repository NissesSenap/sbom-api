# Contributing

## Linting

In the project [golangci-lint](https://golangci-lint.run/) is used, and you can run it locally before creating a PR.
To see which version of golangci-lint to use, see [pull.yaml](.github/workflows/pull.yaml)

## e2e

```shell
# create kind cluster
kind create cluster --config e2e/kind.yaml
# create postgresql instance
sh e2e/run.sh
```

The script will start port-forwads to minio and postgresql.

### Postgres

Using psql

```shell
psql -h localhost -d sbom -U sbom
```

### Minio

If the app is unable to reach the `sbom` bucket, verify that the Kubernetes job has run as inteded.

## Database managment

After some back and forth I went with [sqlc](https://github.com/sqlc-dev/sqlc) to generate the queries in a type safe way.
I use pgx for connection pool maangment.

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

In the long run I will setup a [CI/CD](https://docs.sqlc.dev/en/stable/howto/ci-cd.html) solution to verify that it's up to date.
