package billhandler

import "context"

type HealthCheckResponse struct {
	Message string `json:"message"`
}

//encore:api public method=GET path=/healthcheck/bill
func (h *Handler) HealthCheck(ctx context.Context) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Message: "Hello World"}, nil
}
