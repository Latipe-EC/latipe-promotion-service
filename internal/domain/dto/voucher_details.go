package dto

import (
	"time"
)

type VoucherRespDetail struct {
	ID               string             `json:"id,omitempty"`
	VoucherCode      string             `json:"voucher_code"`
	VoucherType      int                `json:"voucher_type"`
	VoucherCounts    int                `json:"voucher_counts"`
	Detail           string             `json:"detail"`
	OwnerVoucher     string             `json:"owner_voucher"`
	Status           int                `json:"status"`
	DiscountDataResp DiscountDataResp   `json:"discount_data"`
	VoucherRequire   VoucherRequireResp `json:"voucher_require,omitempty"`
	StatedTime       time.Time          `json:"stated_time,omitempty"`
	EndedTime        time.Time          `json:"ended_time,omitempty"`
	CreatedAt        time.Time          `json:"created_at,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty"`
}

type VoucherRequireResp struct {
	MinRequire        int64 `json:"min_require,omitempty"`
	PaymentMethod     int   `json:"payment_method,omitempty"`
	MaxVoucherPerUser int   `json:"max_voucher_per_user,omitempty"`
}

type DiscountDataResp struct {
	DiscountType    int     `json:"discount_type,omitempty"`
	ShippingValue   uint    `json:"shipping_value,omitempty"`
	DiscountValue   uint    `json:"discount_value,omitempty"`
	DiscountPercent float32 `json:"discount_percent,omitempty"`
	MaximumValue    uint    `json:"maximum_value,omitempty"`
}

type VoucherUserDetail struct {
	ID               string             `json:"id,omitempty"`
	VoucherCode      string             `json:"voucher_code"`
	VoucherType      int                `json:"voucher_type"`
	VoucherCounts    int                `json:"voucher_counts,omitempty"`
	Detail           string             `json:"detail,omitempty"`
	OwnerVoucher     string             `json:"owner_voucher,omitempty"`
	DiscountDataResp DiscountDataResp   `json:"discount_data"`
	VoucherRequire   VoucherRequireResp `json:"voucher_require,omitempty"`
	StatedTime       time.Time          `json:"stated_time,omitempty"`
	EndedTime        time.Time          `json:"ended_time,omitempty"`
	Usable           bool               `json:"usable"`
}
