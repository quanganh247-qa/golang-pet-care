-- -- name: CreateServicePackage :one
-- INSERT INTO service_packages (
--     name,
--     description,
--     price,
--     duration_days,
--     is_active
-- ) VALUES (
--     $1, $2, $3, $4, $5
-- ) RETURNING *;

-- -- name: GetServicePackages :many
-- SELECT * FROM service_packages 
-- WHERE is_active = true
-- ORDER BY created_at DESC
-- LIMIT $1 OFFSET $2;

-- -- name: GetServicePackageByID :one
-- SELECT * FROM service_packages
-- WHERE id = $1 AND is_active = true;

-- -- name: UpdateServicePackage :one
-- UPDATE service_packages
-- SET 
--     name = COALESCE($2, name),
--     description = COALESCE($3, description),
--     price = COALESCE($4, price),
--     duration_days = COALESCE($5, duration_days),
--     updated_at = CURRENT_TIMESTAMP
-- WHERE id = $1
-- RETURNING *;

-- -- name: DeleteServicePackage :exec
-- UPDATE service_packages
-- SET 
--     is_active = false,
--     updated_at = CURRENT_TIMESTAMP
-- WHERE id = $1;

-- -- name: AddServiceToPackage :one
-- INSERT INTO package_services (
--     package_id,
--     service_id,
--     quantity,
--     discount_percentage
-- ) VALUES (
--     $1, $2, $3, $4
-- ) RETURNING *;

-- -- name: GetPackageServices :many
-- SELECT 
--     ps.*,
--     s.name as service_name,
--     s.description as service_description,
--     s.cost as service_cost
-- FROM package_services ps
-- JOIN services s ON ps.service_id = s.id
-- WHERE ps.package_id = $1;

-- -- name: RemoveServiceFromPackage :exec
-- DELETE FROM package_services
-- WHERE package_id = $1 AND service_id = $2;

-- -- name: CreatePackageDiscount :one
-- INSERT INTO package_discounts (
--     package_id,
--     discount_type,
--     discount_value,
--     start_date,
--     end_date,
--     min_purchase_amount,
--     max_discount_amount,
--     usage_limit
-- ) VALUES (
--     $1, $2, $3, $4, $5, $6, $7, $8
-- ) RETURNING *;

-- -- name: GetActivePackageDiscounts :many
-- SELECT * FROM package_discounts
-- WHERE package_id = $1 
-- AND is_active = true
-- AND (CURRENT_DATE BETWEEN start_date AND end_date OR (start_date IS NULL AND end_date IS NULL))
-- AND (usage_limit IS NULL OR current_usage < usage_limit);

-- -- name: UpdatePackageDiscount :one
-- UPDATE package_discounts
-- SET 
--     discount_type = COALESCE($2, discount_type),
--     discount_value = COALESCE($3, discount_value),
--     start_date = COALESCE($4, start_date),
--     end_date = COALESCE($5, end_date),
--     min_purchase_amount = COALESCE($6, min_purchase_amount),
--     max_discount_amount = COALESCE($7, max_discount_amount),
--     usage_limit = COALESCE($8, usage_limit),
--     is_active = COALESCE($9, is_active)
-- WHERE id = $1
-- RETURNING *;

-- -- name: IncrementDiscountUsage :exec
-- UPDATE package_discounts
-- SET current_usage = current_usage + 1
-- WHERE id = $1;