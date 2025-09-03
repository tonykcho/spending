ALTER TABLE categories ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
ALTER TABLE categories ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE spending_records ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
ALTER TABLE spending_records ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE UNIQUE INDEX idx_categories_name ON categories (name) WHERE (is_deleted = FALSE);