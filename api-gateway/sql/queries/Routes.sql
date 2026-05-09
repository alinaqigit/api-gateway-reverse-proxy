-- name: CreateRoute :one
INSERT INTO routes (
    route_name,
    path_prefix,
    upstream_url,
    auth_required,
    rate_limit,
    timeout_ms,
    strip_prefix,
    version,
    status
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;