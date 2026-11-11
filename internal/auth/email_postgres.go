package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"abclinic.local/backend/internal/platform/postgres"
	"github.com/jackc/pgx/v5"
)

func insertEmailToken(ctx context.Context, tx pgx.Tx, token EmailToken) error {
	_, err := tx.Exec(ctx, `INSERT INTO email_verifications(user_id,token_hash,expires_at,last_sent_at)
 VALUES($1,$2,$3,$4)`, token.AccountID, token.Hash, token.ExpiresAt, token.SentAt)
	return err
}

func (r *PostgresRepository) ReplaceEmailToken(ctx context.Context, token EmailToken, cutoff time.Time) (bool, error) {
	if r.pool == nil {
		return false, ErrUnavailable
	}
	stored := false
	err := postgres.InTx(ctx, r.pool, func(tx pgx.Tx) error {
		var active, verified bool
		err := tx.QueryRow(ctx, `SELECT COALESCE(is_active,FALSE),email_verified FROM users WHERE id=$1 FOR UPDATE`, token.AccountID).Scan(&active, &verified)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if !active || verified {
			return nil
		}
		result, err := tx.Exec(ctx, `INSERT INTO email_verifications(user_id,token_hash,expires_at,last_sent_at)
   VALUES($1,$2,$3,$4) ON CONFLICT(user_id) DO UPDATE
   SET token_hash=EXCLUDED.token_hash,expires_at=EXCLUDED.expires_at,last_sent_at=EXCLUDED.last_sent_at
   WHERE email_verifications.last_sent_at<=$5`, token.AccountID, token.Hash, token.ExpiresAt, token.SentAt, cutoff)
		if err != nil {
			return err
		}
		stored = result.RowsAffected() == 1
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("replace email verification: %w", err)
	}
	return stored, nil
}

func (r *PostgresRepository) VerifyEmail(ctx context.Context, hash string, now time.Time) error {
	if r.pool == nil {
		return ErrUnavailable
	}
	return postgres.InTx(ctx, r.pool, func(tx pgx.Tx) error {
		var accountID string
		err := tx.QueryRow(ctx, `SELECT user_id FROM email_verifications WHERE token_hash=$1`, hash).Scan(&accountID)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidEmailToken
		}
		if err != nil {
			return err
		}
		var active, verified bool
		if err := tx.QueryRow(ctx, `SELECT COALESCE(is_active,FALSE),email_verified FROM users WHERE id=$1 FOR UPDATE`, accountID).Scan(&active, &verified); err != nil {
			return err
		}
		if !active {
			return ErrInactive
		}
		if verified {
			return ErrInvalidEmailToken
		}
		var expiresAt time.Time
		err = tx.QueryRow(ctx, `SELECT expires_at FROM email_verifications WHERE user_id=$1 AND token_hash=$2 FOR UPDATE`, accountID, hash).Scan(&expiresAt)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrInvalidEmailToken
		}
		if err != nil {
			return err
		}
		if !expiresAt.After(now) {
			return ErrInvalidEmailToken
		}
		if _, err := tx.Exec(ctx, `UPDATE users SET email_verified=TRUE,email_verified_at=$2::timestamptz,updated_at=CURRENT_TIMESTAMP WHERE id=$1`, accountID, now); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `DELETE FROM email_verifications WHERE user_id=$1`, accountID)
		return err
	})
}
