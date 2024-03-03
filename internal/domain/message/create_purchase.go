package message

const (
	COMMIT_SUCCESS = 1
	COMMIT_FAIL    = 0
)

type CreatePurchaseMessage struct {
	CheckoutID   string   `json:"checkout_id"`
	UserID       string   `json:"user_id"`
	OrderID      string   `json:"order_id"`
	VoucherCodes []string `json:"voucher_codes"`
}

type OrderData struct {
	OrderID string `json:"order_id"`
	StoreID string `json:"store_id"`
}

type RollbackPurchaseMessage struct {
	OrderID string `json:"order_id"`
	Status  int    `json:"status"`
}

type ReplyPurchaseMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
