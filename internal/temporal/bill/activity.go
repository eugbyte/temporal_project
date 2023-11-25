package temporalbill

import (
	"context"
	"time"
)

type BillService interface {
	Create(billID string) error
	Add(billID string, date time.Time, item string, amount float64) error
	Close(billID string) error
}

type BillActivity struct {
	billService BillService
}

func (a *BillActivity) Confirm(ctx context.Context, confirmed bool) bool {
	return confirmed
}

func (a *BillActivity) Create(ctx context.Context, billID string) error {
	return a.billService.Create(billID)
}

type BillDetail struct {
	Date   time.Time
	Item   string
	Amount float64
}

func (a *BillActivity) Add(ctx context.Context, billID string, billDetail BillDetail) error {
	return a.billService.Add(billID, billDetail.Date, billDetail.Item, billDetail.Amount)
}

func (a *BillActivity) Close(ctx context.Context, billID string) error {
	return a.billService.Close(billID)
}
