-- name: MarkQueueFailed :exec
UPDATE crawl_queue
SET status      = 'failed',
    attempts    = attempts + 1,
    last_error  = $2,
    next_run_at = now() + ($3::interval),
    updated_at  = now()
WHERE url = $1;