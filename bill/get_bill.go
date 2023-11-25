package billhandler

import (
	"context"

	db "encore.app/internal/db/bill"
)

//encore:api public method=GET path=/bill/:billID
func (h *Handler) Get(ctx context.Context, billID string) (db.Bill, error) {
	return h.billService.Get(billID)
}
