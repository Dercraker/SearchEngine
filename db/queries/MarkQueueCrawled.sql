-- name: MarkQueueCrawled :exec
UPDATE crawl_queue
SET status      = 'crawled',
    last_error  = NULL,
    next_run_at = now(),
    updated_at  = now()
WHERE url = $1;