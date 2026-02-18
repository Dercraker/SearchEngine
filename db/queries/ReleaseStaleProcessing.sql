-- name: ReleaseStaleProcessing :exec
UPDATE crawl_queue
SET status     = 'pending',
    locked_at  = NULL,
    updated_at = now()
WHERE status = 'processing'
  AND locked_at IS NOT NULL
  AND locked_at < (now() - ($1)::interval);
