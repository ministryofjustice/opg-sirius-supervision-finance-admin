package server

import (
	"encoding/csv"
	"errors"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"io"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

type UploadHandler struct {
	router
}

func (h *UploadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	reportUploadType := r.PostFormValue("reportUploadType")
	uploadDate := r.PostFormValue("uploadDate")
	email := r.PostFormValue("email")

	// Handle file upload
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		return h.handleError(w, r, "No file uploaded", http.StatusBadRequest)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	expectedHeaders := reportHeadersByType(reportUploadType)

	readHeaders, err := csvReader.Read()
	if err != nil {
		return h.handleError(w, r, "Failed to read CSV headers", http.StatusBadRequest)
	}

	for i, header := range readHeaders {
		readHeaders[i] = cleanString(header)
	}

	// Compare the extracted headers with the expected headers
	if reportUploadType != model.ReportTypeTest.Key() && !reflect.DeepEqual(readHeaders, expectedHeaders) {
		return h.handleError(w, r, "CSV headers do not match for the file trying to be uploaded", http.StatusBadRequest)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	data, err := shared.NewUpload(reportUploadType, uploadDate, email, file, handler.Filename)
	if err != nil {
		return h.handleError(w, r, "Failed to read file", http.StatusBadRequest)
	}

	// Upload the file
	if err := h.Client().Upload(ctx, data); err != nil {
		return h.handleUploadError(w, r, err)
	}

	return nil
}

func reportHeadersByType(reportType string) []string {
	switch reportType {
	case model.ReportTypeUploadDeputySchedule.Key():
		return []string{"Deputy number", "Deputy name", "Case number", "Client forename", "Client surname", "Do not invoice", "Total outstanding"}
	case model.ReportTypeUploadDebtChase.Key():
		return []string{"Client_no", "Deputy_name", "Total_debt"}
	case model.ReportTypeUploadPaymentsOPGBACS.Key():
		return []string{"Line", "Type", "Code", "Number", "Transaction", "Value Date", "Amount", "Amount Reconciled", "Charges", "Status", "Desc Flex", "Consolidated line"}
	default:
		return []string{"Unknown report type"}
	}
}

// handleError simplifies repetitive error handling in the render method.
func (h *UploadHandler) handleError(w http.ResponseWriter, r *http.Request, msg string, code int) error {
	fileError := model.ValidationErrors{
		"FileUpload": map[string]string{"required": msg},
	}
	data := AppVars{ValidationErrors: RenameErrors(fileError)}
	w.WriteHeader(code)
	return h.execute(w, r, data)
}

// handleUploadError processes specific upload-related errors.
func (h *UploadHandler) handleUploadError(w http.ResponseWriter, r *http.Request, err error) error {
	var (
		valErr model.ValidationError
		stErr  api.StatusError
	)
	if errors.As(err, &valErr) {
		data := AppVars{ValidationErrors: RenameErrors(valErr.Errors)}
		w.WriteHeader(http.StatusUnprocessableEntity)
		return h.execute(w, r, data)
	} else if errors.As(err, &stErr) {
		data := AppVars{Error: stErr.Error()}
		w.WriteHeader(stErr.Code)
		return h.execute(w, r, data)
	}

	return err
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
