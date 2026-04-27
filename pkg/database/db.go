package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Solo una declaración de QueryLogger en todo el paquete
type QueryLogger struct{}

func (ql *QueryLogger) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Printf("--- SQL DEBUG ---")
	log.Printf("Query: %s", data.SQL)
	log.Printf("Args:  %v", data.Args)
	log.Printf("-----------------")
	return ctx
}

func (ql *QueryLogger) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://myuser:mypassword@localhost:5432/reminders_db?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("error parseando config: %w", err)
	}

	config.ConnConfig.Tracer = &QueryLogger{}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("error conectando a DB: %w", err)
	}

	DB = pool
	log.Println("✅ DB Conectada y SQL Tracer listo")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
