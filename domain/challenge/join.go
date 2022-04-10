package challenge

type JoinToChallengeRequest struct {
	UserID      string `json:"user_id"`
	ChallengeID string `json:"challenge_id"`
}

type JoinToChallengeResponse struct {
	Joined bool `json:"joined"`
}
