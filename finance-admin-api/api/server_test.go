package api

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockAWSClient struct {
	incomingObject *s3.PutObjectInput
	outgoingObject *s3.GetObjectOutput
	optFns         []func(*s3.Options)
	err            error
}

func (m *MockAWSClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.incomingObject = params
	m.optFns = optFns
	return nil, m.err
}

func (m *MockAWSClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.outgoingObject, m.err
}

func (m *MockAWSClient) Options() s3.Options {
	return s3.Options{}
}
