package openbarapi

import (
	"testing"

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
	api := New(logger, dbSuite, mux.NewRouter())
	suite.Run(t, &testSuite{
		DBSuite: dbSuite,
		Api:     api,
	})
}
