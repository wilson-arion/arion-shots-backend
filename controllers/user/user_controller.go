package user

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"

	domain "arion_shot_api/internal/domain/user"
	"arion_shot_api/internal/platform/auth"
	service "arion_shot_api/internal/services/user"
	"github.com/pkg/errors"
	"github.com/wgarcia4190/web-router-go/web"
	"go.opencensus.io/trace"
)

type Users struct {
	Authenticator *auth.Authenticator
}

func (users *Users) Login(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Users.Token")
	defer span.End()

	values, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewRequestError(errors.New("web values missing from context"), http.StatusInternalServerError)
	}

	var loginRequest domain.LoginRequest
	if err := web.Decode(request, &loginRequest); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if err := loginRequest.Validate(); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	result, err := service.UserService.Login(loginRequest, values.Start, users.Authenticator)
	if err != nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, writer, result, http.StatusOK)
}

func (users *Users) Exist(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
	email := chi.URLParam(request, "email")

	result := service.UserService.Exist(email)

	return web.Respond(ctx, writer, &domain.UserExistResponse{
		Exist: result,
	}, http.StatusOK)
}

func (users *Users) Register(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Users.Token")
	defer span.End()

	values, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewRequestError(errors.New("web values missing from context"), http.StatusInternalServerError)
	}

	var registerRequest domain.RegisterRequest
	if err := web.Decode(request, &registerRequest); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	if err := registerRequest.Validate(); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	result, err := service.UserService.Register(registerRequest, values.Start, users.Authenticator)
	if err != nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, writer, result, http.StatusOK)
}
