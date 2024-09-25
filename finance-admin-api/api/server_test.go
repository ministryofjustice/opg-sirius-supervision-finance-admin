package api

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockAWSClient struct {
	params *s3.PutObjectInput
	optFns []func(*s3.Options)
}

func (m *MockAWSClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.params = params
	m.optFns = optFns
	return nil, nil
}

func (m *MockAWSClient) Options() s3.Options {
	return s3.Options{}
}
