-- +goose Up
-- +goose StatementBegin

-- 1) Assure une ligne document_search par document
INSERT INTO public.document_search (document_id, search_vector)
SELECT dt.document_id,
       to_tsvector('simple',
                   coalesce(dt.title, '') || ' ' || coalesce(dt.cleaned_text, '') || ' ' || coalesce(dt.body, ''))
FROM public.document_text dt
         LEFT JOIN public.document_search ds ON ds.document_id = dt.document_id
WHERE ds.document_id IS NULL;

-- 2) Fonction qui calcule le vector (pondère titre > texte)
CREATE
OR REPLACE FUNCTION public.sync_document_search_vector()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
INSERT INTO public.document_search(document_id, search_vector)
VALUES (NEW.document_id,
        setweight(to_tsvector('simple', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(NEW.cleaned_text, coalesce(NEW.body, ''))), 'B')) ON CONFLICT (document_id)
  DO
UPDATE SET search_vector = EXCLUDED.search_vector;

RETURN NEW;
END;
$$;

-- 3) Trigger sur document_text (insert + update)
DROP TRIGGER IF EXISTS trg_sync_document_search_vector ON public.document_text;
CREATE TRIGGER trg_sync_document_search_vector
    AFTER INSERT OR
UPDATE OF title, body, cleaned_text
ON public.document_text
    FOR EACH ROW
    EXECUTE FUNCTION public.sync_document_search_vector();

-- 4) Index GIN (tu l’as déjà, mais idempotent)
CREATE INDEX IF NOT EXISTS idx_document_search_vector
    ON public.document_search USING gin (search_vector);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_sync_document_search_vector ON public.document_text;
DROP FUNCTION IF EXISTS public.sync_document_search_vector();
-- (on garde la table et l’index, sauf si tu veux vraiment les supprimer)
-- +goose StatementEnd
