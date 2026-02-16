-- name: GetDocumentHashByURL :one
SELECT content_hash
FROM documents
WHERE url = $1;