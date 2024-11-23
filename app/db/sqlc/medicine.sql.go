// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: medicine.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const listMedicinesByPet = `-- name: ListMedicinesByPet :many
SELECT 
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
    pm.dosage,
    pm.frequency,
    pm.duration,
    pm.notes AS medicine_notes,
    pt.start_date AS treatment_start_date,
    pt.end_date AS treatment_end_date,
    pt.status AS treatment_status
FROM 
    pet_treatments pt
JOIN 
    treatment_phases tp ON pt.disease_id = tp.disease_id
JOIN 
    phase_medicines pm ON tp.id = pm.phase_id
JOIN 
    medicines m ON pm.medicine_id = m.id
WHERE 
    pt.pet_id = $1 and pt.status = $2 -- Replace with the specific pet_id
ORDER BY 
    tp.phase_number, pm.medicine_id LIMIT $3 OFFSET $4
`

type ListMedicinesByPetParams struct {
	PetID  pgtype.Int8 `json:"pet_id"`
	Status pgtype.Text `json:"status"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

type ListMedicinesByPetRow struct {
	MedicineUsage       pgtype.Text `json:"medicine_usage"`
	MedicineName        string      `json:"medicine_name"`
	MedicineDescription pgtype.Text `json:"medicine_description"`
	Dosage              pgtype.Text `json:"dosage"`
	Frequency           pgtype.Text `json:"frequency"`
	Duration            pgtype.Text `json:"duration"`
	MedicineNotes       pgtype.Text `json:"medicine_notes"`
	TreatmentStartDate  pgtype.Date `json:"treatment_start_date"`
	TreatmentEndDate    pgtype.Date `json:"treatment_end_date"`
	TreatmentStatus     pgtype.Text `json:"treatment_status"`
}

func (q *Queries) ListMedicinesByPet(ctx context.Context, arg ListMedicinesByPetParams) ([]ListMedicinesByPetRow, error) {
	rows, err := q.db.Query(ctx, listMedicinesByPet,
		arg.PetID,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListMedicinesByPetRow{}
	for rows.Next() {
		var i ListMedicinesByPetRow
		if err := rows.Scan(
			&i.MedicineUsage,
			&i.MedicineName,
			&i.MedicineDescription,
			&i.Dosage,
			&i.Frequency,
			&i.Duration,
			&i.MedicineNotes,
			&i.TreatmentStartDate,
			&i.TreatmentEndDate,
			&i.TreatmentStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
