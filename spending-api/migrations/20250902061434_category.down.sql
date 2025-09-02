DROP TABLE categories;

ALTER TABLE spending_records
DROP COLUMN category_id;

ALTER TABLE spending_records
ADD COLUMN category TEXT;