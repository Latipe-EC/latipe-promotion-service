package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	FREE_SHIP      = 1
	DISCOUNT_ORDER = 2

	PENDING   = 0
	ACTIVE    = 1
	IN_ACTIVE = 2
)

type Voucher struct {
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VoucherCode     string             `json:"voucher_code" bson:"voucher_code,omitempty"`
	VoucherType     int                `json:"voucher_type" bson:"voucher_type,omitempty"`
	VoucherCounts   int                `json:"voucher_counts" bson:"voucher_counts,omitempty"`
	Detail          string             `json:"detail"  bson:"detail,omitempty"`
	OwnerVoucher    string             `json:"owner_voucher" bson:"owner_voucher"`
	DiscountPercent float64            `json:"discount_percent" bson:"discount_percent,omitempty"`
	DiscountValue   int                `json:"discount_value" bson:"discount_value,omitempty"`
	VoucherRequire  VoucherRequire     `json:"voucher_require" bson:"voucher_require"`
	CreateAt        time.Time          `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt        time.Time          `json:"update_at" bson:"update_at,omitempty"`
	StatedTime      time.Time          `json:"stated_time" bson:"stated_time,omitempty"`
	EndedTime       time.Time          `json:"ended_time" bson:"ended_time,omitempty"`
	Status          int                `json:"status" bson:"is_active,omitempty"`
}

type VoucherUsedLog struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VoucherCode string             `json:"delivery_name" bson:"delivery_name,omitempty"`
	CreateAt    time.Time          `json:"create_at" bson:"create_at,omitempty"`
	IsSuccess   bool               `json:"is_success" bson:"is_success,omitempty"`
}

type VoucherRequire struct {
	MinRequire          int64  `json:"min_require" bson:"min_require"`
	MemberType          int    `json:"member_type,omitempty" bson:"member_type,omitempty"`
	PaymentMethod       int    `json:"payment_method,omitempty" bson:"payment_method,omitempty"`
	RequiredOwnerProdId string `json:"required_owner_prod_id,omitempty" bson:"required_owner_prod_id,omitempty"`
}
