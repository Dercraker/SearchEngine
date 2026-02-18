-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS crawl_queue (
  id          bigserial PRIMARY KEY,
  url         text NOT NULL UNIQUE,
  status      text NOT NULL DEFAULT 'pending', -- pending | processing | crawled | failed
  attempts    int NOT NULL DEFAULT 0,
  last_error  text NULL,
  next_run_at timestamptz NOT NULL DEFAULT now(),
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_crawl_queue_next_run
ON crawl_queue (status, next_run_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_crawl_queue_next_run;
DROP TABLE crawl_queue;
-- +goose StatementEnd
