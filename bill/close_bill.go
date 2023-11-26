package billhandler

import (
	"context"

	db "encore.app/internal/db/bill"
	"encore.app/internal/temporal/bill/workflows"
	"github.com/bojanz/currency"
	"go.temporal.io/sdk/client"
)

type CloseBillResp struct {
	Items []string        `json:"items"`
	Total currency.Amount `json:"total"`
}

//encore:api public method=PUT path=/close/bill/:billID
func (h *Handler) CloseBill(ctx context.Context, billID string) (*CloseBillResp, error) {
	options := client.StartWorkflowOptions{
		ID:        genWorkFlowID(billID),
		TaskQueue: taskQ,
	}

	we, err := h.client.ExecuteWorkflow(ctx, options, workflows.CloseBill, billID)
	if err != nil {
		return nil, err
	}

	logger.Info("started workflow. ", "id: ", we.GetID(), ". run_id:", we.GetRunID())
	var bill db.Bill
	err = we.Get(ctx, &bill)
	if err != nil {
		return nil, err
	}

	total, _ := currency.NewAmount("0", "USD")

	resp := CloseBillResp{
		Items: make([]string, 0),
		Total: total,
	}

	for _, detail := range bill.Transactions {
		resp.Items = append(resp.Items, detail.ItemName)

		amount := detail.Amount
		logger.Info(amount.String())
		total, err = resp.Total.Add(amount)
		if err != nil {
			return nil, err
		}
		resp.Total = total
	}

	return &resp, nil
}
