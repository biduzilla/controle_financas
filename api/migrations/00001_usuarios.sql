-- +goose Up
-- +goose StatementBegin

-- Habilita extensão para geração de UUIDs nativa (Postgres 13+)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS usuarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    nome VARCHAR(100) NOT NULL 
        CHECK (char_length(nome) >= 3 AND char_length(nome) <= 100),
    telefone VARCHAR(20) NOT NULL,
    email VARCHAR(255) NOT NULL,
    is_ativo BOOLEAN NOT NULL DEFAULT true,
    password_hash BYTEA NOT NULL,

    -- Auditoria (BaseModel)
    version INTEGER NOT NULL DEFAULT 1,
    deleted BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,

    -- Constraints com nomes explícitos
    CONSTRAINT user_email_key UNIQUE (email),
    CONSTRAINT user_telefone_key UNIQUE (telefone),
    CONSTRAINT email_format CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    )   -- ← vírgula removida
);

-- Índices parciais para soft delete
CREATE INDEX IF NOT EXISTS idx_usuarios_email_active 
    ON usuarios(email) WHERE deleted = false;

CREATE INDEX IF NOT EXISTS idx_usuarios_telefone_active 
    ON usuarios(telefone) WHERE deleted = false AND is_ativo = true;

CREATE INDEX IF NOT EXISTS idx_usuarios_is_ativo 
    ON usuarios(is_ativo) WHERE deleted = false;

CREATE INDEX IF NOT EXISTS idx_usuarios_created_at 
    ON usuarios(created_at) WHERE deleted = false;

COMMENT ON TABLE usuarios IS 'Tabela de usuários do sistema de controle financeiro';
COMMENT ON COLUMN usuarios.nome IS 'Nome completo (3-100 caracteres)';
COMMENT ON COLUMN usuarios.telefone IS 'Telefone formato brasileiro: (DD) 9XXXX-XXXX';
COMMENT ON COLUMN usuarios.email IS 'Email único para autenticação';
COMMENT ON COLUMN usuarios.is_ativo IS 'Flag de ativação lógica do usuário';
COMMENT ON COLUMN usuarios.password_hash IS 'Hash bcrypt da senha (custo 12), armazenado como BYTEA';
COMMENT ON COLUMN usuarios.version IS 'Versão para controle de concorrência otimista';
COMMENT ON COLUMN usuarios.deleted IS 'Flag de soft delete';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS usuarios CASCADE;
-- +goose StatementEnd