package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	_ "github.com/dolthub/driver"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/stretchr/testify/suite"
)

type DBSuite struct {
	suite.Suite
	*dbutils.DBProvider

	SchemaDir string
	Database  string
	Branch    string
	Hash      string
	DbDir     string

	mu   *sync.Mutex
	conn *dbr.Connection
	sess *dbr.Session
}

func NewDBSuite(database, branch, schemaDir, dbDir string) *DBSuite {
	if len(dbDir) == 0 {
		tmpDir, err := os.MkdirTemp("", database+"*")

		dsn := "file://" + tmpDir + "?commitname=Billy%20Batson&commitemail=shazam@gmail.com&database=hostedapidb"
		db, err := sql.Open("dolt", dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		r, err := db.Query("CREATE DATABASE " + database)
		if err != nil {
			panic(fmt.Errorf("failed to create database '%s': %w", database, err))
		}
		r.Close()

		dbDir = tmpDir
	}

	return &DBSuite{
		Suite:     suite.Suite{},
		Database:  database,
		Branch:    branch,
		SchemaDir: schemaDir,
		mu:        &sync.Mutex{},
		DbDir:     dbDir,
	}
}

func (s *DBSuite) SetupSuite() {
	database := s.Database
	if len(s.Branch) > 0 {
		database = database + "/" + s.Branch
	}

	openDB, err := sql.Open("dolt", "file://"+s.DbDir+"?commitname=Test%20Committer&commitemail=test@test.com&database="+s.Database+"&multistatements=true")
	if err != nil {
		panic(err)
	}

	s.conn = &dbr.Connection{DB: openDB, Dialect: dialect.MySQL, EventReceiver: &dbr.NullEventReceiver{}}
	s.DBProvider, err = dbutils.NewDBProvider(s.conn, s.Database, s.SchemaDir)
	if err != nil {
		panic(err)
	}
	sess := s.Session()

	tx, err := sess.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	var bh struct {
		Hash string `db:"hash"`
	}
	_, err = tx.Select("*").From("dolt_branches").Where(dbr.Eq("name", s.Branch)).Load(&bh)
	if err != nil {
		panic(err)
	}

	s.Hash = bh.Hash

	s.LogWorking(tx)
}

func (s *DBSuite) TearDownSuite() {
	sess := s.conn.NewSession(nil)

	ctx := context.Background()
	tx, err := sess.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}

	err = s.resetToHash(ctx, tx)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	err = s.conn.Close()
	if err != nil {
		panic(err)
	}
}

func (s *DBSuite) SetupTest()                           {}
func (s *DBSuite) TearDownTest()                        {}
func (s *DBSuite) AfterTest(suiteName, testName string) {}

func (s *DBSuite) BeforeTest(suiteName, testName string) {
	ctx := context.Background()
	err := s.Transaction(ctx, func(tx *dbr.Tx) error {
		err := s.resetToHash(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to reset to hash '%s': %w", s.Hash, err)
		}

		return tx.Commit()
	})

	if err != nil {
		panic(err)
	}
}

func (s *DBSuite) resetToHash(ctx context.Context, tx *dbr.Tx) error {
	_, err := tx.ExecContext(ctx, "call dolt_add('-A')")
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, fmt.Sprintf("call dolt_reset('--hard','%s')", s.Hash))
	return err
}

func (s *DBSuite) BeginTx(ctx context.Context) (*dbr.Tx, error) {
	return s.Session().BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (s *DBSuite) LogWorking(tx *dbr.Tx) {
	var w struct {
		Hash string `db:"working"`
	}
	_, err := tx.SelectBySql("SELECT @@cocktails_working as working;").Load(&w)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("working:", w.Hash)
	}
}

func (s *DBSuite) Run(name string, f func()) {
	s.BeforeTest("suite", name)
	defer s.AfterTest("suite", name)

	s.Suite.Run(name, func() {
		f()
	})
}
