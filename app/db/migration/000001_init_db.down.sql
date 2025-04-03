-- DROP TABLE IF EXISTS "verify_emails";

-- DROP TABLE IF EXISTS "session";

-- DROP TABLE IF EXISTS "users";

-- DROP TABLE IF EXISTS "pages";

-- Drop functions
DROP FUNCTION IF EXISTS public.get_low_stock_medicines();
DROP FUNCTION IF EXISTS public.get_expiring_medicines(integer);

-- Drop indexes
DROP INDEX IF EXISTS idx_medicine_transactions_supplier_id;
DROP INDEX IF EXISTS idx_medicine_transactions_transaction_type;
DROP INDEX IF EXISTS idx_medicine_transactions_transaction_date;
DROP INDEX IF EXISTS idx_medicine_transactions_medicine_id;
DROP INDEX IF EXISTS idx_medicine_suppliers_name;
DROP INDEX IF EXISTS idx_medicines_reorder_level;

-- Drop tables with foreign keys first
DROP TABLE IF EXISTS public.medicine_transactions;
DROP TABLE IF EXISTS public.medicine_suppliers;

-- Alter the existing medicines table to revert changes
ALTER TABLE public.medicines 
  DROP COLUMN IF EXISTS unit_price,
  DROP COLUMN IF EXISTS reorder_level,
  DROP COLUMN IF EXISTS created_at,
  DROP COLUMN IF EXISTS updated_at;




