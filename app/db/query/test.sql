-- -- name: CreateTest :one
-- INSERT INTO tests (
--     pet_id,
--     doctor_id,
--     test_type,
--     status,
--     created_at
-- ) VALUES (
--     $1, $2, $3, $4, CURRENT_TIMESTAMP
-- ) RETURNING *;

-- -- name: UpdateTestStatus :exec
-- UPDATE tests 
-- SET status = $2,
--     updated_at = CURRENT_TIMESTAMP
-- WHERE id = $1;

-- -- name: GetTestByID :one
-- SELECT * FROM tests WHERE id = $1;

-- -- name: GetTestsByPetID :many
-- SELECT * FROM tests WHERE pet_id = $1;

-- -- name: GetTestsByDoctorID :many
-- SELECT * FROM tests WHERE doctor_id = $1;

-- -- name: AddTestResult :one
-- INSERT INTO test_results (
--     test_id,
--     parameters,
--     notes,
--     files,
--     created_at
-- ) VALUES (
--     $1, $2, $3, $4, CURRENT_TIMESTAMP
-- ) RETURNING *;

-- -- name: GetTestResults :many
-- SELECT * FROM test_results WHERE test_id = $1;

-- -- name: GetStatusHistory :many
-- SELECT 
--     status,
--     updated_at as timestamp,
--     updated_by
-- FROM test_statuses
-- WHERE test_id = $1
-- ORDER BY updated_at DESC;

-- name: AddTestCategory :exec
INSERT INTO test_categories (category_id, name, description, icon_name)
VALUES ($1, $2, $3, $4);

-- name: CreateTest :one
INSERT INTO tests (test_id, category_id, name, description, price, turnaround_time)
VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetTestCategoryByID :one
SELECT * FROM test_categories WHERE category_id = $1;

-- name: GetTestByID :one
SELECT * FROM tests WHERE id = $1 AND is_active = true;

-- name: UpdateTest :one
UPDATE tests
SET name = $2, description = $3, price = $4, turnaround_time = $5
WHERE test_id = $1
RETURNING *;

-- name: SoftDeleteTest :exec
UPDATE tests
SET is_active = false
WHERE test_id = $1;

-- name: GetTestCategories :many
SELECT * FROM test_categories;

-- name: ListTests :many
SELECT * FROM tests WHERE is_active is true;


-- name: CreateTestOrder :one
INSERT INTO test_orders (
    appointment_id, 
    total_amount,
    notes
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: AddOrderedTest :one
INSERT INTO ordered_tests (
    order_id,
    test_id,
    price_at_order
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTestsByCategory :many
SELECT * FROM tests 
WHERE category_id = $1 AND is_active = TRUE;


-- name: CancelTestOrder :exec
UPDATE test_orders 
SET status = 'cancelled'
WHERE order_id = $1 AND status IN ('pending', 'processing');

-- name: UpdateTestStatus :one
UPDATE ordered_tests
SET status = $2
WHERE id = $1
RETURNING *;

-- name: AddTestResult :one
INSERT INTO test_results (
    ordered_test_id, 
    parameter_name, 
    result_value, 
    normal_range, 
    units, 
    interpretation
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


