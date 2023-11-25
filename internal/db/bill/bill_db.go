package billdb

import (
	"fmt"
	"sync"
	"time"

	customerrors "encore.app/internal/custom_errors"
)

type BillDB struct {
	mu    sync.Mutex
	Bills map[string]Bill
}

type Status string

const (
	OPEN   Status = "OPEN"
	CLOSED Status = "CLOSED"
)

type Bill struct {
	ID           string                      `json:"ID"`
	Status       Status                      `json:"status"`
	Transactions map[int64]TransactionDetail `json:"transactions"` // Unix timestamp against $amount
}

type TransactionDetail struct {
	ItemName string
	Amount   float64
}

func New() *BillDB {
	b := BillDB{}
	b.Bills = make(map[string]Bill)
	return &b
}

func (b *BillDB) Create(billID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; ok {
		return customerrors.NewAppError(fmt.Sprintf("%s already exist", billID))
	}

	b.Bills[billID] = Bill{
		ID:           billID,
		Status:       OPEN,
		Transactions: make(map[int64]TransactionDetail),
	}

	return nil
}

func (b *BillDB) Add(billID string, date time.Time, item string, amount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; !ok {
		return customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	bill := b.Bills[billID]

	if bill.Status == CLOSED {
		return customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	timestamp := date.Unix()
	bill.Transactions[timestamp] = TransactionDetail{
		ItemName: item,
		Amount:   amount,
	}

	b.Bills[billID] = bill
	return nil
}

func (b *BillDB) Close(billID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; !ok {
		return customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	bill := b.Bills[billID]
	bill.Status = CLOSED

	return nil
}

func (b *BillDB) Get(billID string) (Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; !ok {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}
	return b.Bills[billID], nil
}
