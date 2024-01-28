package dto

import "latipe-promotion-services/pkgs/pagable"

type GetVoucherListRequest struct {
	Query *pagable.Query
}
type GetVoucherListResponse struct {
	pagable.ListResponse
}
