package contentfile

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3FileReader interface {
	GetObject(bucket, key string) (io.Reader, error)
}

type S3Client struct {
}

func NewS3Client() S3Client {
	return S3Client{}
}

func (s S3Client) GetObject(bucket, key string) (io.Reader, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	if object, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return nil, err
	} else {
		return object.Body, nil
	}
}
