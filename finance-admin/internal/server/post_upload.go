package server

import (
	"errors"
	"fmt"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"strings"
)

type UploadFormHandler struct {
	router
}

func (h *UploadFormHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := getContext(r)

	reportUploadType := shared.ParseReportUploadType(r.PostFormValue("reportUploadType"))
	uploadDate := r.PostFormValue("uploadDate")
	email := r.PostFormValue("email")

	// Handle file upload
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		return h.handleError(w, r, "No file uploaded", http.StatusBadRequest)
	}
	defer file.Close()

	expectedFilename, err := reportUploadType.Filename(uploadDate)
	if err != nil {
		return h.handleError(w, r, "Could not parse upload date", http.StatusBadRequest)
	}

	if handler.Filename != expectedFilename && expectedFilename != "" {
		expectedFilename := strings.Replace(expectedFilename, ":", "/", -1)
		return h.handleError(w, r, fmt.Sprintf("Filename should be named \"%s\"", expectedFilename), http.StatusBadRequest)
	}

	data, err := shared.NewUpload(reportUploadType, uploadDate, email, file, handler.Filename)
	if err != nil {
		return h.handleError(w, r, "Failed to read file", http.StatusBadRequest)
	}

	// Upload the file
	if err := h.Client().Upload(ctx, data); err != nil {
		return h.handleUploadError(w, r, err)
	}

	w.Header().Add("HX-Redirect", fmt.Sprintf("%s/uploads?success=upload", v.EnvironmentVars.Prefix))

	return nil
}

// handleError simplifies repetitive error handling in the render method.
func (h *UploadFormHandler) handleError(w http.ResponseWriter, r *http.Request, msg string, code int) error {
	fileError := model.ValidationErrors{
		"FileUpload": map[string]string{"required": msg},
	}
	data := AppVars{ValidationErrors: RenameErrors(fileError)}
	w.WriteHeader(code)
	return h.execute(w, r, data)
}

// handleUploadError processes specific upload-related errors.
func (h *UploadFormHandler) handleUploadError(w http.ResponseWriter, r *http.Request, err error) error {
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
