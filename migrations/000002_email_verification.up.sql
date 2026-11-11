ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS email_verifications (
 user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
 token_hash CHAR(64) NOT NULL UNIQUE,
 expires_at TIMESTAMPTZ NOT NULL,
 last_sent_at TIMESTAMPTZ NOT NULL
);
CREATE INDEX IF NOT EXISTS email_verification_expiry_idx ON email_verifications(expires_at);
