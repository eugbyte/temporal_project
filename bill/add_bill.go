package billhandler

import (
	"context"

	customerrors "encore.app/internal/custom_errors"
	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"go.temporal.io/sdk/client"
)

//encore:api public method=PUT path=/bill/:billID
func (h *Handler) AddBill(ctx context.Context, billID string, transactionDetail db.TransactionDetail) (*MessageResponse, error) {
	logger.Info("PUT:", transactionDetail)

	options := client.StartWorkflowOptions{
		ID:        genWorkFlowID("add", billID),
		TaskQueue: taskQ,
	}

	h.workflowIDs.Set(billID, genWorkFlowID("add", billID))

	workflows := temporalbill.NewWorkFlow(h.billService)
	_, err := h.client.ExecuteWorkflow(ctx, options, workflows.CreateBill, billID)
	if err != nil {
		return nil, err
	}

	return &MessageResponse{
		Message: "transaction to add to bill started, awaiting confirmation",
	}, nil
}

func (h *Handler) Confirm(ctx context.Context, billID string, confirmed bool) (*MessageResponse, error) {
	workflowId, ok := h.workflowIDs.Get(billID)
	if !ok {
		return nil, customerrors.NewAppError("workflow id not found")
	}

	runId := ""
	err := h.client.SignalWorkflow(ctx, workflowId, runId, temporalbill.SignalChannel, confirmed)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{
		Message: "invoiced confirmed",
	}, nil

}
