-- +goose Up
-- +goose StatementBegin

UPDATE documents
SET status_code = 0
WHERE status_code IS NULL;
UPDATE documents
SET content_type = ''
WHERE content_type IS NULL;
UPDATE documents
SET fetched_at = now()
WHERE fetched_at IS NULL;

ALTER TABLE documents
    ALTER COLUMN status_code SET NOT NULL,
ALTER
COLUMN content_type SET NOT NULL,
  ALTER
COLUMN fetched_at SET NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE documents
    ALTER COLUMN status_code DROP NOT NULL,
ALTER
COLUMN content_type DROP
NOT NULL,
  ALTER
COLUMN fetched_at DROP
NOT NULL;

-- +goose StatementEnd
