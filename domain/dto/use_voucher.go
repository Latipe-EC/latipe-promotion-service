package dto

type UseVoucherRequest struct {
	OrderUUID string   `json:"order_uuid" validate:"required"`
	Vouchers  []string `json:"vouchers" validate:"required"`
}

type UseVoucherResponse struct {
	IsSuccess bool                `json:"is_success"`
	Items     []VoucherRespDetail `json:"items"`
}

type CheckingVoucherRequest struct {
	Vouchers         []string `json:"vouchers" validate:"required"`
	OrderTotalAmount int      `json:"order_total_amount" validate:"required"`
	PaymentMethod    int      `json:"payment_method" validate:"required"`
	UserId           string   `json:"user_id" validate:"required"`
}
