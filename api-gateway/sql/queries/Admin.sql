-- name: CreateAdmin :one
INSERT INTO admins (name, email, password_hash, is_superadmin)
VALUES ($1, $2, $3, FALSE)
RETURNING id, name, email, password_hash, is_superadmin;

-- name: GetAdminByID :one
SELECT id, name, email, is_superadmin FROM admins WHERE id = $1;

-- name: GetAdminByEmail :one
SELECT id, name, email, password_hash, is_superadmin FROM admins WHERE email = $1;

-- name: UpdateAdmin :exec
UPDATE admins SET name = $2, email = $3 WHERE id = $1;

-- name: DeleteAdmin :exec
DELETE FROM admins WHERE id = $1;