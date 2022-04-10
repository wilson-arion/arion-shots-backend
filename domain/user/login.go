package user

import (
	"arion_shot_api/internal/utils/crypto"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const (
	ArionDomain = "arionkoder.com"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"user_role"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserExistResponse struct {
	Exist bool `json:"exist"`
}

func (request LoginRequest) GetEncryptedPassword() string {
	pass, _ := crypto.GetMd5(request.Password)
	return pass
}

func (request RegisterRequest) GetEncryptedPassword() string {
	pass, _ := crypto.GetMd5(request.Password)
	return pass
}

func (request LoginRequest) Validate() error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if !re.MatchString(request.Email) {
		return errors.New("invalid email address")
	}

	atIndex := strings.LastIndex(request.Email, "@")
	domain := request.Email[atIndex+1:]
	if domain != ArionDomain {
		return errors.New("Invalid email address. You need an arionkoder email")
	}

	return nil
}

func (request RegisterRequest) Validate() error {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	if !re.MatchString(request.Email) {
		return errors.New("invalid email address")
	}

	atIndex := strings.LastIndex(request.Email, "@")
	domain := request.Email[atIndex+1:]
	if domain != ArionDomain {
		return errors.New("Invalid email address. You need an arionkoder email")
	}

	return nil
}
