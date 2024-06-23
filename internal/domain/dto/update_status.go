package dto

const INVALID_VOUCHER_STATUS = -2

type UpdateVoucherRequest struct {
	VoucherCode string `params:"code" validate:"required"`
	Status      int    `json:"status" validate:"required"`
	StoreID     string `json:"store_id"`
}

func (rq UpdateVoucherRequest) GetStatus() int {
	if rq.Status != -1 && rq.Status != 0 && rq.Status != 1 {
		return INVALID_VOUCHER_STATUS
	}
	return rq.Status
}

type CancelUpdateVoucherRequest struct {
	VoucherCode string `json:"voucher_code" validate:"required"`
}
