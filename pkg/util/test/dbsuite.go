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
	cp "github.com/otiai10/copy"
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
	// If a dbDir is not provided, create a temporary directory to use for the database.
	// When SetupSuite is called, the migration scripts located in schemaDir will be run.
	// If a dbDir is provided, it is assumed that the database has already been created and
	// is seeded with data.
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
	// copy our database dir to a new directory so that any chunks written during the test, even
	// if they are rolled back, will not be in the git diff.
	if s.DbDir != "" {
		newDbDir := s.T().TempDir()
		err := cp.Copy(s.DbDir, newDbDir)
		if err != nil {
			panic(err)
		}

		s.DbDir = newDbDir
	}

	// Append the branch to the database name if a branch is provided.
	database := s.Database
	if len(s.Branch) > 0 {
		database = database + "/" + s.Branch
	}

	// Use the dolt driver to open a connection to the database.
	openDB, err := sql.Open("dolt", "file://"+s.DbDir+"?commitname=Test%20Committer&commitemail=test@test.com&database="+s.Database+"&multistatements=true")
	if err != nil {
		panic(err)
	}

	s.conn = &dbr.Connection{DB: openDB, Dialect: dialect.MySQL, EventReceiver: &dbr.NullEventReceiver{}}

	// Create a DBProvider which will provide the database to our tests
	s.DBProvider, err = dbutils.NewDBProvider(s.conn, s.Database, s.SchemaDir)
	if err != nil {
		panic(err)
	}
	sess := s.Session()

	// Query the database to get the hash of the last commit on our branch and store it
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

// BeforeTest is called before each test. It calls resetToHash to reset the database to the
// hash of the last commit on the branch.
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
	// Add all tables to the staging area so that everything will be reset and new tables will be deleted
	_, err := tx.ExecContext(ctx, "call dolt_add('-A')")
	if err != nil {
		return err
	}

	// Reset the database to the hash of the last commit on the branch that we stored at set up
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

func (s *DBSuite) SetupSubTest() {
	s.BeforeTest("", "")
}

func (s *DBSuite) TearDownSubTest() {
	s.AfterTest("", "")
}
