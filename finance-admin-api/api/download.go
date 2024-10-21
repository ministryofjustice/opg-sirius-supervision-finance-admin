package api

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/opg-sirius-finance-admin/apierror"
	"io"
	"net/http"
	"os"
)

func (s *Server) download(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	filename := r.PathValue("filename")

	result, err := s.awsClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("REPORTS_S3_BUCKET")),
		Key:    aws.String(filename),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "NoSuchKey" {
				return apierror.NotFoundError(err)
			}
		}
		return fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", *result.ContentType)

	// Stream the S3 object to the response writer using io.Copy
	_, err = io.Copy(w, result.Body)

	return err
}
