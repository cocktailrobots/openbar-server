package test

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/stretchr/testify/suite"
)

type MultiDBSuite struct {
	suite.Suite

	dbSuites map[string]*DBSuite
}

func NewMultiDBSuite(dbSuites map[string]*DBSuite) *MultiDBSuite {
	return &MultiDBSuite{
		dbSuites: dbSuites,
	}
}

func (mdbs *MultiDBSuite) SetupSuite() {
	for _, dbs := range mdbs.dbSuites {
		dbs.SetupSuite()
	}
}

func (mdbs *MultiDBSuite) TearDownSuite() {
	for _, dbs := range mdbs.dbSuites {
		dbs.TearDownSuite()
	}
}

func (mdbs *MultiDBSuite) SetupTest() {
	for _, dbs := range mdbs.dbSuites {
		dbs.SetupTest()
	}
}

func (mdbs *MultiDBSuite) TearDownTest() {
	for _, dbs := range mdbs.dbSuites {
		dbs.TearDownTest()
	}
}

func (mdbs *MultiDBSuite) Run(name string, fn func()) {
	for k := range mdbs.dbSuites {
		mdbs.dbSuites[k].BeforeTest(k, name)
	}

	mdbs.Suite.Run(name, fn)

	for k := range mdbs.dbSuites {
		mdbs.dbSuites[k].AfterTest(k, name)
	}
}

func (mdbs *MultiDBSuite) BeginTx(ctx context.Context, db string) (*dbr.Tx, error) {
	suite, ok := mdbs.dbSuites[db]
	if !ok {
		return nil, fmt.Errorf("unknown db: %s", db)
	}

	return suite.BeginTx(ctx)
}
