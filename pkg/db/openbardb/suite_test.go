package openbardb

import (
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testSuite struct {
	*test.DBSuite
}

func TestOpenbarDBPackage(t *testing.T) {
	dbsuite := test.NewDBSuite("openbar", "test", "../../../schema/openbardb/", "")
	suite.Run(t, &testSuite{DBSuite: dbsuite})
}
