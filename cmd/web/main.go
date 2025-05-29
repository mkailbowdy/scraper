package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" // Using pgx/v5/stdlib because we're using pgx driver through database sql
	"jobscraper.kailmendoza.com/internal/models"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger   *slog.Logger
	shigotos *models.ShigotoModel
}

func main() {
	/* Flags (Configuration) */
	addr := flag.String("addr", "localhost:4000", "HTTP network address")
	dsn := flag.String("dsn", "postgres://user:pass@localhost:5432/letsgosaka", "postgres data source name")
	flag.Parse()

	/* Initialize Dependencies */
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	/* Initialize database pool */
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	/* Initialize App and Inject Dependencies */
	app := &application{
		logger: logger,
		shigotos: &models.ShigotoModel{
			DB: db,
		},
	}

	/* Initialize Router (mux) */
	mux := app.routes()

	/* Start Server */
	app.logger.Info("start server", slog.String("addr", *addr))
	err = http.ListenAndServe(*addr, mux) // Pass in router (mux)
	app.logger.Error(err.Error())         // Display error
	os.Exit(1)                            // Exit the program
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
