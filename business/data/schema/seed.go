package schema

import "github.com/jmoiron/sqlx"

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

const seeds = `
INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO cafes (cafe_id, name, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'rafe', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'damdes', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO staffs (staff_id, position, fio, user_id, cafe_id, date_created, date_updated) VALUES
	('dfa196ed-3496-404f-8d6b-23da5b71445b', 'saler', 'Kabir Salauat', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('448a2f2e-c348-4d8a-9a42-ae22dd8398ac', 'cashier', 'anna petrova', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;`

const deleteAll = `
DELETE FROM staffs;
DELETE FROM cafes;
DELETE FROM users;`

func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
