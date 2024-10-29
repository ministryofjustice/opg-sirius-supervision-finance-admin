package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDownload(t *testing.T) {
	client := mockApiClient{}
	ro := &mockRoute{client: client}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "download?uid=dGVzdC5jc3Y=", nil)

	appVars := AppVars{}

	sut := GetDownloadHandler{ro}
	err := sut.render(appVars, w, r)

	assert.Nil(t, err)
	assert.True(t, ro.executed)

	expected := GetDownloadVars{
		Uid:      "dGVzdC5jc3Y=",
		Filename: "test.csv",
		AppVars:  appVars,
	}
	assert.Equal(t, expected, ro.data)
}
