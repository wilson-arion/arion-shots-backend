package user

import (
	"time"

	"arion_shot_api/internal/domain/user"
	"arion_shot_api/internal/platform/auth"
	repository "arion_shot_api/internal/repository/user"
	"github.com/pkg/errors"
)

var (
	UserService userServiceInterface = &userService{}
)

type userService struct{}

type userServiceInterface interface {
	Login(request user.LoginRequest, now time.Time, authenticator *auth.Authenticator) (*user.LoginResponse, error)
	Register(request user.RegisterRequest, now time.Time, authenticator *auth.Authenticator) (*user.LoginResponse, error)
	Exist(email string) bool
}

func (service *userService) Login(request user.LoginRequest, now time.Time, authenticator *auth.Authenticator) (*user.LoginResponse, error) {
	u, err := repository.UserRepository.FindByEmailAndPassword(request)

	if err != nil {
		return nil, err
	}

	var token string
	var claims auth.Claims

	claims, err = repository.UserRepository.Authenticate(u.ID, u.Role, now)

	if err != nil {
		return nil, errors.Wrap(err, "authenticating")
	}

	token, err = authenticator.GenerateToken(claims)
	if err != nil {
		return nil, errors.Wrap(err, "generating token")
	}

	response := &user.LoginResponse{
		User:  *u,
		Token: token,
	}

	return response, nil
}

func (service *userService) Register(request user.RegisterRequest, now time.Time, authenticator *auth.Authenticator) (*user.LoginResponse, error) {
	u, err := repository.UserRepository.CreateUser(request)

	if err != nil {
		return nil, err
	}

	var token string
	var claims auth.Claims

	claims, err = repository.UserRepository.Authenticate(u.ID, u.Role, now)

	if err != nil {
		return nil, errors.Wrap(err, "authenticating")
	}

	token, err = authenticator.GenerateToken(claims)
	if err != nil {
		return nil, errors.Wrap(err, "generating token")
	}

	response := &user.LoginResponse{
		User:  *u,
		Token: token,
	}

	return response, nil
}

func (service *userService) Exist(email string) bool {
	result := repository.UserRepository.EmailExist(email)

	return result
}
