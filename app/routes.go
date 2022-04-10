package app

import (
	"arion_shot_api/internal/controllers/challenge"
	"arion_shot_api/internal/controllers/content"
	"arion_shot_api/internal/controllers/health"
	"arion_shot_api/internal/controllers/user"
	"arion_shot_api/internal/platform/middleware"
	"log"
	"net/http"
	"os"

	"arion_shot_api/internal/platform/auth"
	"github.com/wgarcia4190/web-router-go/web"
)

// API constructs a handler that knows about all API routes.
func API(shutdown chan os.Signal, logger *log.Logger, authenticator *auth.Authenticator) http.Handler {
	router := web.NewRouter(shutdown, logger, middleware.Logger(logger), middleware.Errors(logger), middleware.Panics())

	router.Handler(http.MethodGet, "/v1/health", health.Health)

	u := user.Users{
		Authenticator: authenticator,
	}
	router.Handler(http.MethodPost, "/v1/token", u.Login)
	router.Handler(http.MethodPost, "/v1/register", u.Register)
	router.Handler(http.MethodGet, "/v1/user/exist/{email}", u.Exist)

	router.Handler(http.MethodGet, "/v1/challenges", challenge.GetChallenges, middleware.Authenticate(authenticator))
	router.Handler(http.MethodPut, "/v1/challenge/{id}/join", challenge.JoinToChallenge, middleware.Authenticate(authenticator))

	router.Handler(http.MethodPost, "/v1/content/vote", content.UpdateVote, middleware.Authenticate(authenticator))
	router.Handler(http.MethodPost, "/v1/content/create", content.CreateContent, middleware.Authenticate(authenticator))

	return router
}
