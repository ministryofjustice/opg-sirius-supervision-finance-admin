package filestorage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"io"
	"os"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	Options() s3.Options
}

type Client struct {
	s3 S3Client
}

func NewClient(ctx context.Context) (*Client, error) {
	awsRegion := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, err
	}

	if iamRole, ok := os.LookupEnv("AWS_IAM_ROLE"); ok {
		client := sts.NewFromConfig(cfg)
		cfg.Credentials = stscreds.NewAssumeRoleProvider(client, iamRole)
	}

	client := s3.NewFromConfig(cfg, func(u *s3.Options) {
		u.UsePathStyle = true
		u.Region = awsRegion

		endpoint := os.Getenv("AWS_S3_ENDPOINT")
		if endpoint != "" {
			u.BaseEndpoint = &endpoint
		}
	})

	return &Client{client}, nil
}

func (c *Client) PutFile(ctx context.Context, bucketName string, fileName string, file io.Reader) (*string, error) {
	output, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               &bucketName,
		Key:                  &fileName,
		Body:                 file,
		ServerSideEncryption: "aws:kms",
		SSEKMSKeyId:          aws.String(os.Getenv("S3_ENCRYPTION_KEY")),
	})

	if output == nil {
		return nil, err
	}

	return output.VersionId, err
}
