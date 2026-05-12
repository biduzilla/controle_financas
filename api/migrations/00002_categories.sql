-- +goose Up
-- +goose StatementBegin
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(500) NOT NULL,
    type INTEGER NOT NULL CHECK (type IN (1, 2)),
    color VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,

    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID
);

CREATE INDEX IF NOT EXISTS categories_nome_key ON categories USING GIN (to_tsvector('simple', name));

CREATE UNIQUE INDEX IF NOT EXISTS categories_name_user_id_unique 
ON categories (user_id, name) WHERE NOT deleted;

CREATE INDEX idx_categories_user_id ON categories(user_id) WHERE NOT deleted;
CREATE INDEX idx_categories_name ON categories(name) WHERE NOT deleted;
CREATE INDEX idx_categories_type ON categories(type) WHERE NOT deleted;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd