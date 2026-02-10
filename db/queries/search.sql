-- name: SearchDocuments :many
SELECT dt.title AS title,
       d.url    AS url,
       ts_rank(
               ds.search_vector,
               websearch_to_tsquery('simple', sqlc.arg(q))
       ) ::float8 AS score
FROM public.document_search ds
         JOIN public.documents d
              ON d.id = ds.document_id
         LEFT JOIN public.document_text dt
                   ON dt.document_id = d.id
WHERE ds.search_vector @@ websearch_to_tsquery('simple', sqlc.arg(q))
ORDER BY score DESC, d.fetched_at DESC
    LIMIT sqlc.arg(page_limit)
OFFSET sqlc.arg(page_offset);
