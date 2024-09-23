package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"os"
)

type AwsClient struct {
	Client s3.Client
}

func NewClient(ctx context.Context) (*s3.Client, error) {
	awsRegion, ok := os.LookupEnv("AWS_REGION")

	if !ok || awsRegion == "" {
		awsRegion = "eu-west-1" // default region
	}

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
		u.BaseEndpoint = aws.String(os.Getenv("AWS_S3_ENDPOINT"))
	})

	return client, nil
}
