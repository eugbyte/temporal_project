package billhandler

import (
	"context"
	"encoding/json"
	"fmt"

	customerrors "encore.app/internal/custom_errors"
	db "encore.app/internal/db/bill"
	workflows "encore.app/internal/temporal/bill/workflow"
	"go.temporal.io/sdk/client"
)

type IncreaseBillResp struct {
	BillID     string
	WorkflowID string
}

//encore:api public method=PUT path=/bill/:billID
func (h *Handler) IncreaseBill(ctx context.Context, billID string, billDetail db.TransactionDetail) (*IncreaseBillResp, error) {
	byts, _ := json.MarshalIndent(billDetail, "", "\t")
	fmt.Println(string(byts))

	// Convert the currency to USD
	currency := billDetail.Amount.CurrencyCode()
	if _, ok := h.currencies[currency]; !ok {
		return nil, customerrors.NewAppError("currency not recognised")
	}

	usd, err := billDetail.Amount.Convert("USD", h.currencies[currency])
	if err != nil {
		return nil, customerrors.NewAppError("currency conversion failed")
	}
	billDetail.Amount = usd

	// Start the workflow to increase the bill, which will continue running pending user confirmation via ConfirmBillIncrease
	workflowID := genWorkFlowID(billID)
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQ,
	}

	_, err = h.client.ExecuteWorkflow(ctx, options, workflows.IncreaseBill, billID, billDetail)
	if err != nil {
		return nil, err
	}

	return &IncreaseBillResp{
		BillID:     billID,
		WorkflowID: workflowID,
	}, nil
}

//encore:api public method=GET path=/confirm/bill/:billID/:workflowID
func (h *Handler) ConfirmBillIncrease(ctx context.Context, billID string, workflowID string) (*MessageResponse, error) {
	runId := "" // we did not store runId we can safely leave it empty
	confirmed := true
	err := h.client.SignalWorkflow(ctx, workflowID, runId, workflows.SignalChannel, confirmed)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{
		Message: "invoiced confirmed",
	}, nil

}
