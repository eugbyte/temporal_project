package billhandler

import "context"

type MessageResponse struct {
	Message string `json:"message"`
}

//encore:api public method=GET path=/healthcheck/bill
func (h *Handler) HealthCheck(ctx context.Context) (*MessageResponse, error) {
	return &MessageResponse{Message: "Hello World"}, nil
}
