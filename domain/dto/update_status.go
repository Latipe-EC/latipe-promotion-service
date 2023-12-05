package dto

type UpdateVoucherRequest struct {
	VoucherCode string `params:"code" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}
