-- +goose Up

CREATE TABLE route_roles (
    id SERIAL PRIMARY KEY,

    route_id UUID REFERENCES routes(id) ON DELETE CASCADE,

    role_name VARCHAR(50) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS route_roles;