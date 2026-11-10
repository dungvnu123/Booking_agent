package integration

import (
	"abclinic.local/backend/internal/platform/postgres"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"os"
	"testing"
	"time"
)

func TestPostgresPoolAndTransaction(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to a disposable PostgreSQL database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := postgres.Open(ctx, dsn, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatal("PostgreSQL is unreachable")
	}
	marker := errors.New("rollback requested")
	err = postgres.InTx(ctx, pool, func(tx pgx.Tx) error {
		var one int
		if err := tx.QueryRow(ctx, "SELECT 1").Scan(&one); err != nil {
			return err
		}
		if one != 1 {
			return errors.New("unexpected value")
		}
		return marker
	})
	if !errors.Is(err, marker) {
		t.Fatalf("transaction error=%v", err)
	}
	if err := postgres.InTx(ctx, pool, func(tx pgx.Tx) error { _, err := tx.Exec(ctx, "SELECT 1"); return err }); err != nil {
		t.Fatal(err)
	}
}
