-- +goose Up

ALTER TABLE admins
ADD COLUMN is_superadmin BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE admins
DROP COLUMN IF EXISTS is_superadmin;
