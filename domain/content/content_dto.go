package content

type CreateContentRequest struct {
	URL         string `json:"url"`
	ChallengeID string `json:"challenge_id"`
	OwnerID     string `json:"owner_id"`
}

type CreateContentResponse struct {
	ContentURL string `json:"content_url"`
}
