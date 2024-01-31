package message

const (
	COMMIT_SUCCESS = 1
	COMMIT_FAIL    = 0
)

type CreatePurchaseMessage struct {
	OrderID       string   `json:"order_id" validate:"required"`
	PaymentMethod int      `json:"payment_method"`
	Amount        int      `json:"amount"`
	Vouchers      []string `json:"vouchers" validate:"required"`
}

type RollbackPurchaseMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}

type ReplyPurchaseMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
