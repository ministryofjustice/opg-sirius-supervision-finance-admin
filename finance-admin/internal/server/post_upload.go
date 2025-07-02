package server

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type UploadFormHandler struct {
	router
}

func (h *UploadFormHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var err error
	var pisNumber int
	uploadType := shared.ParseUploadType(r.PostFormValue("uploadType"))

	if uploadType == shared.ReportTypeUploadPaymentsSupervisionCheque {
		pisNumberForm := strings.ReplaceAll(r.PostFormValue("pisNumber"), " ", "")
		pisNumber, err = strconv.Atoi(pisNumberForm)
		if len([]rune(pisNumberForm)) != 6 {
			return h.handleError(w, r, "PisNumber", "PIS number must be 6 digits", http.StatusUnprocessableEntity)
		}
		if err != nil {
			return h.handleError(w, r, "PisNumber", "Error parsing PIS number", http.StatusUnprocessableEntity)
		}
	}

	uploadDate := r.PostFormValue("uploadDate")
	email := r.PostFormValue("email")

	// Handle file upload
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		return h.handleError(w, r, "FileUpload", "No file uploaded", http.StatusUnprocessableEntity)
	}

	defer file.Close()

	var expectedFilename string
	if uploadDate != "" || uploadType.NoDateRequired() {
		expectedFilename, err = uploadType.Filename(uploadDate)
		if err != nil {
			return h.handleError(w, r, "UploadDate", "Could not parse upload date", http.StatusUnprocessableEntity)
		}
	} else {
		return h.handleError(w, r, "UploadDate", "Upload date required", http.StatusUnprocessableEntity)
	}

	if !uploadType.NoDateRequired() && shared.NewDate(uploadDate).After(shared.Date{Time: time.Now()}) {
		fileError := model.ValidationErrors{
			"UploadDate": map[string]string{"date-in-the-future": "Can not upload for a date in the future"},
		}
		data := AppVars{ValidationErrors: RenameErrors(fileError)}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	}

	if expectedFilename != "" && !matchFilenameWithWildcard(handler.Filename, expectedFilename) {
		expectedFilename := strings.ReplaceAll(expectedFilename, ":", "/")
		return h.handleError(w, r, "FileUpload", fmt.Sprintf("Filename should be \"%s\"", expectedFilename), http.StatusUnprocessableEntity)
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		return h.handleError(w, r, "FileUpload", "Failed to read file", http.StatusUnprocessableEntity)
	}

	if ok, field, reason := validateCSVHeaders(fileData, uploadType); !ok {
		fileError := model.ValidationErrors{
			"FileUpload": map[string]string{field: reason},
		}
		data := AppVars{ValidationErrors: RenameErrors(fileError)}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	}

	data := shared.Upload{
		UploadType:   uploadType,
		EmailAddress: email,
		Base64Data:   base64.StdEncoding.EncodeToString(fileData),
		Filename:     handler.Filename,
		UploadDate:   shared.NewDate(uploadDate),
		PisNumber:    pisNumber,
	}

	// Upload the file
	if err := h.Client().Upload(ctx, data); err != nil {
		return err
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("%s/uploads?success=upload", v.EnvironmentVars.Prefix))

	return nil
}

// handleError simplifies repetitive error handling in the render method.
func (h *UploadFormHandler) handleError(w http.ResponseWriter, r *http.Request, field string, msg string, code int) error {
	fileError := model.ValidationErrors{
		field: map[string]string{"required": msg},
	}
	data := AppVars{ValidationErrors: RenameErrors(fileError)}
	w.WriteHeader(code)
	return h.execute(w, r, data)
}

func validateCSVHeaders(file []byte, uploadType shared.ReportUploadType) (ok bool, field string, reason string) {
	if !uploadType.HasHeader() {
		return true, "", ""
	}

	fileReader := bytes.NewReader(file)
	csvReader := csv.NewReader(fileReader)

	expectedHeaders := uploadType.CSVHeaders()

	readHeaders, err := csvReader.Read()
	if err != nil {
		return false, "read-failed", "Failed to read CSV headers"
	}

	for i, header := range readHeaders {
		cleanedHeader := cleanString(header)
		if cleanedHeader == "" {
			continue
		}
		if i >= len(expectedHeaders) {
			if uploadType.HasOptionalExtraHeaders() {
				continue
			}
			return false, "incorrect-headers", "CSV headers do not match for the file being uploaded"
		}
		if uploadType.StrictHeaderComparison() {
			if cleanString(readHeaders[i]) != cleanString(expectedHeaders[i]) {
				return false, "incorrect-headers", "CSV headers do not match for the file being uploaded"
			}
		} else {
			if !strings.Contains(cleanString(readHeaders[i]), cleanString(expectedHeaders[i])) {
				return false, "incorrect-headers", "CSV headers do not match for the file being uploaded"
			}
		}
	}

	return true, "", ""
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

func matchFilenameWithWildcard(actualFilename, expectedFilename string) bool {
	// If no wildcard, do a direct comparison
	if !strings.Contains(expectedFilename, "*") {
		return actualFilename == expectedFilename
	}

	if strings.Count(expectedFilename, "*") == 1 {
		parts := strings.Split(expectedFilename, "*")
		prefix := parts[0]
		suffix := parts[1]
		return strings.HasPrefix(actualFilename, prefix) &&
			strings.HasSuffix(actualFilename, suffix)
	}

	return false
}
