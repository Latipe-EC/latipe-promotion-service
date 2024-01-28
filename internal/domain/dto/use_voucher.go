package dto

type UseVoucherRequest struct {
	OrderID  string   `json:"order_id" validate:"required"`
	Vouchers []string `json:"vouchers" validate:"required"`
}

type UseVoucherResponse struct {
	IsSuccess bool                `json:"is_success"`
	Items     []VoucherRespDetail `json:"items"`
}

type CheckingVoucherRequest struct {
	Vouchers         []string `json:"vouchers" validate:"required"`
	OrderTotalAmount int      `json:"order_total_amount" validate:"required"`
	PaymentMethod    int      `json:"payment_method" validate:"required"`
	UserId           string   `json:"user_id"`
}
