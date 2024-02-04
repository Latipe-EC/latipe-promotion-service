package dto

type RollbackVoucherRequest struct {
	VoucherCodes []string `json:"voucher_codes"`
}

type RollbackVoucherResponse struct {
	Success bool `json:"is_success"`
}
