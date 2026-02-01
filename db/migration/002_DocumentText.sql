-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS document_text (
  document_id    uuid PRIMARY KEY REFERENCES documents(id) ON DELETE CASCADE,

  title          text NULL,
  body           text NULL,

  language       text NULL,
  cleaned_text   text NULL,

  extracted_at   timestamptz NULL,
  updated_at     timestamptz NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE document_text;
-- +goose StatementEnd
