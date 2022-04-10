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
    (0, 'John', 'Doe', 'johndoe@arionkoder.com', MD5('Admin123456'),
     "A", now(), now());

-- Create challenge
INSERT INTO challenges
 	(challenge_id, title, description, banner, creator_id, category, end_date, status,
 	 submission_limit, submission_rules, submission_format, date_created, date_updated)
 VALUES
 	(0, 'Inspired by Nature', 'Upload nature pictures',
 	 'https://res.cloudinary.com/arionshots/image/upload/v1647875975/arion_shots_picture_contest/EFFECTS_nc7vs7.jpg',
 	 1,'Nature', now(), "A",
 	 DEFAULT, NULL, DEFAULT, now(), now());

-- Create content
INSERT INTO contents
    (content_id, url, challenge_id, owner_id, date_created, date_updated)
VALUES
    0,
    'https://res.cloudinary.com/arionshots/image/upload/v1647875975/arion_shots_picture_contest/EFFECTS_nc7vs7.jpg',
    1, 1, now(), now());

INSERT INTO content_voters
	(content_id, voter_id)
VALUES
	(1, 1);
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
