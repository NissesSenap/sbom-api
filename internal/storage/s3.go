package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	client *s3.Client
}

func NewS3Storage(endpoint, region, accessKey, secretKey string) (*S3Storage, error) {
	// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
	cfg, err := config.LoadDefaultConfig(context.TODO())
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
		o.BaseEndpoint = &endpoint
		o.Region = "us-east-1"
		o.UsePathStyle = true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &S3Storage{client: client}, nil
}

func (s *S3Storage) Upload(ctx context.Context, bucket, key, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
