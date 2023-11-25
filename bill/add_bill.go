package billhandler

import (
	"context"

	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"go.temporal.io/sdk/client"
)

type AddBillResp struct {
	BillID     string
	WorkflowID string
}

//encore:api public method=PUT path=/bill/:billID
func (h *Handler) AddBill(ctx context.Context, billID string, transactionDetail db.TransactionDetail) (*AddBillResp, error) {
	logger.Info("PUT:", transactionDetail)

	workflowID := genWorkFlowID(billID)
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQ,
	}

	workflows := temporalbill.NewWorkFlow(h.billService)
	_, err := h.client.ExecuteWorkflow(ctx, options, workflows.AddBill, billID, transactionDetail)
	if err != nil {
		return nil, err
	}

	return &AddBillResp{
		BillID:     billID,
		WorkflowID: workflowID,
	}, nil
}

//encore:api public method=GET path=/confirm-bill/:billID/:workflowID
func (h *Handler) Confirm(ctx context.Context, billID string, workflowID string) (*MessageResponse, error) {
	runId := "" // we did not store runId we can safely leave it empty
	confirmed := true
	err := h.client.SignalWorkflow(ctx, workflowID, runId, temporalbill.SignalChannel, confirmed)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{
		Message: "invoiced confirmed",
	}, nil

}
