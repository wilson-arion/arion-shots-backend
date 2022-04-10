package content

import (
	"arion_shot_api/internal/domain/content"
	sqlutils "arion_shot_api/internal/utils/sql"
	"github.com/pkg/errors"
)

const (
	queryAddVote = `
        INSERT INTO content_voters
	        (content_id, voter_id)
        VALUES
	        (UUID_TO_BIN(?), UUID_TO_BIN(?));
    `
	queryDeleteVote = `
        DELETE FROM
            content_voters
        WHERE
            content_id = UUID_TO_BIN(?)
        AND
            voter_id = UUID_TO_BIN(?);
    `

	queryCountVote = `
        SELECT
            COUNT(*) as votes
        FROM
            content_voters
        WHERE
            content_id = UUID_TO_BIN(?)
        LIMIT 1;
    `

	queryCreateContent = `
        INSERT INTO contents
            (content_id, url, challenge_id, owner_id, date_created, date_updated)
        VALUES
            (UUID_TO_BIN(UUID()), ?, UUID_TO_BIN(?), UUID_TO_BIN(?), now(), now());

    `

	queryGetContentsPerChallenge = `
        SELECT
            BIN_TO_UUID(c.content_id) as content_id,
            c.url,
            BIN_TO_UUID(c.challenge_id) as challenge_id,
            c.date_created,
            c.date_updated,
            (
                SELECT
                    COUNT(*) as cnt
                FROM
                   content_voters cv
                WHERE
                    cv.content_id = c.content_id
            ) as votes,
            BIN_TO_UUID(c.owner_id) as owner_id,
            u.first_name,
            u.last_name,
            u.email,
            u.user_role,
            u.date_created,
            u.date_updated
        FROM
            contents c
        JOIN
            users u
        ON
            c.owner_id = u.user_id;
        WHERE c.challenge_id = UUID_TO_BIN(?);
    `
)

var (
	ContentRepository contentRepositoryInterface = &contentRepository{}
)

type contentRepository struct{}

type contentRepositoryInterface interface {
	AddVote(contentId string, voterId string) (*content.VoteResponse, error)
	DeleteVote(contentId string, voterId string) (*content.VoteResponse, error)
	CreateContent(challengeId string, ownerId string, url string) (*content.CreateContentResponse, error)
	getCount(contentId string) (int, error)
}

func (repository *contentRepository) AddVote(contentId string, voterId string) (*content.VoteResponse, error) {
	stmt, err := sqlutils.CreateStmt(queryAddVote)

	if err != nil {
		return nil, err
	}
	defer sqlutils.CloseStmt(stmt)

	insertResult, insertErr := stmt.Exec(contentId, voterId)
	if insertErr != nil {
		return nil, sqlutils.ParseError(insertErr)
	}

	if insertResult != nil {
		count, err := repository.getCount(contentId)
		if err != nil {
			return nil, err
		}

		response := &content.VoteResponse{
			TotalVotes: count,
		}

		return response, nil
	}

	return nil, errors.New("error when trying to add a vote")
}

func (repository *contentRepository) DeleteVote(contentId string, voterId string) (*content.VoteResponse, error) {
	stmt, err := sqlutils.CreateStmt(queryDeleteVote)

	if err != nil {
		return nil, err
	}
	defer sqlutils.CloseStmt(stmt)

	deleteResult, deleteErr := stmt.Exec(contentId, voterId)
	if deleteErr != nil {
		return nil, sqlutils.ParseError(deleteErr)
	}

	if deleteResult != nil {
		count, err := repository.getCount(contentId)
		if err != nil {
			return nil, err
		}

		response := &content.VoteResponse{
			TotalVotes: count,
		}

		return response, nil
	}

	return nil, errors.New("error when trying to delete a vote")
}

func (repository *contentRepository) CreateContent(challengeId string, ownerId string, url string) (*content.CreateContentResponse, error) {
	stmt, err := sqlutils.CreateStmt(queryCreateContent)

	if err != nil {
		return nil, err
	}
	defer sqlutils.CloseStmt(stmt)

	insertResult, insertErr := stmt.Exec(url, challengeId, ownerId)
	if insertErr != nil {
		return nil, sqlutils.ParseError(insertErr)
	}

	if insertResult != nil {
		response := &content.CreateContentResponse{
			ContentURL: url,
		}

		return response, nil
	}

	return nil, errors.New("error when trying to add a content")

}

func (repository *contentRepository) getCount(contentId string) (int, error) {
	stmt, err := sqlutils.CreateStmt(queryCountVote)

	if err != nil {
		return 0, err
	}
	defer sqlutils.CloseStmt(stmt)

	var count int

	result := stmt.QueryRow(contentId)
	if err := result.Scan(&count); err != nil {
		return 0, sqlutils.ParseError(err)
	}

	return count, nil
}
