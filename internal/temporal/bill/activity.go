package temporalbill

import (
	"context"

	db "encore.app/internal/db/bill"
)

type BillService interface {
	Create(billID string) (db.Bill, error)
	Add(billID string, billDetail db.TransactionDetail) (db.Bill, error)
	Close(billID string) error
}

type BillActivity struct {
	billService BillService
}

func NewActivity(billService BillService) *BillActivity {
	return &BillActivity{billService: billService}
}

func (a *BillActivity) CreateBill(ctx context.Context, billID string) (db.Bill, error) {
	logger.Info("Activity: ", billID)
	return a.billService.Create(billID)
}

func (a *BillActivity) AddBill(ctx context.Context, billID string, billDetail db.TransactionDetail) (db.Bill, error) {
	return a.billService.Add(billID, billDetail)
}

func (a *BillActivity) CloseBill(ctx context.Context, billID string) error {
	return a.billService.Close(billID)
}
