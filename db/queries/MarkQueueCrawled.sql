-- name: MarkQueueCrawled :exec
UPDATE crawl_queue
SET status      = 'crawled',
    last_error  = NULL,
    locked_at  = NULL,
    updated_at  = now()
WHERE url = $1;