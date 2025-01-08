package api

import (
	"context"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/event"
	"io"
	"net/http"
)

type MockDispatch struct {
	event any
}

func (m *MockDispatch) FinanceAdminUpload(ctx context.Context, event event.FinanceAdminUpload) error {
	m.event = event
	return nil
}

type MockFileStorage struct {
	versionId  string
	bucketname string
	filename   string
	file       io.Reader
	err        error
}

func (m *MockFileStorage) PutFile(ctx context.Context, bucketName string, fileName string, file io.Reader) (*string, error) {
	m.bucketname = bucketName
	m.filename = fileName
	m.file = file

	return &m.versionId, m.err
}

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetDoFunc fetches the mock client's `Do` func. Implement this within a test to modify the client's behaviour.
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
