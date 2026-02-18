-- name: ClaimNextBatch :many
WITH picked AS (SELECT id
                FROM crawl_queue
                WHERE next_run_at <= now()
                  AND status IN ('pending', 'failed')
                ORDER BY next_run_at ASC
    LIMIT $1
    FOR
UPDATE SKIP LOCKED
    )
UPDATE crawl_queue q
SET status     = 'processing',
    locked_at  = now(),
    updated_at = now() FROM picked
WHERE q.id = picked.id
    RETURNING q.id
    , q.url
    , q.status
    , q.attempts
    , q.next_run_at
    , q.last_error;