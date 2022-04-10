package challenge

import (
	"context"
	"net/http"

	"arion_shot_api/domain/challenge"
	"arion_shot_api/platform/auth"
	service "arion_shot_api/services/challenge"
	"github.com/go-chi/chi"
	"github.com/wgarcia4190/web-router-go/web"
)

func GetChallenges(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims not in context")
	}

	challenges, err := service.ChallengeService.GetChallenges(claims.Subject)
	if err != nil {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, writer, challenges, http.StatusOK)
}

func JoinToChallenge(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {
	challengeID := chi.URLParam(request, "id")

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("auth claims not in context")
	}

	r := &challenge.JoinToChallengeRequest{
		UserID:      claims.Subject,
		ChallengeID: challengeID,
	}

	result, err := service.ChallengeService.JoinToChallenge(r)
	if err != nil || !result {
		return web.NewRequestError(err, http.StatusInternalServerError)
	}

	return web.Respond(ctx, writer, &challenge.JoinToChallengeResponse{
		Joined: result,
	}, http.StatusOK)
}
