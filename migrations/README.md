# Database migrations

Run `scripts/migrate.ps1` to apply the identity schema for roles, users, patients and invalidated tokens. The runner records versions and applies each version in a transaction. API startup does not run migrations.
