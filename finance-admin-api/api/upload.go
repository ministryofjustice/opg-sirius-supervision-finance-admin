package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/opg-sirius-finance-admin/shared"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"unicode"
)

func validateCSVHeaders(file []byte, reportUploadType string) error {
	fileReader := bytes.NewReader(file)
	csvReader := csv.NewReader(fileReader)
	expectedHeaders := reportHeadersByType(reportUploadType)

	readHeaders, err := csvReader.Read()
	if err != nil {
		return apierror.ValidationError{Errors: apierror.ValidationErrors{
			"FileUpload": {
				"read-failed": "Failed to read CSV headers",
			},
		},
		}
	}

	for i, header := range readHeaders {
		readHeaders[i] = cleanString(header)
	}

	fmt.Println(expectedHeaders)
	fmt.Println(readHeaders)

	// Compare the extracted headers with the expected headers
	if !reflect.DeepEqual(readHeaders, expectedHeaders) {
		return apierror.ValidationError{Errors: apierror.ValidationErrors{
			"FileUpload": {
				"incorrect-headers": "CSV headers do not match for the report trying to be uploaded",
			},
		},
		}
	}

	_, err = fileReader.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func reportHeadersByType(reportType string) []string {
	switch reportType {
	case shared.ReportTypeUploadDeputySchedule.Key():
		return []string{"Deputy number", "Deputy name", "Case number", "Client forename", "Client surname", "Do not invoice", "Total outstanding"}
	case shared.ReportTypeUploadDebtChase.Key():
		return []string{"Client_no", "Deputy_name", "Total_debt"}
	case shared.ReportTypeUploadPaymentsOPGBACS.Key():
		return []string{"Line", "Type", "Code", "Number", "Transaction", "Value Date", "Amount", "Amount Reconciled", "Charges", "Status", "Desc Flex", "Consolidated line"}
	default:
		return []string{"Unknown report type"}
	}
}

func cleanString(s string) string {
	// Trim leading and trailing spaces
	s = strings.TrimSpace(s)
	// Remove non-printable characters
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, s)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var upload shared.Upload
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&upload); err != nil {
		return err
	}

	err := validateCSVHeaders(upload.File, upload.ReportUploadType)
	if err != nil {
		return err
	}

	fmt.Println("before upload")

	_, err = s.awsClient.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("ASYNC_S3_BUCKET")),
		Key:                  aws.String(fmt.Sprintf("%s/%s", "finance-admin", upload.Filename)),
		Body:                 bytes.NewReader(upload.File),
		ServerSideEncryption: "AES256",
	})

	fmt.Println("after upload")

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}
