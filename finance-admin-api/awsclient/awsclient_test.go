package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_S3_ENDPOINT")

	region := "eu-west-1"
	os.Setenv("AWS_REGION", "eu-west-1")

	endpoint := "some-endpoint"
	os.Setenv("AWS_S3_ENDPOINT", endpoint)

	got, err := NewClient(context.Background())

	assert.Nil(t, err)

	assert.IsType(t, new(s3.Client), got)
	assert.Equal(t, region, got.Options().Region)
	assert.Equal(t, endpoint, *got.Options().BaseEndpoint)

}
