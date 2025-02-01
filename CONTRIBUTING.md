# Contributing

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
