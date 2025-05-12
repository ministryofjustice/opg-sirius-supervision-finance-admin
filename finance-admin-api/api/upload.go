package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/event"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"
)

const s3Directory = "finance-admin"

func validateCSVHeaders(file []byte, reportUploadType shared.ReportUploadType, useStrictComparison bool) error {
	fileReader := bytes.NewReader(file)
	csvReader := csv.NewReader(fileReader)

	expectedHeaders := reportUploadType.CSVHeaders()

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
		cleanedHeader := cleanString(header)
		if cleanedHeader == "" {
			continue
		}
		if i >= len(expectedHeaders) {
			return apierror.ValidationError{Errors: apierror.ValidationErrors{
				"FileUpload": {
					"incorrect-headers": "CSV headers do not match for the report trying to be uploaded",
				},
			},
			}
		}
		if useStrictComparison {
			if cleanString(readHeaders[i]) != cleanString(expectedHeaders[i]) {
				return apierror.ValidationError{Errors: apierror.ValidationErrors{
					"FileUpload": {
						"incorrect-headers": "CSV headers do not match for the report trying to be uploaded",
					},
				},
				}
			}
		} else {
			if !strings.Contains(cleanString(readHeaders[i]), cleanString(expectedHeaders[i])) {
				return apierror.ValidationError{Errors: apierror.ValidationErrors{
					"FileUpload": {
						"incorrect-headers": "CSV headers do not match for the report trying to be uploaded",
					},
				},
				}
			}
		}
	}

	_, err = fileReader.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	return nil
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)

	// Replace double-spaces in headers with single spaces (BACS uploads have double spaces)
	s = strings.ReplaceAll(s, "  ", " ")

	s = strings.ToLower(s)

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

	if upload.ReportUploadType.HasHeader() {
		useStrictHeaderComparison := upload.ReportUploadType != shared.ReportTypeUploadPaymentsSupervisionCheque

		err := validateCSVHeaders(upload.File, upload.ReportUploadType, useStrictHeaderComparison)
		if err != nil {
			return err
		}
	}

	_, err := s.filestorage.PutFile(
		ctx,
		os.Getenv("ASYNC_S3_BUCKET"),
		fmt.Sprintf("%s/%s", s3Directory, upload.Filename),
		bytes.NewReader(upload.File))

	if err != nil {
		return err
	}

	uploadEvent := event.FinanceAdminUpload{
		EmailAddress: upload.Email,
		Filename:     fmt.Sprintf("%s/%s", s3Directory, upload.Filename),
		UploadType:   upload.ReportUploadType.Key(),
		UploadDate:   upload.UploadDate,
		PisNumber:    upload.PisNumber,
	}
	err = s.dispatch.FinanceAdminUpload(ctx, uploadEvent)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	return nil
}
