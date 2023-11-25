package billdb

import (
	"fmt"
	"sync"

	debug "encore.app/internal/logger"

	customerrors "encore.app/internal/custom_errors"
)

var logger = debug.Logger

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
	ID           string              `json:"ID"`
	Status       Status              `json:"status"`
	Transactions []TransactionDetail `json:"transactions"` // Unix timestamp against $amount
}

type TransactionDetail struct {
	Timestamp int64   `json:"timestamp"`
	ItemName  string  `json:"itemName"`
	Amount    float64 `json:"amount"`
}

func New() *BillDB {
	b := BillDB{}
	b.Bills = make(map[string]Bill)
	return &b
}

func (b *BillDB) Create(billID string) (Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if bill, ok := b.Bills[billID]; ok {
		logger.Info("bill: ", bill)
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s already exist", billID))
	}

	b.Bills[billID] = Bill{
		ID:           billID,
		Status:       OPEN,
		Transactions: make([]TransactionDetail, 0),
	}

	return b.Bills[billID], nil
}

func (b *BillDB) Add(billID string, detail TransactionDetail) (Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; !ok {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	bill := b.Bills[billID]

	if bill.Status == CLOSED {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	bill.Transactions = append(bill.Transactions, detail)

	b.Bills[billID] = bill
	return b.Bills[billID], nil
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
