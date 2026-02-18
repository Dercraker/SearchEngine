-- +goose Up
-- +goose StatementBegin

CREATE
EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE crawl_queue
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE crawl_queue
DROP
CONSTRAINT IF EXISTS crawl_queue_pkey;

ALTER TABLE crawl_queue
ALTER
COLUMN id TYPE uuid
  USING gen_random_uuid();

ALTER TABLE crawl_queue
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE crawl_queue
    ADD CONSTRAINT crawl_queue_pkey PRIMARY KEY (id);

ALTER TABLE crawl_queue
    ADD COLUMN IF NOT EXISTS locked_at TIMESTAMPTZ NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE crawl_queue
DROP
CONSTRAINT IF EXISTS crawl_queue_pkey;

ALTER TABLE crawl_queue
DROP
COLUMN IF EXISTS locked_at;

ALTER TABLE crawl_queue
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE crawl_queue
ALTER
COLUMN id TYPE bigint
  USING NULL;

DO
$$
BEGIN
  IF
NOT EXISTS (
    SELECT 1
    FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE c.relkind = 'S' AND c.relname = 'crawl_queue_id_seq'
  ) THEN
CREATE SEQUENCE crawl_queue_id_seq;
END IF;
END $$;

ALTER TABLE crawl_queue
    ALTER COLUMN id SET DEFAULT nextval('crawl_queue_id_seq'::regclass);

ALTER TABLE crawl_queue
    ADD CONSTRAINT crawl_queue_pkey PRIMARY KEY (id);

-- +goose StatementEnd
