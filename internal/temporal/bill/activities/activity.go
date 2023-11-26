package activities

import (
	"context"
	"fmt"

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

func CreateBillActivity(ctx context.Context, billID string) (db.Bill, error) {
	logger.Info("Activity: ", billID)
	return billService.Create(billID)
}

func IncreaseBillActivity(ctx context.Context, billID string, billDetail db.TransactionDetail) (db.Bill, error) {
	return billService.Add(billID, billDetail)
}

func CloseBillActivity(ctx context.Context, billID string) (db.Bill, error) {
	return billService.Close(billID)
}

func SanityCheckActivity(ctx context.Context) error {
	fmt.Println("Started sanity check activity")
	return nil
}
