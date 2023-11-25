package temporalbill

import (
	"time"

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

func (w *WorkFlow) Create(ctx workflow.Context, billID string) error {
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)
	return workflow.ExecuteActivity(ctx, w.billService.Create, billID).Get(ctx, nil)
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
	return workflow.ExecuteActivity(ctx, w.billService.Add, billID, billDetail).Get(ctx, nil)
}
