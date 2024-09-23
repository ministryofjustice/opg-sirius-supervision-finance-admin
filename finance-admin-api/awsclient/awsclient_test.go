package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		region     *string
		role       *string
		wantRegion string
		wantErr    bool
	}{
		{nil, nil, "eu-west-1", false},
		{aws.String(""), nil, "eu-west-1", false},
		{aws.String("us-west-1"), nil, "us-west-1", false},
		{nil, aws.String("some-iam-role"), "eu-west-1", false},
	}

	for _, test := range tests {
		os.Unsetenv("AWS_REGION")
		os.Unsetenv("AWS_IAM_ROLE")

		if test.region != nil {
			os.Setenv("AWS_REGION", *test.region)
		}
		if test.role != nil {
			os.Setenv("AWS_IAM_ROLE", *test.role)
		}

		got, err := NewClient(context.Background())
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
		}
		assert.IsType(t, new(s3.Client), got)
		assert.Equal(t, test.wantRegion, got.Options().Region)
	}
}
