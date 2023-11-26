package workflows

import (
	db "encore.app/internal/db/bill"
	debug "encore.app/internal/logger"
	"encore.app/internal/temporal/bill/activities"
	"go.temporal.io/sdk/workflow"
)

var logger = debug.Logger

func CreateBill(ctx workflow.Context, billID string) (db.Bill, error) {
	logger.Info("creating...", billID)
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	var bill db.Bill
	err := workflow.ExecuteActivity(ctx, activities.CreateBillActivity, billID).Get(ctx, &bill)
	return bill, err
}

func IncreaseBill(ctx workflow.Context, billID string, billDetail db.TransactionDetail) error {
	logger.Info("starting increase bill")
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Wait for confirmation before adding to invoice
	selector := workflow.NewSelector(ctx)
	signalCh := workflow.GetSignalChannel(ctx, SignalChannel)

	var confirmed bool = false
	// implement selector reciever via signal channel
	selector.AddReceive(signalCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &confirmed)
	})

	// blocks untill a signal is received
	logger.Info("waiting...")
	selector.Select(ctx)
	logger.Info("signal received")

	// If confirmed, add invoice
	if !confirmed {
		logger.Info("confirmation denied")
		return nil
	}
	return workflow.ExecuteActivity(ctx, activities.IncreaseBillActivity, billID, billDetail).Get(ctx, nil)
}

func CloseBill(ctx workflow.Context, billID string) (db.Bill, error) {
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	var bill db.Bill
	err := workflow.ExecuteActivity(ctx, activities.CloseBillActivity, billID).Get(ctx, &bill)
	return bill, err
}

func SanityCheck(ctx workflow.Context) error {
	logger.Info("sanity check")
	ctx = workflow.WithActivityOptions(ctx, options)
	return workflow.ExecuteActivity(ctx, activities.SanityCheckActivity).Get(ctx, nil)
}
