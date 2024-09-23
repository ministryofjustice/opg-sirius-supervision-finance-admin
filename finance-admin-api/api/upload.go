package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/opg-sirius-finance-admin/finance-admin-api/awsclient"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
)

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	awsClient, err := awsclient.NewClient(ctx)

	if err != nil {
		return err
	}

	uploader := manager.NewUploader(awsClient)

	var upload shared.Upload
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		return err
	}

	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("ASYNC_S3_BUCKET")),
		Key:                  aws.String(fmt.Sprintf("%s/%s", "finance-admin", upload.Filename)),
		Body:                 bytes.NewReader(upload.File),
		ServerSideEncryption: "AES256",
	})

	if err != nil {
		return err
	}

	return nil
}
