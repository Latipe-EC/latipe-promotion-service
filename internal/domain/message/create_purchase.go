package message

const (
	COMMIT_SUCCESS = 1
	COMMIT_FAIL    = 0
)

type CreatePurchaseMessage struct {
	UserId       string          `json:"user_id"`
	CheckoutData CheckoutRequest `json:"checkout_data"`
	Vouchers     []string        `json:"vouchers"`
}

type CheckoutRequest struct {
	CheckoutID string      `json:"checkout_id"`
	OrderData  []OrderData `json:"order_data"`
}

type OrderData struct {
	OrderID string `json:"order_id"`
	StoreID string `json:"store_id"`
}

type RollbackPurchaseMessage struct {
	VoucherCodes []string `json:"voucher_codes"`
	OrderID      string   `json:"order_id"`
}

type ReplyPurchaseMessage struct {
	Status  int    `json:"status"`
	OrderID string `json:"order_id"`
}
