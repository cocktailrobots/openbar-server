package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cocktailrobots/openbar-server/pkg/apis/cocktailsapi"
	"github.com/cocktailrobots/openbar-server/pkg/apis/openbarapi"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/cocktailrobots/openbar-server/pkg/hardware"
	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	"github.com/gocraft/dbr/v2"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	mainBranch       = "main"
	migrationsDirArg = "migration-dir"
)

func installSignalHandler(cancelCtx context.CancelFunc) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		s := <-sigc
		log.Println("Got signal: ", s)
		cancelCtx()
	}()
}

func main() {
	migrationsDir := flag.String(migrationsDirArg, "", "run migrations")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Usage: openbar-server [-migration-dir=<migration_file_dir>] <config file>")
	}

	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to create logger - " + err.Error())
	}
	defer logger.Sync()

	ctx, cancelCtx := context.WithCancel(ctx)

	configFile := flag.Args()[0]
	config, err := ReadConfig(configFile, logger)
	if err != nil {
		log.Fatal("Failed to read " + configFile + " - " + err.Error())
	}

	if migrationsDir != nil && len(*migrationsDir) > 0 {
		err := runMigrations(ctx, *migrationsDir, config)
		if err != nil {
			log.Fatal(fmt.Printf("Failed to run migrations: %s", err.Error()))
		} else {
			log.Printf("Successfully ran migrations")
			os.Exit(0)
		}
	}

	installSignalHandler(cancelCtx)
	err = run(ctx, logger, config)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func runMigrations(ctx context.Context, migrationsDir string, config *Config) error {
	openBarConn, err := connectToDB(ctx, "", "", config, true)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer openBarConn.Close()

	return dbutils.MigrateUp(openBarConn, db.OpenBarDB, migrationsDir)
}

func run(ctx context.Context, logger *zap.Logger, config *Config) error {
	cockConn, err := connectToDB(ctx, db.CocktailsDB, mainBranch, config, false)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer cockConn.Close()

	openBarConn, err := connectToDB(ctx, db.OpenBarDB, mainBranch, config, false)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer openBarConn.Close()

	hw, err := initHardware(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to initialize hardware: %w", err)
	}
	defer func() {
		hardware.TurnPumpsOff(hw)
		hw.Close()
	}()

	// debug hardware changes the iostreams so we need to reinitialize the logger
	logger, err = zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			hw.Update(logger)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	cockDBP, err := dbutils.NewDBProvider(cockConn, db.CocktailsDB+"/"+mainBranch, "")
	if err != nil {
		return fmt.Errorf("failed to initialize database '%s' provider: %w", db.CocktailsDB, err)
	}

	openbarDBP, err := dbutils.NewDBProvider(openBarConn, db.OpenBarDB+"/"+mainBranch, "")
	if err != nil {
		return fmt.Errorf("failed to initialize database '%s' provider: %w", db.OpenBarDB, err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		rtr := mux.NewRouter()
		openbarapi.New(logger, openbarDBP, rtr)
		return startHttpServer(ctx, config.OpenBarApi, rtr)
	})

	eg.Go(func() error {
		rtr := mux.NewRouter()
		cocktailsapi.New(logger, cockDBP, rtr)
		return startHttpServer(ctx, config.CocktailsApi, rtr)
	})

	return eg.Wait()
}

func initHardware(ctx context.Context, config *Config) (hardware.Hardware, error) {
	var hw hardware.Hardware
	var err error

	switch {
	case config.Hardware.Debug != nil:
		dbgConfig := config.Hardware.Debug
		hw, err = hardware.NewDebugHardware(dbgConfig.NumPumps, dbgConfig.OutFile)
		if err != nil {
			return nil, fmt.Errorf("error creating debug hardware: %w", err)
		}
	}

	err = hardware.TurnPumpsOff(hw)
	if err != nil {
		return nil, fmt.Errorf("error turning pumps off: %w", err)
	}

	return hw, nil
}

func connectToDB(ctx context.Context, database, branch string, config *Config, multiStatements bool) (*dbr.Connection, error) {
	if config.DB.Host == nil || *config.DB.Host == "" {
		return nil, fmt.Errorf("no database host specified")
	} else if config.DB.User == nil || *config.DB.User == "" {
		return nil, fmt.Errorf("no database user specified")
	} else if config.DB.Port == nil || *config.DB.Port == 0 {
		return nil, fmt.Errorf("no database port specified")
	}

	params := &db.ConnParams{
		Host:            *config.DB.Host,
		User:            *config.DB.User,
		Port:            *config.DB.Port,
		DbName:          database,
		Branch:          branch,
		MultiStatements: multiStatements,
	}

	if config.DB.Pass != nil && *config.DB.Pass != "" {
		params.Pass = *config.DB.Pass
	}

	return db.NewConn(ctx, params)
}

// startHttpServer starts an HTTP Server on the given port.
func startHttpServer(ctx context.Context, listener *ListenerConfig, mux http.Handler) error {
	addr := fmt.Sprintf("%s:%d", listener.GetHost(), listener.GetPort())
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down HTTP server")
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Printf("Error shutting down HTTP server: %s", err.Error())
		}
	}()

	log.Printf("Starting HTTP server listening on '%s'", addr)
	return srv.ListenAndServe()
}
