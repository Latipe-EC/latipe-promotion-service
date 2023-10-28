package dto

type CreateVoucherRequest struct {
	VoucherCode     string         `json:"voucher_code" validate:"required"`
	VoucherType     int            `json:"voucher_type" validate:"required"`
	VoucherCounts   int            `json:"voucher_counts" validate:"required"`
	Detail          string         `json:"detail" validate:"required"`
	DiscountPercent float64        `json:"discount_percent" validate:"required"`
	DiscountValue   int            `json:"discount_value" validate:"required"`
	VoucherRequire  VoucherRequire `json:"voucher_require" validate:"required"`
	StatedTime      string         `json:"stated_time" validate:"required"`
	EndedTime       string         `json:"ended_time" validate:"required"`
	OwnerVoucherId  string
}

type VoucherRequire struct {
	MinRequire    int64 `json:"min_require" validate:"required"`
	MemberType    int   `json:"member_type,omitempty"`
	PaymentMethod int   `json:"payment_method,omitempty"`
}

type CreateVoucherResponse struct {
	VoucherRespDetail
}
