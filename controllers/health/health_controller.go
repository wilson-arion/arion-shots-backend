package health

import (
	"arion_shot_api/services/health"
	"context"
	"github.com/wgarcia4190/web-router-go/web"

	"net/http"
)

// Health responds with a 200 OK if the service is healthy and ready for traffic
func Health(ctx context.Context, writer http.ResponseWriter, request *http.Request) error {

	var h struct {
		Status string `json:"status"`
	}

	result, err := health.HealthService.Check()
	if err != nil {
		h.Status = "db not ready"
		return web.Respond(ctx, writer, h, http.StatusInternalServerError)
	}

	if !result {
		h.Status = "db not ready"
		return web.Respond(ctx, writer, h, http.StatusInternalServerError)
	}

	h.Status = "ok"
	return web.Respond(ctx, writer, h, http.StatusOK)
}
