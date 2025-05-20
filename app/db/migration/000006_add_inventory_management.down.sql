-- -- +migrate Down

-- -- Drop the trigger first
-- DROP TRIGGER IF EXISTS after_inventory_transaction ON inventory_transactions;

-- -- Drop the function
-- DROP FUNCTION IF EXISTS update_inventory_quantity();

-- -- Drop the view
-- DROP VIEW IF EXISTS inventory_status;

-- -- Drop the indexes
-- DROP INDEX IF EXISTS idx_inventory_transactions_reference;
-- DROP INDEX IF EXISTS idx_inventory_transactions_date;
-- DROP INDEX IF EXISTS idx_inventory_transactions_type;
-- DROP INDEX IF EXISTS idx_inventory_transactions_item_id;

-- DROP INDEX IF EXISTS idx_inventory_items_supplier;
-- DROP INDEX IF EXISTS idx_inventory_items_quantity;
-- DROP INDEX IF EXISTS idx_inventory_items_expiration;
-- DROP INDEX IF EXISTS idx_inventory_items_type;
-- DROP INDEX IF EXISTS idx_inventory_items_name;

-- -- Drop the tables
-- DROP TABLE IF EXISTS inventory_transactions;
-- DROP TABLE IF EXISTS inventory_items;

-- -- Drop the enum
-- DROP TYPE IF EXISTS inventory_item_type_enum; 