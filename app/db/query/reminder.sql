-- name: CreateReminder :one
INSERT INTO Reminders (
    PetID,
    Title,
    Description,
    DueDate,
    RepeatInterval
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;