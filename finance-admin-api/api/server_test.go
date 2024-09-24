package api

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockAWSClient struct{}

func (m *MockAWSClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return nil, nil
}
