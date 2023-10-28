package dto

import (
	"time"
)

type VoucherRespDetail struct {
	ID              string         `json:"_id,omitempty"`
	VoucherCode     string         `json:"voucher_code"`
	VoucherType     int            `json:"voucher_type"`
	VoucherCounts   int            `json:"voucher_counts,omitempty"`
	Detail          string         `json:"detail,omitempty"`
	OwnerVoucher    string         `json:"owner_voucher,omitempty"`
	DiscountPercent float64        `json:"discount_percent,omitempty"`
	DiscountValue   int            `json:"discount_value,omitempty"`
	VoucherRequire  VoucherReqResp `json:"voucher_require,omitempty"`
	CreateAt        time.Time      `json:"create_at,omitempty"`
	UpdateAt        time.Time      `json:"update_at,omitempty"`
	StatedTime      time.Time      `json:"stated_time,omitempty"`
	EndedTime       time.Time      `json:"ended_time,omitempty"`
	Status          int            `json:"status"`
}

type VoucherReqResp struct {
	MinRequire          int64  `json:"min_require"`
	MemberType          int    `json:"member_type,omitempty"`
	PaymentMethod       int    `json:"payment_method,omitempty"`
	RequiredOwnerProdId string `json:"required_owner_prod_id,omitempty"`
}
