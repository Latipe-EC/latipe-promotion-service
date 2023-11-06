package dto

type UseVoucherRequest struct {
	Vouchers []string `json:"vouchers" validate:"required"`
}

type UseVoucherResponse struct {
	IsSuccess bool                `json:"is_success"`
	Items     []VoucherRespDetail `json:"items"`
}
