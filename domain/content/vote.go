package content

type VoteRequest struct {
	ContentID string  `json:"content_id"`
	VoterID   *string `json:"voter_id"`
	Action    string  `json:"action"`
}

type VoteResponse struct {
	TotalVotes int `json:"total_votes"`
}

const (
	ADD    string = "ADD"
	REMOVE        = "REMOVE"
)
