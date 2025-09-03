ALTER Table categories DROP COLUMN is_deleted;
ALTER Table categories DROP COLUMN deleted_at;
ALTER Table spending_records DROP COLUMN is_deleted;
ALTER Table spending_records DROP COLUMN deleted_at;

DROP INDEX IF EXISTS idx_categories_name;