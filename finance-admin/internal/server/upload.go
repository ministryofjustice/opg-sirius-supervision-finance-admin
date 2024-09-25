package server

import (
	"errors"
	"fmt"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/api"
	"github.com/opg-sirius-finance-admin/finance-admin/internal/model"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
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

	data, err := shared.NewUpload(shared.ParseReportUploadType(reportUploadType), uploadDate, email, file, handler.Filename)
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
