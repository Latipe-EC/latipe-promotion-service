package dto

import "latipe-promotion-services/pkgs/pagable"

type GetUserVoucherRequest struct {
	Query *pagable.Query
}
type GetUserVoucherResponse struct {
	pagable.ListResponse
}
