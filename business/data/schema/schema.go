package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

var migrations = []darwin.Migration{
	{
		Version:     1.1,
		Description: "Create table users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,
	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     1.2,
		Description: "Create table cafes",
		Script: `
CREATE TABLE cafes (
	cafe_id      UUID,
	name         TEXT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (supermarket_id)
);`,
	},
	{
		Version:     1.3,
		Description: "Create table staff",
		Script: `
CREATE TABLE staffs (
	staff_id     UUID,
	position  	 TEXT,
	fio  		 TEXT,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,
	PRIMARY KEY (category_id)
);`,
	},
}
