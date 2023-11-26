package billdb

import (
	"encoding/json"
	"fmt"
	"sync"

	debug "encore.app/internal/logger"
	"github.com/bojanz/currency"

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
	Timestamp int64  `json:"timestamp"`
	ItemName  string `json:"itemName"`
	// stored as USD
	Amount currency.Amount `json:"amount"`
}

// Singleton instance to be used by both the Handlers and Temporal workers.
var BillService = New()

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

	_currency := detail.Amount.CurrencyCode()
	if _currency != "USD" {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("currency must be USD, got %s instead", _currency))
	}

	bill.Transactions = append(bill.Transactions, detail)

	b.Bills[billID] = bill
	return b.Bills[billID], nil
}

func (b *BillDB) Close(billID string) (Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.Bills[billID]; !ok {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}

	bill := b.Bills[billID]

	if bill.Status == CLOSED {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s already closed", billID))
	}
	bill.Status = CLOSED

	b.Bills[billID] = bill
	return bill, nil
}

func (b *BillDB) Get(billID string) (Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	logger.Info(b.Bills)

	if _, ok := b.Bills[billID]; !ok {
		return Bill{}, customerrors.NewAppError(fmt.Sprintf("%s does not exist", billID))
	}
	return b.Bills[billID], nil
}

func DeepCopy(bill Bill) (Bill, error) {
	var billCopy Bill = Bill{}
	byts, err := json.Marshal(bill)
	if err != nil {
		return billCopy, err
	}
	json.Unmarshal(byts, &billCopy)
	return billCopy, err

}
