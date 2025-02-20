package server

import (
	"errors"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/apierror"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/shared"
	"net/http"
)

const (
	downloadError = "Sorry, this link has expired. Please request a new report."
	systemError   = "Sorry, there is a problem with the service. Please try again later."
)

type GetDownloadVars struct {
	Uid          string
	Filename     string
	ErrorMessage string
	AppVars
}

type GetDownloadHandler struct {
	router
}

func (h *GetDownloadHandler) render(v AppVars, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	uid := r.URL.Query().Get("uid")

	var downloadRequest shared.DownloadRequest
	data := GetDownloadVars{AppVars: v}

	err := downloadRequest.Decode(uid)
	if err != nil {
		data.ErrorMessage = downloadError
	} else {
		err = h.Client().CheckDownload(ctx, uid)
		if err != nil {
			var notFound apierror.NotFound
			if errors.As(err, &notFound) {
				data.ErrorMessage = downloadError
			} else {
				data.ErrorMessage = systemError
			}
		} else {
			data.Uid = uid
			data.Filename = downloadRequest.Key
		}
	}

	return h.execute(w, r, data)
}
