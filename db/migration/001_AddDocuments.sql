-- +goose Up
-- +goose StatementBegin
-- Extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Table: documents
CREATE TABLE documents (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    url           text NOT NULL UNIQUE,

    fetched_at    timestamptz NULL,
    status_code   integer NULL,
    content_type  text NULL,

    content_hash  text NULL,
    raw_content   text NULL,

    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX idx_documents_fetched_at
    ON documents (fetched_at DESC);

CREATE INDEX idx_documents_content_hash
    ON documents (content_hash);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_documents_fetched_at;
DROP INDEX idx_documents_content_hash;
DROP TABLE documents;
-- +goose StatementEnd
