-- name: TouchDocumentFetchedAt :exec
UPDATE documents
SET fetched_at = now(),
    status_code = $2,
    content_type = $3
WHERE url = $1;