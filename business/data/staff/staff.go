package staff

import (
	"context"
	"database/sql"
	"github.com/Salauatinho/soft-final-project/business/auth"
	"github.com/Salauatinho/soft-final-project/business/data/user"
	"github.com/Salauatinho/soft-final-project/foundation/database"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Staff struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Staff {
	return Staff{
		log: log,
		db:  db,
	}
}

func (c Staff) Create(ctx context.Context, traceID string, ns NewStaff, now time.Time) (Info, error) {
	staff := Info{
		ID:          uuid.New().String(),
		Position:    ns.Position,
		FIO:         ns.FIO,
		UserID:      ns.UserID,
		CafeID:      ns.CafeID,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO staffs (staff_id, position, fio, user_id, cafe_id, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4)`

	c.log.Printf("%s : %s query : %s", traceID, "staff.Create",
		database.Log(q, staff.ID, staff.Position, staff.FIO, staff.UserID, staff.CafeID, staff.DateCreated, staff.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, staff.ID, staff.Position, staff.FIO, staff.UserID, staff.CafeID, staff.DateCreated, staff.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting staff")
	}
	return staff, nil
}

func (c Staff) Update(ctx context.Context, traceID string, claims auth.Claims, staffID string, us UpdateStaff, now time.Time) error {
	staff, err := c.QueryByID(ctx, traceID, claims, staffID)
	if err != nil {
		return err
	}

	if us.Position != nil {
		staff.Position = *us.Position
	}
	if us.FIO != nil {
		staff.FIO = *us.FIO
	}
	if us.CafeID != nil {
		staff.CafeID = *us.CafeID
	}
	if us.UserID != nil {
		staff.UserID = *us.UserID
	}
	staff.DateUpdated = now

	const q = `
	UPDATE
		staffs
	SET 
		"position" = $2,
		"fio" = $3,
		"cafe_id" = $4,
		"user_id" = $5,
		date_update =$6
	WHERE
		cafe_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "staff.Update",
		database.Log(q, staff.ID, staff.Position, staff.FIO, staff.UserID, staff.CafeID, staff.DateCreated, staff.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, staff.ID, staff.Position, staff.FIO, staff.UserID, staff.CafeID, staff.DateUpdated); err != nil {
		return errors.Wrap(err, "updating staff")
	}

	return nil
}

func (c Staff) Delete(ctx context.Context, traceID string, staffID string) error {
	if _, err := uuid.Parse(staffID); err != nil {
		return user.ErrInvalidID
	}
	const q = `DELETE FROM staffs where staff_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "staff.Delete",
		database.Log(q, staffID),
	)

	if _, err := c.db.ExecContext(ctx, q, staffID); err != nil {
		return errors.Wrapf(err, "deleting staff %s", staffID)
	}
	return nil
}

func (c Staff) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT * FROM staffs ORDER BY staff_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	c.log.Printf("%s : %s query : %s", traceID, "staff.Query", database.Log(q, offset, rowsPerPage))

	cafe := []Info{}

	if err := c.db.SelectContext(ctx, &cafe, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting staff")
	}

	return cafe, nil
}

func (c Staff) QueryByID(ctx context.Context, traceID string, claims auth.Claims, staffID string) (Info, error) {
	if _, err := uuid.Parse(staffID); err != nil {
		return Info{}, user.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, user.ErrForbidden
	}

	const q = `SELECT * FROM staffs WHERE staffs_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "staff.QueryByID",
		database.Log(q, staffID),
	)

	var staff Info
	if err := c.db.GetContext(ctx, &staff, q, staffID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, user.ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting staff %q", staffID)
	}

	return staff, nil
}
