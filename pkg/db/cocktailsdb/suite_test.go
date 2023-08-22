package cocktailsdb

import (
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	*test.DBSuite
}

func TestCocktailsDBPackage(t *testing.T) {
	testdataDbs := test.FindTestdataDBs()
	dbsuite := test.NewDBSuite("cocktails", "test", "", testdataDbs)
	suite.Run(t, &testSuite{DBSuite: dbsuite})
}
