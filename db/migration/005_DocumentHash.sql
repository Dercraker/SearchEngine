-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS uq_documents_content_hash
ON documents (content_hash)
WHERE content_hash IS NOT NULL;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Exemple index trigram (si tu fais ILIKE / similarity sur title/body)
CREATE INDEX IF NOT EXISTS idx_doc_text_title_trgm
ON document_text USING GIN (title gin_trgm_ops);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX uq_documents_content_hash;
DROP INDEX idx_doc_text_title_trgm;
DROP EXTENSION pg_trgm;
-- +goose StatementEnd
