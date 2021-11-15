package cafe

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

type Cafe struct {
	log *log.Logger
	db  *sqlx.DB
}

func New(log *log.Logger, db *sqlx.DB) Cafe {
	return Cafe{
		log: log,
		db:  db,
	}
}

func (c Cafe) Create(ctx context.Context, traceID string, nc NewCafe, now time.Time) (Info, error) {
	cafe := Info{
		ID:   uuid.New().String(),
		Name: nc.Name,

		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO cafes (cafe_id, name, date_created, date_updated)
			   VAlUES ($1, $2, $3, $4)`

	c.log.Printf("%s : %s query : %s", traceID, "cafe.Create",
		database.Log(q, cafe.ID, cafe.Name, cafe.DateCreated, cafe.DateUpdated),
	)

	if _, err := c.db.ExecContext(ctx, q, cafe.ID, cafe.Name, cafe.DateCreated, cafe.DateUpdated); err != nil {
		return Info{}, errors.Wrap(err, "inserting cafe")
	}
	return cafe, nil
}

func (c Cafe) Update(ctx context.Context, traceID string, claims auth.Claims, cafeID string, uc UpdateCafe, now time.Time) error {
	cafe, err := c.QueryByID(ctx, traceID, claims, cafeID)
	if err != nil {
		return err
	}

	if uc.Name != nil {
		cafe.Name = *uc.Name
	}
	cafe.DateUpdated = now

	const q = `
	UPDATE
		cafes
	SET 
		"name" = $2,
		"date_updated" = $3
	WHERE
		cafe_id = $1`

	c.log.Printf("%s: %s: %s", traceID, "cafe.Update",
		database.Log(q, cafe.ID, cafe.Name, cafe.DateCreated, cafe.DateUpdated),
	)

	if _, err = c.db.ExecContext(ctx, q, cafe.ID, cafe.Name, cafe.DateUpdated); err != nil {
		return errors.Wrap(err, "updating cafe")
	}

	return nil
}

func (c Cafe) Delete(ctx context.Context, traceID string, cafeID string) error {
	if _, err := uuid.Parse(cafeID); err != nil {
		return user.ErrInvalidID
	}
	const q = `DELETE FROM cafes where cafe_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "cafe.Delete",
		database.Log(q, cafeID),
	)

	if _, err := c.db.ExecContext(ctx, q, cafeID); err != nil {
		return errors.Wrapf(err, "deleting cafe %s", cafeID)
	}
	return nil
}

func (c Cafe) Query(ctx context.Context, traceID string, pageNumber, rowsPerPage int) ([]Info, error) {
	const q = `SELECT * FROM cafes ORDER BY cafe_id OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`
	offset := (pageNumber - 1) * rowsPerPage

	c.log.Printf("%s : %s query : %s", traceID, "cafe.Query", database.Log(q, offset, rowsPerPage))

	cafe := []Info{}

	if err := c.db.SelectContext(ctx, &cafe, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting cafe")
	}

	return cafe, nil
}

func (c Cafe) QueryByID(ctx context.Context, traceID string, claims auth.Claims, cafeID string) (Info, error) {
	if _, err := uuid.Parse(cafeID); err != nil {
		return Info{}, user.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) {
		return Info{}, user.ErrForbidden
	}

	const q = `SELECT * FROM cafes WHERE cafe_id = $1`

	c.log.Printf("%s : %s query : %s", traceID, "cafe.QueryByID",
		database.Log(q, cafeID),
	)

	var cafe Info
	if err := c.db.GetContext(ctx, &cafe, q, cafeID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, user.ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting cafe %q", cafeID)
	}

	return cafe, nil
}
