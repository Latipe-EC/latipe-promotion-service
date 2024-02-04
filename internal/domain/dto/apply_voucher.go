package dto

type ApplyVoucherRequest struct {
	UserId       string          `json:"user_id"`
	CheckoutData CheckoutRequest `json:"checkout_data"`
	Vouchers     []string        `json:"vouchers"`
}

type CheckoutRequest struct {
	CheckoutID string   `json:"checkout_id"`
	OrderIds   []string `json:"order_ids"`
}
