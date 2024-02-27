package dto

import "time"

type CheckoutVoucherRequest struct {
	UserID   string
	Vouchers []string `json:"vouchers"`
}

type CheckoutVoucherResponse struct {
	Items []VoucherUserDetail `json:"items"`
}

type PurchaseVoucherRequest struct {
	UserId              string               `json:"user_id"`
	OrderTotalAmount    int                  `json:"order_total_amount" validate:"required"`
	PaymentMethod       int                  `json:"payment_method" validate:"required"`
	FreeShippingVoucher *FreeShippingVoucher `json:"free_shipping_voucher"`
	PaymentVoucher      *PaymentVoucher      `json:"payment_voucher"`
	ShopVouchers        []ShopVoucher        `json:"shop_vouchers"`
}

type FreeShippingVoucher struct {
	VoucherCode string `json:"voucher_code"`
}

type PaymentVoucher struct {
	VoucherCode string `json:"voucher_code"`
}

type ShopVoucher struct {
	StoreId     string `json:"store_id"`
	SubTotal    int    `json:"sub_total"`
	VoucherCode string `json:"voucher_code"`
}

type VoucherDetail struct {
	ID               string                 `json:"id,omitempty"`
	VoucherCode      string                 `json:"voucher_code"`
	VoucherType      int                    `json:"voucher_type"`
	VoucherCounts    int                    `json:"voucher_counts"`
	Detail           string                 `json:"detail"`
	OwnerVoucher     string                 `json:"owner_voucher"`
	Status           int                    `json:"status"`
	DiscountDataResp CheckingDiscountData   `json:"discount_data"`
	VoucherRequire   CheckingVoucherRequire `json:"voucher_require,omitempty"`
	StatedTime       time.Time              `json:"stated_time,omitempty"`
	EndedTime        time.Time              `json:"ended_time,omitempty"`
	CreatedAt        time.Time              `json:"created_at,omitempty"`
	UpdatedAt        time.Time              `json:"updated_at,omitempty"`
}

type CheckingVoucherRequire struct {
	MinRequire        int64 `json:"min_require,omitempty"`
	PaymentMethod     int   `json:"payment_method,omitempty"`
	MaxVoucherPerUser int   `json:"max_voucher_per_user,omitempty"`
}

type CheckingDiscountData struct {
	DiscountType    int     `json:"discount_type,omitempty"`
	ShippingValue   uint    `json:"shipping_value,omitempty"`
	DiscountValue   uint    `json:"discount_value,omitempty"`
	DiscountPercent float32 `json:"discount_percent,omitempty"`
	MaximumValue    uint    `json:"maximum_value,omitempty"`
}

type UseVoucherResponse struct {
	StoreVouchers   []VoucherDetail `json:"store_vouchers"`
	ShippingVoucher VoucherDetail   `json:"shipping_voucher"`
	PaymentVoucher  VoucherDetail   `json:"payment_voucher"`
}
