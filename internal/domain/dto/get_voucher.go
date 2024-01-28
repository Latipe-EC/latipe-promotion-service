package dto

type GetVoucherByIdRequest struct {
	Id string `params:"id"`
}

type GetVoucherByCodeRequest struct {
	Code string `params:"code"`
}
