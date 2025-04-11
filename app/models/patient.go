package models

// Patient represents a pet patient in the system
type Patient struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Species   string  `json:"species"`
	Breed     string  `json:"breed"`
	Age       float64 `json:"age"`
	Gender    string  `json:"gender"`
	Weight    float64 `json:"weight"`
	Status    string  `json:"status"`
	Notes     string  `json:"notes"`
	OwnerID   int64   `json:"owner_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
