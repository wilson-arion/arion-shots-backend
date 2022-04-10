package schema

import (
	"database/sql"
	"github.com/GuiaBolso/darwin"
)

var (
	migrations = []darwin.Migration{
		{
			Version:     1,
			Description: "Add users",
			Script: `
					CREATE TABLE IF NOT EXISTS users(
						user_id         MEDIUMINT NOT NULL AUTO_INCREMENT,
						first_name      VARCHAR(45) NULL,
						last_name       VARCHAR(45) NULL,
						email           VARCHAR(45) NOT NULL,
						pass            VARCHAR(255) NOT NULL,
                        user_role       VARCHAR(1) NOT NULL,
						date_created    DATETIME NOT NULL,
						date_updated    DATETIME NOT NULL,
						PRIMARY KEY(user_id),
						UNIQUE INDEX email_UNIQUE (email ASC)
					);
			`,
		},
		{
			Version:     2,
			Description: "Add challenges",
			Script: `
                    CREATE TABLE IF NOT EXISTS challenges (
                        challenge_id        MEDIUMINT NOT NULL AUTO_INCREMENT,
                        title               VARCHAR(45) NOT NULL,
                        description         TEXT NOT NULL,
                        banner              VARCHAR(2083) NOT NULL,
                        creator_id          MEDIUMINT NOT NULL,
                        category            VARCHAR(45) NOT NULL,
                        end_date            DATETIME NOT NULL,
                        status              VARCHAR(1) NOT NULL,
                        submission_limit    INT DEFAULT 0,
                        submission_rules    TEXT NULL,
                        submission_format   TEXT NULL,
                        date_created        DATETIME NOT NULL,
						date_updated        DATETIME NOT NULL,
                        PRIMARY KEY (challenge_id),
                        INDEX (challenge_id),
                        FOREIGN KEY (creator_id) REFERENCES users(user_id) ON DELETE CASCADE
                    );
            `,
		},
		{
			Version:     3,
			Description: "Add content",
			Script: `
                    CREATE TABLE IF NOT EXISTS contents (
                        content_id          MEDIUMINT NOT NULL AUTO_INCREMENT,
                        url                 VARCHAR(2083) NOT NULL,
                        challenge_id        MEDIUMINT NOT NULL,
                        owner_id            MEDIUMINT NOT NULL,
                        date_created        DATETIME NOT NULL,
                        date_updated        DATETIME NOT NULL,
                        PRIMARY KEY (content_id),
                        INDEX (content_id, challenge_id, owner_id),
                        FOREIGN KEY (challenge_id) REFERENCES challenges(challenge_id) ON DELETE CASCADE,
                        FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE
                    );
            `,
		},
		{
			Version:     4,
			Description: "Add content and vote",
			Script: `
                    CREATE TABLE IF NOT EXISTS content_voters (
                        content_id  MEDIUMINT NOT NULL,
                        voter_id    MEDIUMINT NOT NULL,
                        PRIMARY KEY (content_id, voter_id),
                        INDEX (content_id, voter_id),
                        FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE,
                        FOREIGN KEY (voter_id) REFERENCES users(user_id) ON DELETE CASCADE
                    );
            `,
		},
		{
			Version:     5,
			Description: "Add challenge and user",
			Script: `
                    CREATE TABLE IF NOT EXISTS user_challenges_joined (
                        user_id         MEDIUMINT NOT NULL,
                        challenge_id    MEDIUMINT NOT NULL,
                        PRIMARY KEY (user_id, challenge_id),
                        INDEX (user_id, challenge_id),
                        FOREIGN KEY (challenge_id) REFERENCES challenges(challenge_id) ON DELETE CASCADE,
                        FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
                    );
            `,
		},
		{
			Version:     6,
			Description: "Add comments",
			Script: `
                    CREATE TABLE IF NOT EXISTS comments (
                        comment_id      MEDIUMINT NOT NULL AUTO_INCREMENT,
                        creator_id      MEDIUMINT NOT NULL,
                        content_id      MEDIUMINT NOT NULL,
                        comment         TEXT NOT NULL,
                        date_created    DATETIME NOT NULL,
                        PRIMARY KEY (comment_id),
                        INDEX (creator_id, content_id),
                        FOREIGN KEY (content_id) REFERENCES contents(content_id) ON DELETE CASCADE,
                        FOREIGN KEY (creator_id) REFERENCES users(user_id) ON DELETE CASCADE
                    );
            `,
		},
	}
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.MySQLDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}
