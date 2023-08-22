package db

import (
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	*test.MultiDBSuite
}

func TestDBPackage(t *testing.T) {
	testdataDbs := test.FindTestdataDBs()
	dbSuites := map[string]*test.DBSuite{
		CocktailsDB: test.NewDBSuite("cocktails", "test", "", testdataDbs),
		OpenBarDB:   test.NewDBSuite("openbar", "test", "../../schema/openbardb/", ""),
	}

	suite.Run(t, &testSuite{MultiDBSuite: test.NewMultiDBSuite(dbSuites)})
}
