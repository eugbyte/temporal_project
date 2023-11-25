package temporalbill

import (
	"time"

	db "encore.app/internal/db/bill"
	debug "encore.app/internal/logger"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var logger = debug.Logger

var retrypolicy = &temporal.RetryPolicy{
	InitialInterval:        time.Second,
	BackoffCoefficient:     2.0,
	MaximumInterval:        100 * time.Second,
	MaximumAttempts:        0, // unlimited retries
	NonRetryableErrorTypes: []string{"ApplicationError"},
}

var options = workflow.ActivityOptions{
	// Timeout options specify when to automatically timeout Activity functions.
	StartToCloseTimeout: time.Minute,
	// Optionally provide a customized RetryPolicy.
	// Temporal retries failed Activities by default.
	RetryPolicy: retrypolicy,
}

type WorkFlow struct {
	billService BillService
}

func NewWorkFlow(billService BillService) *WorkFlow {
	return &WorkFlow{billService: billService}
}

func (w *WorkFlow) Create(ctx workflow.Context, billID string) (db.Bill, error) {
	logger.Info("creating...", billID)
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)
	activities := NewActivity(w.billService)

	var bill db.Bill
	err := workflow.ExecuteActivity(ctx, activities.Create, billID).Get(ctx, &bill)
	return bill, err
}

func (w *WorkFlow) AddBill(ctx workflow.Context, billID string, billDetail BillDetail) error {
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	// Wait for confirmation before adding to invoice
	selector := workflow.NewSelector(ctx)
	signalCh := workflow.GetSignalChannel(ctx, "confirmInvoice")

	var confirmed bool = false
	// implement selector reciever via signal channel
	selector.AddReceive(signalCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &confirmed)
	})

	// blocks untill a signal is received
	selector.Select(ctx)

	if !confirmed {
		logger.Info("confirmation denied")
		return nil
	}

	// If confirmed, add invoice
	activities := NewActivity(w.billService)
	return workflow.ExecuteActivity(ctx, activities.Add, billID, billDetail).Get(ctx, nil)
}
