-- +goose Up
-- +goose StatementBegin
DROP INDEX uq_documents_content_hash;

CREATE INDEX IF NOT EXISTS idx_documents_content_hash
    ON documents (content_hash);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS uq_documents_content_hash
    ON documents (content_hash)
    WHERE content_hash IS NOT NULL;

DROP INDEX idx_documents_content_hash;
-- +goose StatementEnd
