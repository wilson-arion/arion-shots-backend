package challenge

import (
	"arion_shot_api/internal/domain/challenge"
	"arion_shot_api/internal/domain/user"
	sqlutils "arion_shot_api/internal/utils/sql"
	"github.com/pkg/errors"
)

const (
	queryGetChallengesPerUser = `
        SELECT
	        BIN_TO_UUID(challenge_id) as challenge_id ,
	        title,
            description,
            banner,
            category,
            end_date,
            status,
            submission_limit,
            submission_rules,
            submission_format,
            ((
                SELECT
                    COUNT(*) as cnt
                FROM
                    user_challenges_joined ucj
                WHERE
                    ucj.user_id = UUID_TO_BIN(?)
                AND
                    ucj.challenge_id = c.challenge_id
                HAVING cnt > 0
            ) IS NOT NULL) as subscribed,
            BIN_TO_UUID(creator_id) as user_id,
            u.first_name,
            u.last_name,
            u.email,
            u.user_role,
            u.date_created,
            u.date_updated
        FROM
            challenges c
        JOIN
            users u
        ON
            c.creator_id = u.user_id;
    `
	queryJoinToChallenge = `
        INSERT INTO user_challenges_joined
	        (user_id, challenge_id)
        VALUES
	        (UUID_TO_BIN(?), UUID_TO_BIN(?))

    `
)

var (
	ChallengeRepository challengeRepositoryInterface = &challengeRepository{}
)

type challengeRepository struct{}

type challengeRepositoryInterface interface {
	GetChallengesPerUser(userId string) ([]challenge.Challenge, error)
	JoinToChallenge(request *challenge.JoinToChallengeRequest) (bool, error)
}

func (repository *challengeRepository) GetChallengesPerUser(userId string) ([]challenge.Challenge, error) {
	challenges := make([]challenge.Challenge, 0)
	stmt, err := sqlutils.CreateStmt(queryGetChallengesPerUser) //nolint:sqlclosecheck

	if err != nil {
		return challenges, err
	}
	defer sqlutils.CloseStmt(stmt)

	result, err := stmt.Query(userId)
	if err != nil {
		return challenges, err
	}
	defer sqlutils.CloseRows(result)

	for result.Next() {
		u := &user.User{}
		c := &challenge.Challenge{}

		if err := result.Scan(
			&c.ID,
			&c.Title,
			&c.Description,
			&c.Banner,
			&c.Category,
			&c.EndDate,
			&c.Status,
			&c.SubmissionLimit,
			&c.SubmissionRules,
			&c.SubmissionFormat,
			&c.Subscribed,
			&u.ID,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Role,
			&u.DateCreated,
			&u.DateUpdated,
		); err != nil {
			return challenges, sqlutils.ParseError(err)
		}

		c.Creator = u
		challenges = append(challenges, *c)
	}

	if err := result.Err(); err != nil {
		return []challenge.Challenge{}, err
	}

	return challenges, nil
}

func (repository *challengeRepository) JoinToChallenge(request *challenge.JoinToChallengeRequest) (bool, error) {
	stmt, err := sqlutils.CreateStmt(queryJoinToChallenge)

	if err != nil {
		return false, err
	}
	defer sqlutils.CloseStmt(stmt)

	insertResult, insertErr := stmt.Exec(request.UserID, request.ChallengeID)
	if insertErr != nil {
		return false, sqlutils.ParseError(insertErr)
	}

	if insertResult != nil {
		return true, nil
	}

	return false, errors.New("error when trying to save user")
}
