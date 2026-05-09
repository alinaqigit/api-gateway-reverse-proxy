-- +goose Up

CREATE TABLE routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    route_name VARCHAR(100) UNIQUE NOT NULL,

    path_prefix TEXT NOT NULL,

    upstream_url TEXT NOT NULL,

    auth_required BOOLEAN DEFAULT FALSE,

    rate_limit INTEGER DEFAULT 0,

    timeout_ms INTEGER DEFAULT 5000,

    strip_prefix BOOLEAN DEFAULT FALSE,

    version VARCHAR(20) DEFAULT 'v1',

    status VARCHAR(20) DEFAULT 'active',

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS routes;