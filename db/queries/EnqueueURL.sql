-- name: EnqueueURL :exec
INSERT INTO crawl_queue (url, status, attempts, next_run_at, last_error) 
VALUES ($1, 'pending', 0, now(), null)
    ON CONFLICT (url) DO UPDATE
        SET status = CASE
                        WHEN crawl_queue.status IN ('crawled') THEN crawl_queue.status
                        ELSE 'pending'
                     END,
            next_run_at = LEAST(crawl_queue.next_run_at, now()),
            updated_at = now();