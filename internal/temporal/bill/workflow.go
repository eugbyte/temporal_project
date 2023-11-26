package temporalbill

import (
	db "encore.app/internal/db/bill"
	debug "encore.app/internal/logger"
	"go.temporal.io/sdk/workflow"
)

var logger = debug.Logger

const SignalChannel = "confirm-invoice"

type BillService interface {
	Create(billID string) (db.Bill, error)
	Add(billID string, billDetail db.TransactionDetail) (db.Bill, error)
	Close(billID string) (db.Bill, error)
}

type WorkFlow struct {
	billService BillService
}

func NewWorkFlows(billService BillService) *WorkFlow {
	return &WorkFlow{billService: billService}
}

func (w *WorkFlow) CreateBill(ctx workflow.Context, billID string) (db.Bill, error) {
	logger.Info("creating...", billID)
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)
	activities := NewActivities(w.billService)

	var bill db.Bill
	err := workflow.ExecuteActivity(ctx, activities.CreateBill, billID).Get(ctx, &bill)
	return bill, err
}

func (w *WorkFlow) IncreaseBill(ctx workflow.Context, billID string, billDetail db.TransactionDetail) error {
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
	selector.Select(ctx)

	if !confirmed {
		logger.Info("confirmation denied")
		return nil
	}

	// If confirmed, add invoice
	activities := NewActivities(w.billService)
	return workflow.ExecuteActivity(ctx, activities.IncreaseBill, billID, billDetail).Get(ctx, nil)
}

func (w *WorkFlow) CloseBill(ctx workflow.Context, billID string) (db.Bill, error) {
	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	activities := NewActivities(w.billService)

	var bill db.Bill
	err := workflow.ExecuteActivity(ctx, activities.CloseBill, billID).Get(ctx, &bill)
	return bill, err
}
