package challenge

import "arion_shot_api/internal/domain/user"

type Challenge struct {
	ID               string     `json:"challenge_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Banner           string     `json:"banner"`
	Category         string     `json:"category"`
	EndDate          string     `json:"end_date"`
	Status           string     `json:"status"`
	SubmissionLimit  int        `json:"submission_limit"`
	SubmissionRules  *string    `json:"submission_rules"`
	SubmissionFormat *string    `json:"submission_format"`
	Subscribed       bool       `json:"subscribed"`
	Creator          *user.User `json:"creator"`
}
