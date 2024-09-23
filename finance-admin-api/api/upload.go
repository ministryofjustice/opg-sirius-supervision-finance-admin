package api

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/opg-sirius-finance-admin/finance-admin-api/session"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	var upload shared.Upload
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		return err
	}

	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	sess.AwsSession.Config.Endpoint = &endpoint
	sess.AwsSession.Config.S3ForcePathStyle = aws.Bool(true)

	uploader := s3manager.NewUploader(sess.AwsSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(os.Getenv("ASYNC_S3_BUCKET")),
		Key:                  &upload.Filename,
		Body:                 bytes.NewReader(upload.File),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}

	return nil
}
