package dto

const getStore = "/api/v1/stores/validate-store/"

type GetStoreIdByUserRequest struct {
	BaseHeader
	UserID string `json:"user_id"`
}

type GetStoreIdByUserResponse struct {
	StoreID string `json:"store_id"`
}

func (GetStoreIdByUserRequest) URL() string {
	return getStore
}
