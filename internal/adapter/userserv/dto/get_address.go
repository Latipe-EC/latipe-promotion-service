package dto

const url = "/api/v1/users/address"

type GetAddressRequest struct {
	AddressId string `json:"addressId"`
}

type GetAddressResponse struct {
	DetailAddress       string `json:"detailAddress"`
	City                string `json:"city"`
	ZipCode             string `json:"zipCode"`
	DistrictId          int    `json:"districtId"`
	DistrictName        string `json:"districtName"`
	StateOrProvinceId   int    `json:"stateOrProvinceId"`
	StateOrProvinceName string `json:"stateOrProvinceName"`
	CountryId           int    `json:"countryId"`
	CountryName         string `json:"countryName"`
}

func (GetAddressRequest) URL() string {
	return url
}
