package cocktailsapi

import (
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type testSuite struct {
	*test.DBSuite
	api *CocktailsAPI
}

func TestCocktailsapiPackage(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	testdataDbs := test.FindTestdataDBs()
	dbSuite := test.NewDBSuite("cocktails", "test", "", testdataDbs)

	api := New(logger, dbSuite, mux.NewRouter())
	suite.Run(t, &testSuite{
		DBSuite: dbSuite,
		api:     api,
	})
}
