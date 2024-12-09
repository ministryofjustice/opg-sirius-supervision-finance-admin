package api

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/db"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/event"
	"github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin-api/testhelpers"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"os"
	"testing"
)

type IntegrationSuite struct {
	suite.Suite
	cm     *testhelpers.ContainerManager
	seeder *testhelpers.Seeder
	ctx    context.Context
}

func (suite *IntegrationSuite) SetupSuite() {
	suite.ctx = telemetry.ContextWithLogger(context.Background(), telemetry.NewLogger("finance-api-test"))
	suite.cm = testhelpers.NewContainerManager(suite.ctx)
}

func (suite *IntegrationSuite) SetupTest() {
	suite.seeder = testhelpers.NewSeeder(suite.ctx)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}

func (suite *IntegrationSuite) TearDownSuite() {
	suite.cm.TearDown(suite.ctx)
}

func (suite *IntegrationSuite) AfterTest(suiteName, testName string) {
	err := suite.cm.Restore(suite.ctx)
	if err != nil {
		suite.T().Error(fmt.Sprintf("Failed to restore snapshot after test %s", testName))
	}
}

type MockDispatch struct {
	event any
}

func (m *MockDispatch) FinanceAdminUpload(ctx context.Context, event event.FinanceAdminUpload) error {
	m.event = event
	return nil
}

type MockFileStorage struct {
	versionId      string
	bucketname     string
	filename       string
	file           io.Reader
	outgoingObject *s3.GetObjectOutput
	err            error
	exists         bool
}

func (m *MockFileStorage) GetFile(ctx context.Context, bucketName string, fileName string, versionId string) (*s3.GetObjectOutput, error) {
	return m.outgoingObject, m.err
}

func (m *MockFileStorage) PutFile(ctx context.Context, bucketName string, fileName string, file io.Reader) (*string, error) {
	m.bucketname = bucketName
	m.filename = fileName
	m.file = file

	return &m.versionId, nil
}

// add a FileExists method to the MockFileStorage struct
func (m *MockFileStorage) FileExists(ctx context.Context, bucketName string, filename string, versionID string) bool {
	return m.exists
}

type MockReports struct {
	query db.ReportQuery
}

func (m *MockReports) Generate(ctx context.Context, filename string, query db.ReportQuery) (*os.File, error) {
	m.query = query
	return nil, nil
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
