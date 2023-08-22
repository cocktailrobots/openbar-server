package dbutils

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gocraft/dbr/v2"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TxProvider interface {
	SetBranch(branch string)
	GetBranch() string
	Transaction(ctx context.Context, txFunc func(tx *dbr.Tx) error) error
}

var _ TxProvider = &DBProvider{}

type DBProvider struct {
	*DBBranchTracker
	conn *dbr.Connection
	sess *dbr.Session
}

func NewDBProvider(conn *dbr.Connection, database, migrationSchemaDir string) (*DBProvider, error) {
	if conn == nil {
		panic("nil connection")
	}

	if len(migrationSchemaDir) > 0 {
		driver, err := mysql.WithInstance(conn.DB, &mysql.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create migration driver: %w", err)
		}

		absPath, _ := filepath.Abs(migrationSchemaDir)
		migrate, err := migrate.NewWithDatabaseInstance("file://"+absPath, database, driver)
		if err != nil {
			return nil, fmt.Errorf("failed to create migration instance: %w", err)
		}

		err = migrate.Up()
		if err != nil {
			return nil, err
		}
	}

	return &DBProvider{
		DBBranchTracker: NewDBBranchTracker(),
		conn:            conn,
	}, nil
}

func (dbp *DBProvider) Session() *dbr.Session {
	if dbp.sess == nil {
		dbp.sess = dbp.conn.NewSession(nil)
	}

	return dbp.sess
}

// Transaction executes the txFunc within a transaction. If not committed in the txFunc, the transaction will be rolled back.
func (dbp *DBProvider) Transaction(ctx context.Context, txFunc func(tx *dbr.Tx) error) error {
	if dbp.conn == nil {
		return fmt.Errorf("db connection not initialized")
	}

	branch := dbp.GetBranch()
	sess := dbp.Session()
	if sess == nil {
		return fmt.Errorf("nil session")
	}

	tx, err := sess.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.ExecContext(ctx, fmt.Sprintf("call dolt_checkout('%s')", branch))
	if err != nil {
		return fmt.Errorf("failed to use branch %s: %w", branch, err)
	}

	return txFunc(tx)
}
