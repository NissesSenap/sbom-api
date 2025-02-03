# /bin/sh

kubectl apply -f e2e/minio.yaml
kubectl apply -f e2e/postgres.yaml

# wait for minio and postgres to be ready
kubectl wait --for=condition=ready pod -n minio -l app=minio --timeout=300s
kubectl wait --for=condition=ready pod -n default -l app=postgres --timeout=300s
kubectl wait --for=condition=ready pod -n minio -l job=create-sbom-bucket --timeout=300s

# port fowrard minio and postgres and let the script run in front, when pusing ctrl+c, the port forward will be stopped
# postgresql
kubectl port-forward -n default svc/postgres 5432:5432 &
# minio api
kubectl port-forward -n minio svc/minio 9000:9000 &
# minio console
kubectl port-forward -n minio svc/minio 9001:9001 &
wait
