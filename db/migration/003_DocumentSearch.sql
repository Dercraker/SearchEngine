-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS document_search (
  document_id    uuid PRIMARY KEY REFERENCES documents(id) ON DELETE CASCADE,
  search_vector  tsvector NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_document_search_vector
  ON document_search
  USING GIN (search_vector);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_document_search_vector;
DROP TABLE document_search;
-- +goose StatementEnd
