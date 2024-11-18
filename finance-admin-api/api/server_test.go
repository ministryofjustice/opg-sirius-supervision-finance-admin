package api

import (
	"context"
	"github.com/opg-sirius-finance-admin/finance-admin-api/event"
	"github.com/opg-sirius-finance-admin/finance-admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
)

type IntegrationSuite struct {
	suite.Suite
	testDB *testhelpers.TestDatabase
}

type MockDispatch struct {
	event any
}

func (m *MockDispatch) FinanceAdminUpload(ctx context.Context, event event.FinanceAdminUpload) error {
	m.event = event
	return nil
}

type MockFileStorage struct {
	bucketname string
	filename   string
	file       io.Reader
}

func (m *MockFileStorage) PutFile(ctx context.Context, bucketName string, fileName string, file io.Reader) error {
	m.bucketname = bucketName
	m.filename = fileName
	m.file = file

	return nil
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
