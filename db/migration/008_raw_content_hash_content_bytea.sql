-- +goose Up
-- +goose StatementBegin

ALTER TABLE documents
ALTER
COLUMN raw_content TYPE bytea
  USING CASE
    WHEN raw_content IS NULL THEN NULL
    ELSE convert_to(raw_content, 'UTF8')
END;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE documents
ALTER
COLUMN raw_content TYPE text
  USING CASE
    WHEN raw_content IS NULL THEN NULL
    ELSE convert_from(raw_content, 'UTF8')
END;
-- +goose StatementEnd
