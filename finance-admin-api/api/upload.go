package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opg-sirius-finance-admin/finance-admin-api/session"
	"net/http"
	"os"
	"strings"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	sess, err := session.NewSession()

	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	sess.AwsSession.Config.Endpoint = &endpoint
	sess.AwsSession.Config.S3ForcePathStyle = aws.Bool(true)

	uploader := s3manager.NewUploader(sess.AwsSession)

	file := struct {
		FileName     string
		FileContents string
	}{
		"Test.txt",
		"Test Contents",
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("async-upload"),
		Key:    &file.FileName,
		Body:   strings.NewReader(file.FileContents),
	})
	if err != nil {
		return err
	}

	return nil
}
