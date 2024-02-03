package dto

type UpdateVoucherRequest struct {
	VoucherCode string `params:"code" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

type CancelUpdateVoucherRequest struct {
	VoucherCode string `json:"voucher_code" validate:"required"`
}
