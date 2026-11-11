package migrations

import (
	"context"
	_ "embed"

	"abclinic.local/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed 000001_identity.up.sql
var IdentitySQL string

func Apply(ctx context.Context, pool *pgxpool.Pool) error {
	return postgres.InTx(ctx, pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended('abclinic:migrations',0))`); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version BIGINT PRIMARY KEY,applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP)`); err != nil {
			return err
		}
		for index, query := range []string{IdentitySQL} {
			version := index + 1
			var applied bool
			if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, version).Scan(&applied); err != nil {
				return err
			}
			if applied {
				continue
			}
			if _, err := tx.Exec(ctx, query); err != nil {
				return err
			}
			if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version); err != nil {
				return err
			}
		}
		return nil
	})
}
