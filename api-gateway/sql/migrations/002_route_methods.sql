-- +goose Up
CREATE TABLE route_methods (
    id SERIAL PRIMARY KEY,

    route_id UUID REFERENCES routes(id) ON DELETE CASCADE,

    method VARCHAR(10) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS route_methods;