package filestorage

import (
	"context"
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

	assert.IsType(t, new(Client), got)
	assert.Equal(t, region, got.s3.Options().Region)
	assert.Equal(t, endpoint, *got.s3.Options().BaseEndpoint)

}
