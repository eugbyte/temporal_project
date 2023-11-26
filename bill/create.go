package billhandler

import (
	"context"

	db "encore.app/internal/db/bill"
	"encore.app/internal/temporal/bill/workflows"
	"go.temporal.io/sdk/client"
)

type CreateResponse struct {
	BillID string `json:"billID"`
}

//encore:api public method=POST path=/bill/:billID
func (h *Handler) Create(ctx context.Context, billID string) (*db.Bill, error) {
	logger.Info("billID: ", billID)

	options := client.StartWorkflowOptions{
		ID:        genWorkFlowID(billID),
		TaskQueue: taskQ,
	}

	we, err := h.client.ExecuteWorkflow(ctx, options, workflows.CreateBill, billID)
	if err != nil {
		return nil, err
	}

	logger.Info("started workflow. ", "id: ", we.GetID(), ". run_id:", we.GetRunID())
	var bill db.Bill
	err = we.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}
	return &bill, err
}
