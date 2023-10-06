package openbarapi

import (
	"testing"

	"github.com/cocktailrobots/openbar-server/pkg/hardware"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type testSuite struct {
	*test.DBSuite
	Api *OpenBarAPI
}

func TestOpenBarApiPackage(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	dbSuite := test.NewDBSuite("openbardb", "test", "../../../schema/openbardb/", "")
	hw := hardware.NewTestHardware(8)
	api := New(logger, dbSuite, mux.NewRouter(), hw)
	suite.Run(t, &testSuite{
		DBSuite: dbSuite,
		Api:     api,
	})
}

// BeforeTest is called before each test. It calls resetToHash to reset the database to the
// hash of the last commit on the branch.
func (s *testSuite) BeforeTest(suiteName, testName string) {
	s.DBSuite.BeforeTest(suiteName, testName)
	s.Api.hw.(*hardware.TestHardware).ResetRuntimes()
}
