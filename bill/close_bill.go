package billhandler

import (
	"context"

	db "encore.app/internal/db/bill"
	temporalbill "encore.app/internal/temporal/bill"
	"go.temporal.io/sdk/client"
)

type CloseBillResp struct {
	Items []string `json:"items"`
	Total float64  `json:"total"`
}

//encore:api public method=PUT path=/close/bill/:billID
func (h *Handler) CloseBill(ctx context.Context, billID string) (*CloseBillResp, error) {
	options := client.StartWorkflowOptions{
		ID:        genWorkFlowID(billID),
		TaskQueue: taskQ,
	}

	workflows := temporalbill.NewWorkFlow(h.billService)
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

	resp := CloseBillResp{
		Items: make([]string, 0),
		Total: 0,
	}

	for _, detail := range bill.Transactions {
		item := detail.ItemName
		amount := detail.Amount

		resp.Items = append(resp.Items, item)
		resp.Total += amount
	}

	return &resp, nil
}
