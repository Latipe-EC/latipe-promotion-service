package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	FREE_SHIP        = 1
	PAYMENT_DISCOUNT = 2
	STORE_DISCOUNT   = 3

	PENDING  = 0
	ACTIVE   = 1
	INACTIVE = -1

	VOUCHER_APPLY_SUCCESS = 1
	VOUCHER_APPLY_FAILED  = -1

	FIXED_DISCOUNT   = 0
	PERCENT_DISCOUNT = 1

	COD_METHOD = 1
)

type Voucher struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VoucherCode    string             `json:"voucher_code" bson:"voucher_code"`
	VoucherType    int                `json:"voucher_type" bson:"voucher_type"`
	VoucherCounts  int                `json:"voucher_counts" bson:"voucher_counts"`
	Detail         string             `json:"detail"  bson:"detail,omitempty"`
	OwnerVoucher   string             `json:"owner_voucher" bson:"owner_voucher"`
	Status         int                `json:"status" bson:"status"`
	DiscountData   *DiscountData      `json:"discount_data" bson:"discount_data"`
	VoucherRequire *VoucherRequire    `json:"voucher_require" bson:"voucher_require"`
	StatedTime     time.Time          `json:"stated_time" bson:"stated_time,omitempty"`
	EndedTime      time.Time          `json:"ended_time" bson:"ended_time,omitempty"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
}

type VoucherRequire struct {
	MinRequire        int64 `json:"min_require" bson:"min_require"`
	PaymentMethod     int   `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	MaxVoucherPerUser int   `json:"max_voucher_per_user" bson:"max_voucher_per_user"`
}

type DiscountData struct {
	DiscountType    int     `json:"discount_type,omitempty" bson:"discount_type,omitempty"`
	ShippingValue   uint    `json:"shipping_value,omitempty" bson:"shipping_value,omitempty"`
	DiscountValue   uint    `json:"discount_value,omitempty" bson:"discount_value,omitempty"`
	DiscountPercent float32 `json:"discount_percent,omitempty" bson:"discount_percent,omitempty"`
	MaximumValue    uint    `json:"maximum_value,omitempty" bson:"maximum_value,omitempty"`
}

type VoucherUsingLog struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	VoucherCode      string             `json:"voucher_code" bson:"voucher_code"`
	VoucherID        primitive.ObjectID `json:"voucher_type" bson:"voucher_type"`
	UserID           string             `json:"user_id" bson:"user_id"`
	CheckoutPurchase CheckoutPurchase   `json:"checkout_purchase" bson:"checkout_purchase"`
	Status           int                `json:"status" bson:"status"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
}

type CheckoutPurchase struct {
	CheckoutID string   `json:"checkout_id"  bson:"checkout_id"`
	OrderIDs   []string `json:"order_ids" bson:"order_ids"`
}
