package cafe

import "time"

type Info struct {
	ID          string    `db:"cafe_id" json:"id"`
	Name        string    `db:"name" json:"name"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewCafe struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCafe struct {
	Name *string `json:"name"`
}
