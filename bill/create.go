package billhandler

import (
	"context"

	temporalbill "encore.app/internal/temporal/bill"
	"go.temporal.io/sdk/client"
)

type CreateResponse struct {
	BillID string
}

//encore:api public path=/bill/:billID
func (h *Handler) Create(ctx context.Context, billID string) (*CreateResponse, error) {
	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: taskQ,
	}

	workflows := temporalbill.NewWorkFlow(h.billService)
	we, err := h.client.ExecuteWorkflow(ctx, options, workflows.Create, billID)
	if err != nil {
		return nil, err
	}

	err = we.Get(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &CreateResponse{BillID: "Hello World"}, nil
}
