-- name: CreateService :one
INSERT INTO services (
<<<<<<< HEAD
<<<<<<< HEAD
    name, description, duration, cost, category
) VALUES (
    $1, $2, $3, $4, $5
=======
    name, description, duration, cost, category, notes
) VALUES (
    $1, $2, $3, $4, $5, $6
>>>>>>> b393bb9 (add service and add permission)
=======
    name, description, duration, cost, category
) VALUES (
    $1, $2, $3, $4, $5
>>>>>>> ada3717 (Docker file)
)
RETURNING *;

-- name: GetServices :many
SELECT * FROM services where removed_at is NULL ORDER BY name LIMIT $1 OFFSET $2;

-- name: GetServiceByID :one
<<<<<<< HEAD
<<<<<<< HEAD
SELECT * FROM services
WHERE id = $1;
=======
SELECT * FROM Service 
WHERE serviceID = $1 LIMIT 1;

-- name: GetAllServices :many
<<<<<<< HEAD
SELECT * FROM Service ORDER BY serviceID LIMIT $1 OFFSET $2;
>>>>>>> cfbe865 (updated service response)
=======
SELECT * FROM Service ORDER BY name LIMIT $1 OFFSET $2;
>>>>>>> 5e493e4 (get all services)
=======
SELECT * FROM services
<<<<<<< HEAD
WHERE id = $1 and removed_at is NULL;
>>>>>>> b393bb9 (add service and add permission)
=======
WHERE id = $1;
>>>>>>> ffc9071 (AI suggestion)

-- name: DeleteService :exec
UPDATE services
SET removed_at = NOW()
<<<<<<< HEAD
<<<<<<< HEAD
WHERE id = $1;
=======
WHERE id = $1 and removed_at is NULL;
>>>>>>> b393bb9 (add service and add permission)
=======
WHERE id = $1;
>>>>>>> ffc9071 (AI suggestion)

-- name: UpdateService :one
UPDATE services
SET 
    name = $2,
    description = $3,
    duration = $4,
    cost = $5,
    category = $6,
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ada3717 (Docker file)
    updated_at = NOW()
WHERE id = $1
=======
    notes = $7,
    updated_at = NOW()
<<<<<<< HEAD
WHERE id = $1 and removed_at is NULL
>>>>>>> b393bb9 (add service and add permission)
RETURNING *;
<<<<<<< HEAD
=======
WHERE id = $1
RETURNING *;
>>>>>>> ffc9071 (AI suggestion)
=======

>>>>>>> ada3717 (Docker file)
