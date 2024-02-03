package dto

type CreateVoucherRequest struct {
	VoucherCode    string         `json:"voucher_code" validate:"required"`
	VoucherType    int            `json:"voucher_type" validate:"required"`
	VoucherCounts  int            `json:"voucher_counts" validate:"required"`
	Detail         string         `json:"detail" validate:"required"`
	DiscountData   DiscountData   `json:"discount_data" validate:"required"`
	VoucherRequire VoucherRequire `json:"voucher_require" validate:"required"`
	StatedTime     string         `json:"stated_time" validate:"required"`
	EndedTime      string         `json:"ended_time" validate:"required"`
	StoreID        string
}

type VoucherRequire struct {
	MinRequire        int64 `json:"min_require" bson:"min_require"`
	PaymentMethod     int   `json:"payment_method,omitempty"`
	MaxVoucherPerUser int   `json:"max_voucher_per_user" validate:"required"`
}

type DiscountData struct {
	DiscountType    int     `json:"discount_type"`
	ShippingValue   uint    `json:"shipping_value"`
	DiscountValue   uint    `json:"discount_value"`
	DiscountPercent float32 `json:"discount_percent"`
	MaximumValue    uint    `json:"maximum_value"`
}

type CreateVoucherResponse struct {
	VoucherRespDetail
}
