-- name: UpsertDocumentFetch :one
INSERT INTO documents (url, fetched_at, status_code, content_type, content_hash, raw_content)
VALUES ($1, now(), $2, $3, $4, $5)
ON CONFLICT (url) DO UPDATE
    SET fetched_at = now(),
    status_code = EXCLUDED.status_code,
    content_type = EXCLUDED.content_type,
    content_hash = EXCLUDED.content_hash,
    raw_content = EXCLUDED.raw_content
RETURNING id, url, fetched_at, status_code, content_type, content_hash;