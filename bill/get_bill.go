package billhandler

import (
	"context"
	"strings"

	db "encore.app/internal/db/bill"
	"encore.dev/beta/errs"
)

type GetBillRequest struct {
	Currency string `query:"currency"`
}

//encore:api public method=GET path=/bill/:billID
func (h *Handler) Get(ctx context.Context, billID string, q *GetBillRequest) (db.Bill, error) {
	logger.Info("billID: ", billID)
	if q.Currency == "" {
		q.Currency = "USD"
	}
	currency := strings.ToUpper(q.Currency)
	logger.Info("currency: ", currency)

	if _, ok := h.currencyRates[currency]; !ok {
		return db.Bill{}, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "currency not recognised",
		}
	}

	bill, err := h.billService.Get(billID)
	if err != nil {
		return bill, &errs.Error{
			Code:    errs.NotFound,
			Message: err.Error(),
		}
	}

	logger.Info(bill)

	billCopy, err := db.DeepCopy(bill)
	if err != nil {
		return bill, errs.Convert(err)
	}

	for i := 0; i < len(billCopy.Transactions); i++ {
		amount := billCopy.Transactions[i].Amount
		amount, _ = amount.Convert(currency, h.currencyRates[currency])
		billCopy.Transactions[i].Amount = amount
	}

	return billCopy, nil
}
