-- +goose Up
-- +goose StatementBegin
    DROP TRIGGER IF EXISTS trg_documents_updated_at on documents;

    CREATE TRIGGER trg_documents_updated_at
        BEFORE UPDATE OF content_hash
        ON documents
        FOR EACH ROW
        WHEN (OLD.content_hash IS DISTINCT FROM NEW.content_hash)
    EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TRIGGER IF EXISTS trg_documents_updated_at ON documents;
    
    -- Revenir au comportement initial: updated_at bouge sur n'importe quel update
    CREATE TRIGGER trg_documents_updated_at
        BEFORE UPDATE
        ON documents
        FOR EACH ROW
        EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd
