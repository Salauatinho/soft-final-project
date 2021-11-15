package staff

import "time"

type Info struct {
	ID          string    `db:"staff_id" json:"staff_id"`
	Position    string    `db:"position" json:"position"`
	FIO         string    `db:"fio" json:"fio"`
	UserID      string    `db:"user_id" json:"user_id"`
	CafeID      string    `db:"cafe_id" json:"cafe_id"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewStaff struct {
	Position string `json:"position" validate:"required"`
	FIO      string `json:"fio" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	CafeID   string `json:"cafe_id" validate:"required"`
}

type UpdateStaff struct {
	Position *string `json:"position"`
	FIO      *string `json:"fio"`
	UserID   *string `json:"user_id"`
	CafeID   *string `json:"cafe_id"`
}
