package challenge

import (
    "arion_shot_api/internal/domain/challenge"
    repository "arion_shot_api/internal/repository/challenge"
)

var (
    ChallengeService challengeServiceInterface = &challengeService{}
)

type challengeService struct{}

type challengeServiceInterface interface {
    GetChallenges(userId string) ([]challenge.Challenge, error)
    JoinToChallenge(request *challenge.JoinToChallengeRequest) (bool, error)
}

func (service *challengeService) GetChallenges(userId string) ([]challenge.Challenge, error) {
    challenges, err := repository.ChallengeRepository.GetChallengesPerUser(userId)

    if err != nil {
        return []challenge.Challenge{}, err
    }

    return challenges, nil
}

func (service *challengeService) JoinToChallenge(request *challenge.JoinToChallengeRequest) (bool, error) {
    result, err := repository.ChallengeRepository.JoinToChallenge(request)

    if err != nil {
        return false, err
    }

    return result, nil
}
