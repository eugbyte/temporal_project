package billhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	db "encore.app/internal/db/bill"
	workflows "encore.app/internal/temporal/bill/workflow"
	"encore.dev/beta/errs"
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

	// check whether bill is closed
	bill, err := h.billService.Get(billID)
	if err != nil {
		return nil, errs.Convert(err)
	}
	if bill.Status == db.CLOSED {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "bill is already closed",
		}
	}

	// Convert the currency to USD
	currency := billDetail.Amount.CurrencyCode()
	if _, ok := h.currencyRates[currency]; !ok {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "currency is not recognised",
		}
	}

	f, err := strconv.ParseFloat(h.currencyRates[currency], 64)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: err.Error(),
		}
	}

	usd, err := billDetail.Amount.Convert("USD", fmt.Sprintf("%f", 1/f))
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "currency conversion failed: " + err.Error(),
		}
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
		return nil, errs.Convert(err)
	}
	return &MessageResponse{
		Message: "invoiced confirmed",
	}, nil

}
