package server

import (
	"errors"
	"github.com/opg-sirius-finance-admin/apierror"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDownload(t *testing.T) {
	appVars := AppVars{Path: "/download"}
	tests := []struct {
		name         string
		uid          string
		mockError    error
		expectedVars GetDownloadVars
	}{
		{
			name: "successful download",
			uid:  "eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=",
			expectedVars: GetDownloadVars{
				Uid:      "eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=",
				Filename: "test.csv",
				AppVars:  appVars,
			},
		},
		{
			name: "invalid uid",
			uid:  "invalid-uid",
			expectedVars: GetDownloadVars{
				ErrorMessage: downloadError,
				AppVars:      appVars,
			},
		},
		{
			name:      "download not found",
			uid:       "eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=",
			mockError: apierror.NotFound{},
			expectedVars: GetDownloadVars{
				ErrorMessage: downloadError,
				AppVars:      appVars,
			},
		},
		{
			name:      "system error",
			uid:       "eyJLZXkiOiJ0ZXN0LmNzdiIsIlZlcnNpb25JZCI6InZwckF4c1l0TFZzYjVQOUhfcUhlTlVpVTlNQm5QTmN6In0=",
			mockError: errors.New("internal server error"),
			expectedVars: GetDownloadVars{
				ErrorMessage: systemError,
				AppVars:      appVars,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mockApiClient{error: tt.mockError}
			ro := &mockRoute{client: client}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "download?uid="+tt.uid, nil)

			sut := GetDownloadHandler{ro}
			err := sut.render(appVars, w, r)

			assert.Nil(t, err)
			assert.True(t, ro.executed)

			assert.Equal(t, tt.expectedVars, ro.data)
		})
	}
}
