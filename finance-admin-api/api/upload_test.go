package api

import (
	"bytes"
	"encoding/json"
	"github.com/opg-sirius-finance-admin/shared"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_upload(t *testing.T) {
	var b bytes.Buffer

	uploadForm := &shared.Upload{
		ReportUploadType: "Test",
		Email:            "joseph@test.com",
		Filename:         "file.txt",
		File:             []byte("file contents"),
	}

	_ = json.NewEncoder(&b).Encode(uploadForm)
	req := httptest.NewRequest(http.MethodPost, "/uploads", &b)
	w := httptest.NewRecorder()

	server := Server{}
	_ = server.upload(w, req)

	res := w.Result()
	defer res.Body.Close()

	expected := ""

	assert.Equal(t, expected, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
}
