package client

type UpsHomeApiRequest struct {
	UpsApiKey     string       `json:"api_key"`
	SenderAddress *AddressInfo `json:"sender_address"`
	RecipientInfo *AddressInfo `json:"recipient_info"`
}

type AddressInfo struct {
	Street      string `json:"street"`
	PostalCode  string `json:"postal_code"`
	CityName    string `json:"city_name"`
	CountryCode string `json:"country_code"`
}

type UpsHomeApiResponse struct {
	AvailableFrom string `json:"available_from"`
	AvailableTo   string `json:"available_to"`
	FeeAmount     int32  `json:"fee_amount"`
	FeeCurrency   string `json:"fee_currency"`
}

type UpsPickupApiRequest struct {
	UpsApiKey       string        `json:"api_key"`
	SenderAddress   *AddressInfo  `json:"sender_address"`
	RecipientInfo   *AddressInfo  `json:"recipient_info"`
	DropPointsLimit int32         `json:"drop_points_limit"`
	SearchRadius    int32         `json:"search_radius"`
	DropPointType   []UpsDropType `json:"drop_point_type"`
}

type UpsDropType string

const (
	UPS_DROP_TYPE_ALL          UpsDropType = "all"
	UPS_DROP_TYPE_PACKAGE_SHOP UpsDropType = "package_shop"
	UPS_DROP_TYPE_DROP_BOX     UpsDropType = "drop_box"
	UPS_DROP_TYPE_AUTHORIZED   UpsDropType = "authorized_partner"
)

type UpsPickupApiResponse struct {
	AvailableFrom string       `json:"available_from"`
	AvailableTo   string       `json:"available_to"`
	FeeAmount     int32        `json:"fee_amount"`
	FeeCurrency   string       `json:"fee_currency"`
	DropPoints    []*DropPoint `json:"drop_points"`
}

type DropPoint struct {
	DropPointId    string      `json:"drop_point_id"`
	DisplayName    string      `json:"display_name"`
	City           string      `json:"city"`
	ZipCode        string      `json:"zip_code"`
	Country        string      `json:"country"`
	CoordLat       string      `json:"coord_lat"`
	CoordLng       string      `json:"coord_lng"`
	FullAddress    string      `json:"full_address"`
	Category       string      `json:"category"`
	OperationHours []*Schedule `json:"operation_hours"`
	IsActive       bool        `json:"is_active"`
}

const (
	UpsCategoryDropBox     = "drop-box"
	UpsCategoryPackageShop = "package-shop"
)

type Schedule struct {
	Weekday   string `json:"weekday"`
	StartHour string `json:"start_hour"`
	EndHour   string `json:"end_hour"`
}
