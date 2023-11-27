package activity

import (
	"context"
	"fmt"

	customerrors "encore.app/internal/custom_error"
	debug "encore.app/internal/logger"

	db "encore.app/internal/db/bill"
)

var logger = debug.Logger

type BillService interface {
	Create(billID string) (db.Bill, error)
	Add(billID string, billDetail db.TransactionDetail) (db.Bill, error)
	Close(billID string) (db.Bill, error)
}

var billService BillService

func init() {
	// alternatively, create a struct with { billService BillService } for dependency injection,
	// but there seems to be conflicting practices regarding whether DI is best practice (https://github.com/temporalio/sdk-java/issues/745).
	billService = db.BillService
}

// Note that Activities must be named differently from Workflows, otherwise the test mocking fails.

func CreateBillAct(ctx context.Context, billID string) (db.Bill, error) {
	logger.Info("Activity: ", billID)
	bill, err := billService.Create(billID)
	if err != nil {
		// stop Temporal from retrying, as either bill ID already exists or bill has been closed
		err = customerrors.NewNonRetryError(err.Error())
	}
	return bill, err
}

func IncreaseBillAct(ctx context.Context, billID string, billDetail db.TransactionDetail) (db.Bill, error) {
	bill, err := billService.Add(billID, billDetail)
	if err != nil {
		// stop Temporal from retrying, as bill has been been closed
		err = customerrors.NewNonRetryError(err.Error())
	}
	return bill, err
}

func CloseBillAct(ctx context.Context, billID string) (db.Bill, error) {
	return billService.Close(billID)
}

func SanityCheckAct(ctx context.Context) error {
	fmt.Println("Started sanity check activity")
	return nil
}
