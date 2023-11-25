package billhandler

import "context"

type HealthCheckResponse struct {
	Message string
}

//encore:api public path=/bill/healthcheck
func (h *Handler) HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Message: "Hello World"}, nil
}
