package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// func NewDB(ctx context.Context, cfg Database) (*pgxpool.Pool, error) {
func NewDB(ctx context.Context, cfg Database) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s&user=%s&password=%s",
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SslMode,
		cfg.User,
		cfg.Pass,
	)

	driverName, err := otelsql.Register("pgx",
		otelsql.AllowRoot(),
		//otelsql.TraceQueryWithArgs(),
		otelsql.TraceQueryWithoutArgs(),
		//otelsql.TraceRowsClose(),
		//otelsql.TraceRowsAffected(),
		//otelsql.WithDatabaseName("my_database"),        // Optional.
		otelsql.WithSystem(semconv.DBSystemPostgreSQL), // Optional.
	)
	if err != nil {
		_ = fmt.Errorf("otelsql driver: %v", err)
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Printf("open DB driver: %v", err)
	}
	return db, nil
	//config, err := pgx.ParseConfig(driverName)
	//if err != nil {
	//	_ = fmt.Errorf("pgx config: %v", err)
	//}
	//conn, err := pgx.ConnectConfig(context.Background(), config)
	//if err != nil {
	//
	//}

	//return pgxpool.New(ctx, dsn)
}
