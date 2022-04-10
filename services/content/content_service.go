package content

import (
	"arion_shot_api/internal/domain/content"
	repository "arion_shot_api/internal/repository/content"
	"github.com/pkg/errors"
)

var (
	ContentService contentServiceInterface = &contentService{}
)

type contentService struct{}

type contentServiceInterface interface {
	UpdateVote(request *content.VoteRequest) (*content.VoteResponse, error)
	CreateContent(request *content.CreateContentRequest) (*content.CreateContentResponse, error)
}

func (service *contentService) UpdateVote(request *content.VoteRequest) (*content.VoteResponse, error) {
	switch request.Action {
	case content.ADD:
		result, err := repository.ContentRepository.AddVote(request.ContentID, *request.VoterID)
		if err != nil {
			return nil, err
		}

		return result, nil
	case content.REMOVE:
		result, err := repository.ContentRepository.DeleteVote(request.ContentID, *request.VoterID)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, errors.New("Unrecognized action")
}

func (service *contentService) CreateContent(request *content.CreateContentRequest) (*content.CreateContentResponse, error) {
	result, err := repository.ContentRepository.CreateContent(request.ChallengeID, request.OwnerID, request.URL)
	if err != nil {
		return nil, err
	}

	return result, nil
}
