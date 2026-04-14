package client

type GlsHomeApiRequest struct {
	GlsAuthToken  string       `json:"gls_auth_token"`
	ShipperData   *ContactData `json:"shipper_data"`
	ConsigneeData *ContactData `json:"consignee_data"`
}

type ContactData struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	PostalCode   string `json:"postal_code"`
	TownCity     string `json:"town_city"`
	IsoCountry   string `json:"iso_country"`
}

type GlsHomeApiResponse struct {
	WindowFrom string `json:"window_from"`
	WindowTo   string `json:"window_to"`
	RateCents  int32  `json:"rate_cents"`
	Currency   string `json:"currency"`
}

type GlsPickupApiRequest struct {
	GlsAuthToken      string            `json:"gls_auth_token"`
	ShipperData       *ContactData      `json:"shipper_data"`
	ConsigneeData     *ContactData      `json:"consignee_data"`
	DepotsLimit       int32             `json:"depots_limit"`
	SearchAreaKm      int32             `json:"search_area_km"`
	DepotFacilityType []GlsFacilityType `json:"depot_facility_type"`
}

type GlsFacilityType string

const (
	GLS_FACILITY_ALL        GlsFacilityType = "all"
	GLS_FACILITY_PARCELSHOP GlsFacilityType = "parcelshop"
	GLS_FACILITY_DEPOT      GlsFacilityType = "depot"
	GLS_FACILITY_PARTNER    GlsFacilityType = "partner_shop"
)

type GlsPickupApiResponse struct {
	WindowFrom string       `json:"window_from"`
	WindowTo   string       `json:"window_to"`
	RateCents  int32        `json:"rate_cents"`
	Currency   string       `json:"currency"`
	Depots     []*DepotData `json:"depots"`
}

type DepotData struct {
	DepotCode   string          `json:"depot_code"`
	DepotName   string          `json:"depot_name"`
	City        string          `json:"city"`
	PostalCode  string          `json:"postal_code"`
	IsoCountry  string          `json:"iso_country"`
	GeoLat      string          `json:"geo_lat"`
	GeoLon      string          `json:"geo_lon"`
	Street      string          `json:"street"`
	Facility    string          `json:"facility"`
	Timetable   []*TimetableDay `json:"timetable"`
	Operational bool            `json:"operational"`
}

const (
	GlsFacilityShop  = "parcelshop"
	GlsFacilityDepot = "depot"
)

type TimetableDay struct {
	DayName string `json:"day_name"`
	OpenHr  string `json:"open_hr"`
	CloseHr string `json:"close_hr"`
}
