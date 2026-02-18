-- name: EnsureQueueURL :exec
INSERT INTO crawl_queue (url, status, attempts, last_error, next_run_at)
VALUES ($1, 'pending', 0, NULL, now())
    ON CONFLICT (url) DO NOTHING;