-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN password_hash VARCHAR(255),
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN email_verification_token VARCHAR(255),
ADD COLUMN email_verification_expires_at TIMESTAMPTZ,
ADD COLUMN password_reset_token VARCHAR(255),
ADD COLUMN password_reset_expires_at TIMESTAMPTZ,
ADD COLUMN oauth_provider VARCHAR(50),
ADD COLUMN oauth_provider_id VARCHAR(255),
ADD COLUMN last_login_at TIMESTAMPTZ;

CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;

CREATE INDEX idx_users_password_reset_token ON users(password_reset_token) WHERE password_reset_token IS NOT NULL;

CREATE INDEX idx_users_oauth ON users(oauth_provider, oauth_provider_id) WHERE oauth_provider IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_oauth;
DROP INDEX IF EXISTS idx_users_password_reset_token;
DROP INDEX IF EXISTS idx_users_email_verification_token;

ALTER TABLE users
DROP COLUMN IF EXISTS last_login_at,
DROP COLUMN IF EXISTS avatar_url,
DROP COLUMN IF EXISTS oauth_provider_id,
DROP COLUMN IF EXISTS oauth_provider,
DROP COLUMN IF EXISTS password_reset_expires_at,
DROP COLUMN IF EXISTS password_reset_token,
DROP COLUMN IF EXISTS email_verification_expires_at,
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS password_hash;
-- +goose StatementEnd
