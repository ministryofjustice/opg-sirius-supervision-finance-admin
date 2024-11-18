package filestorage

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type mockS3Client struct {
	headObjectOutput *s3.HeadObjectOutput
	headObjectError  error
}

func (m *mockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m.headObjectOutput, m.headObjectError
}

func (m *mockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return nil, nil
}

func (m *mockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return nil, nil
}

func (m *mockS3Client) Options() s3.Options {
	return s3.Options{}
}

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

func TestFileExists(t *testing.T) {
	tests := []struct {
		name string
		mock *mockS3Client
		want bool
	}{
		{
			name: "success",
			mock: &mockS3Client{
				headObjectOutput: &s3.HeadObjectOutput{},
				headObjectError:  nil,
			},
			want: true,
		},
		{
			name: "fail",
			mock: &mockS3Client{
				headObjectOutput: nil,
				headObjectError:  errors.New("error"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{s3: tt.mock}
			got := client.FileExists(context.Background(), "bucket", "filename", "versionID")
			assert.Equal(t, tt.want, got)
		})
	}
}
