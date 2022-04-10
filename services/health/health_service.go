package health

import (
	"arion_shot_api/internal/repository/health"
)

var (
	HealthService healthServiceInterface = &healthService{}
)

type healthService struct{}

type healthServiceInterface interface {
	Check() (bool, error)
}

func (h *healthService) Check() (bool, error) {
	result, err := health.Check()
	if err != nil {
		return false, err
	}
	return result, nil
}
