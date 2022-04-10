package schema

import "database/sql"

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Using a constant in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.

const seeds = `
-- Create admin User with password "Admin123456"
INSERT INTO users
    (user_id, first_name, last_name, email, pass, user_role, date_created, date_updated)
VALUES
    (UUID_TO_BIN('d2a9266b-a928-11ec-814b-0242ac160002'), 'John', 'Doe', 'johndoe@arionkoder.com', MD5('Admin123456'),
     "A", now(), now());

-- Create challenge
INSERT INTO challenges
 	(challenge_id, title, description, banner, creator_id, category, end_date, status,
 	 submission_limit, submission_rules, submission_format, date_created, date_updated)
 VALUES
 	(UUID_TO_BIN('c44d9587-ac83-11ec-adb9-0242ac163333'), 'Inspired by Nature', 'Upload nature pictures',
 	 'https://res.cloudinary.com/arionshots/image/upload/v1647875975/arion_shots_picture_contest/EFFECTS_nc7vs7.jpg',
 	 UUID_TO_BIN('d2a9266b-a928-11ec-814b-0242ac160002'),'Nature', now(), "A",
 	 DEFAULT, NULL, DEFAULT, now(), now());

-- Create content
INSERT INTO contents
    (content_id, url, challenge_id, owner_id, date_created, date_updated)
VALUES
    (UUID_TO_BIN('f9d0d858-ac88-11ec-adb9-0242ac160002'),
    'https://res.cloudinary.com/arionshots/image/upload/v1647875975/arion_shots_picture_contest/EFFECTS_nc7vs7.jpg',
    UUID_TO_BIN('c44d9587-ac83-11ec-adb9-0242ac163333'), UUID_TO_BIN('d2a9266b-a928-11ec-814b-0242ac160002'), now(), now());

INSERT INTO content_voters
	(content_id, voter_id)
VALUES
	(UUID_TO_BIN('f9d0d858-ac88-11ec-adb9-0242ac160002'), UUID_TO_BIN('d2a9266b-a928-11ec-814b-0242ac160002'));
`

func Seed(db *sql.DB) error {
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
