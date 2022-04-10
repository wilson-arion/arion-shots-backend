package user

type User struct {
	ID          string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Role        string `json:"user_role"`
	Password    string `json:"-"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}
