package auth

import (
	"context"
	"errors"
	"fmt"

	"abclinic.local/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct{ pool *pgxpool.Pool }

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, value NewAccount) (Account, error) {
	if r.pool == nil {
		return Account{}, ErrUnavailable
	}
	err := postgres.InTx(ctx, r.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`, "auth:email:"+value.Email); err != nil {
			return err
		}
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE lower(email)=$1)`, value.Email).Scan(&exists); err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
		if value.Phone != "" {
			if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`, "auth:phone:"+value.Phone); err != nil {
				return err
			}
			if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)`, value.Phone).Scan(&exists); err != nil {
				return err
			}
			if exists {
				return ErrPhoneExists
			}
		}
		var roleID int64
		if err := tx.QueryRow(ctx, `SELECT role_id FROM roles WHERE role_name='PATIENT'`).Scan(&roleID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrRoleMissing
			}
			return err
		}
		_, err := tx.Exec(ctx, `INSERT INTO users
   (id,role_id,username,full_name,email,password_hash,phone,is_active,created_at,updated_at)
   VALUES ($1,$2,$1,$3,$4,$5,NULLIF($6,''),TRUE,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`,
			value.ID, roleID, value.FullName, value.Email, value.PasswordHash, value.Phone)
		if err != nil {
			return err
		}
		if value.PatientID != "" {
			if err := linkPatient(ctx, tx, value); err != nil {
				return err
			}
		}
		return insertEmailToken(ctx, tx, value.Verification)
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "auth_user_phone_unique" {
				return Account{}, ErrPhoneExists
			}
			if pgErr.ConstraintName == "auth_user_email_unique" {
				return Account{}, ErrEmailExists
			}
			if _, lookupErr := r.FindByEmail(ctx, value.Email); lookupErr == nil {
				return Account{}, ErrEmailExists
			}
		}
		return Account{}, fmt.Errorf("create account: %w", err)
	}
	return value.Account, nil
}

func linkPatient(ctx context.Context, tx pgx.Tx, value NewAccount) error {
	rows, err := tx.Query(ctx, `SELECT id,user_id FROM patients WHERE phone=$1 ORDER BY id LIMIT 2 FOR UPDATE`, value.Phone)
	if err != nil {
		return err
	}
	type patientLink struct {
		ID     string
		UserID *string
	}
	matches, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (patientLink, error) {
		var p patientLink
		err := row.Scan(&p.ID, &p.UserID)
		return p, err
	})
	if err != nil {
		return err
	}
	if len(matches) > 1 {
		return ErrPatientConflict
	}
	if len(matches) == 1 {
		if matches[0].UserID != nil {
			return ErrPhoneExists
		}
		_, err = tx.Exec(ctx, `UPDATE patients SET user_id=$1,updated_at=CURRENT_TIMESTAMP WHERE id=$2`, value.ID, matches[0].ID)
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO patients(id,full_name,phone,email,user_id,created_at,updated_at)
 VALUES ($1,$2,$3,$4,$5,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, value.PatientID, value.FullName, value.Phone, value.Email, value.ID)
	return err
}

const accountQuery = `SELECT u.id,u.email,COALESCE(u.full_name,''),r.role_name,
 COALESCE(u.password_hash,''),COALESCE(u.is_active,FALSE),u.email_verified
 FROM users u JOIN roles r ON r.role_id=u.role_id`

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (Account, error) {
	if r.pool == nil {
		return Account{}, ErrUnavailable
	}
	return scanAccount(r.pool.QueryRow(ctx, accountQuery+` WHERE lower(u.email)=$1`, email))
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (Account, error) {
	if r.pool == nil {
		return Account{}, ErrUnavailable
	}
	return scanAccount(r.pool.QueryRow(ctx, accountQuery+` WHERE u.id=$1`, id))
}

func scanAccount(row pgx.Row) (Account, error) {
	var value Account
	err := row.Scan(&value.ID, &value.Email, &value.FullName, &value.Role, &value.PasswordHash, &value.Active, &value.EmailVerified)
	if errors.Is(err, pgx.ErrNoRows) {
		return Account{}, ErrNotFound
	}
	if err != nil {
		return Account{}, fmt.Errorf("query account: %w", err)
	}
	return value, nil
}

func (r *PostgresRepository) IsInvalidated(ctx context.Context, id string) (bool, error) {
	if r.pool == nil {
		return false, ErrUnavailable
	}
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM invalidated_token WHERE id_token=$1)`, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query token revocation: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepository) Revoke(ctx context.Context, tokens []InvalidatedToken, singleUse bool) error {
	if r.pool == nil {
		return ErrUnavailable
	}
	return postgres.InTx(ctx, r.pool, func(tx pgx.Tx) error {
		for i, token := range tokens {
			result, err := tx.Exec(ctx, `INSERT INTO invalidated_token(id_token,expiry_time) VALUES ($1,$2) ON CONFLICT(id_token) DO NOTHING`, token.JTI, token.ExpiresAt.UTC())
			if err != nil {
				return fmt.Errorf("revoke token: %w", err)
			}
			if singleUse && i == 0 && result.RowsAffected() != 1 {
				return ErrInvalidToken
			}
		}
		return nil
	})
}

var _ Repository = (*PostgresRepository)(nil)
