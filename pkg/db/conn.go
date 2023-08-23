package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
)

const (
	CocktailsDB = "cocktails"
	OpenBarDB   = "openbardb"
)

// ConnParams contains the parameters needed to connect to the database
type ConnParams struct {
	Host            string
	User            string
	Pass            string
	DbName          string
	Branch          string
	Port            int
	MultiStatements bool
}

// String returns a string representation of the connection parameters
func (cp ConnParams) String() string {
	return fmt.Sprintf("host=%s, user=%s, pass=%s, dbName=%s, branch=%s, port=%d", cp.Host, cp.User, cp.Pass, cp.DbName, cp.Branch, cp.Port)
}

// NewConn initializes the database connection
func NewConn(ctx context.Context, params *ConnParams) (*dbr.Connection, error) {
	cfg := mysql.NewConfig()
	if params.Branch != "" && params.DbName != "" {
		cfg.DBName = fmt.Sprintf("%s/%s", params.DbName, params.Branch)
	} else if params.DbName != "" {
		cfg.DBName = params.DbName
	} else if params.Branch != "" {
		return nil, fmt.Errorf("dbName must be set if branch is set")
	}

	cfg.User = params.User
	cfg.Addr = fmt.Sprintf("%s:%d", params.Host, params.Port)

	if len(params.Pass) > 0 {
		cfg.Passwd = params.Pass
	}

	cfg.ParseTime = true
	cfg.MultiStatements = params.MultiStatements

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		log.Printf("failed to create mysql connector: %s", err.Error())
		return nil, err
	}

	db := sql.OpenDB(connector)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping db params: %s, cause: %w", params.String(), err)
	}

	log.Println("successfully connected to db")
	return &dbr.Connection{
		DB:            db,
		Dialect:       dialect.MySQL,
		EventReceiver: &dbr.NullEventReceiver{},
	}, nil
}
