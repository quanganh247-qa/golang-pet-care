
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

-- -- name: AddTestResult :one
-- INSERT INTO test_results (
--     ordered_test_id, 
--     parameter_name, 
--     result_value, 
--     normal_range, 
--     units, 
--     interpretation
-- ) VALUES (
--     $1, $2, $3, $4, $5, $6
-- ) RETURNING *;

-- name: GetOrderedTestsByAppointment :many
SELECT 
    ot.id AS ordered_test_id,
    t.test_id,
    t.name AS test_name,
    t.category_id,
    tc.name AS category_name,
    ot.price_at_order,
    ot.status,
    ot.created_at AS ordered_date,
    o.notes
FROM test_orders o
JOIN ordered_tests ot ON o.order_id = ot.order_id
JOIN tests t ON ot.test_id = t.id
JOIN test_categories tc ON t.category_id = tc.category_id
WHERE o.appointment_id = $1
ORDER BY ot.created_at DESC;

-- name: GetAllAppointmentsWithOrders :many
SELECT 
    a.appointment_id,
    a.date,
    p.name AS pet_name,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'order_id', o.order_id,
                'total_amount', o.total_amount,
                'status', o.status,
                'order_date', o.order_date,
                'tests', (
                    SELECT jsonb_agg(jsonb_build_object(
                        'test_id', t.test_id,
                        'test_name', t.name,
                        'price', ot.price_at_order,
                        'status', ot.status
                    ))
                    FROM ordered_tests ot
                    JOIN tests t ON ot.test_id = t.id
                    WHERE ot.order_id = o.order_id
                )
            ) 
        ) FILTER (WHERE o.order_id IS NOT NULL),
        '[]'::jsonb
    ) AS orders
FROM appointments a
LEFT JOIN test_orders o ON a.appointment_id = o.appointment_id
JOIN pets p ON a.petid = p.petid 
GROUP BY a.appointment_id, p.name
ORDER BY a.date DESC;

-- name: GetTestOrderByID :one
SELECT * FROM test_orders WHERE order_id = $1;


-- name: GetOrderedTestsByOrderID :many
SELECT  ordered_tests.*, 
        tests.name AS test_name,
        tests.description AS test_description,
        tests.price AS test_price
FROM ordered_tests 
LEFT JOIN tests ON ordered_tests.test_id = tests.id
LEFT JOIN test_categories ON tests.category_id = test_categories.category_id
WHERE ordered_tests.order_id = $1;

-- name: UpdateTestOrderStatus :exec
UPDATE test_orders
SET status = $2
WHERE order_id = $1;
