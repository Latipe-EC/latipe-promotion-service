package dto

type CreateVoucherRequest struct {
	VoucherCode     string         `json:"voucher_code" validate:"required"`
	VoucherType     int            `json:"voucher_type" validate:"required"`
	VoucherCounts   int            `json:"voucher_counts" validate:"required"`
	Detail          string         `json:"detail" validate:"required"`
	DiscountPercent float64        `json:"discount_percent"`
	DiscountValue   int            `json:"discount_value" validate:"required"`
	VoucherRequire  VoucherRequire `json:"voucher_require,omitempty"`
	StatedTime      string         `json:"stated_time" validate:"required"`
	EndedTime       string         `json:"ended_time" validate:"required"`
	OwnerVoucherId  string
}

type VoucherRequire struct {
	MinRequire    int64 `json:"min_require"`
	PaymentMethod int   `json:"payment_method"`
}

type CreateVoucherResponse struct {
	VoucherRespDetail
}
