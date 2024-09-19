package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

type Session struct {
	AwsSession *session.Session
}

func NewSession() (*Session, error) {
	awsRegion, ok := os.LookupEnv("AWS_REGION")

	if !ok {
		awsRegion = "eu-west-1" // default region
	}

	sess, err := session.NewSession(&aws.Config{Region: &awsRegion})
	if err != nil {
		return nil, err
	}

	if iamRole, ok := os.LookupEnv("AWS_IAM_ROLE"); ok {
		sess.Config.Credentials = stscreds.NewCredentials(sess, iamRole)
	}

	return &Session{sess}, nil
}
