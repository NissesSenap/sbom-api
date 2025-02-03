# /bin/sh

kubectl apply -f minio.yaml
kubectl apply -f postgres.yaml
kubectl port-forward -n default svc/postgres 5432:5432 &
kubectl port-forward -n minio svc/minio 9000:9000 &

# generate a s3 bucket on minio called sbom using curl
#curl -X POST http://localhost:9000/minio/webservices/ --data "Action=CreateBucket&Bucket=sbom&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minio%2F20210915%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20210915T000000Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=

# curl -X PUT "http://localhost:9000/sbom" \
#   -H "Authorization: AWS4-HMAC-SHA256 Credential=minioadmin/$(date -u +%Y%m%d)/us-east-1/s3/aws4_request, SignedHeaders=host;x-amz-content-sha256;x-amz-date, Signature=$(echo -n 'PUT\n\n\n\nx-amz-content-sha256:UNSIGNED-PAYLOAD\nx-amz-date:$(date -u +%Y%m%dT%H%M%SZ)\n/sbom' | openssl dgst -sha256 -hmac 'minioadmin' | awk '{print $2}')" \
#   -H "x-amz-content-sha256: UNSIGNED-PAYLOAD" \
#   -H "x-amz-date: $(date -u +%Y%m%dT%H%M%SZ)"
