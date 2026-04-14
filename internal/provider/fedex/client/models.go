package client

type FedexHomeApiRequest struct {
	FedexClientId      string            `json:"fedex_client_id"`
	SourceDetails      *ShipperConsignee `json:"source_details"`
	DestinationDetails *ShipperConsignee `json:"destination_details"`
}

type ShipperConsignee struct {
	Street1     string `json:"street1"`
	Street2     string `json:"street2"`
	PostalCode  string `json:"postal_code"`
	City        string `json:"city"`
	CountryCode string `json:"country_code"`
}

type FedexHomeApiResponse struct {
	EstimatedPickup   string `json:"estimated_pickup"`
	EstimatedDelivery string `json:"estimated_delivery"`
	NetAmount         int32  `json:"net_amount"`
	CurrencyCode      string `json:"currency_code"`
}

type FedexPickupApiRequest struct {
	FedexClientId      string             `json:"fedex_client_id"`
	SourceDetails      *ShipperConsignee  `json:"source_details"`
	DestinationDetails *ShipperConsignee  `json:"destination_details"`
	DroppointsLimit    int32              `json:"droppoints_limit"`
	WithinRadiusMiles  int32              `json:"within_radius_miles"`
	ServiceTypeFilter  []FedexServiceType `json:"service_type_filter"`
}

type FedexServiceType string

const (
	FEDEX_SERVICE_ALL        FedexServiceType = "all"
	FEDEX_SERVICE_EXPRESS    FedexServiceType = "express"
	FEDEX_SERVICE_GROUND     FedexServiceType = "ground"
	FEDEX_SERVICE_SMART_POST FedexServiceType = "smartpost"
)

type FedexPickupApiResponse struct {
	EstimatedPickup   string            `json:"estimated_pickup"`
	EstimatedDelivery string            `json:"estimated_delivery"`
	NetAmount         int32             `json:"net_amount"`
	CurrencyCode      string            `json:"currency_code"`
	Droppoints        []*FedexDropPoint `json:"droppoints"`
}

type FedexDropPoint struct {
	LocationId      string             `json:"location_id"`
	LocationName    string             `json:"location_name"`
	City            string             `json:"city"`
	StateOrProvince string             `json:"state_or_province"`
	PostalCode      string             `json:"postal_code"`
	CountryCode     string             `json:"country_code"`
	LatLong         string             `json:"lat_long"`
	Address         string             `json:"address"`
	LocationType    string             `json:"location_type"`
	HoursOperation  []*OperationWindow `json:"hours_operation"`
	CurrentlyOpen   bool               `json:"currently_open"`
}

const (
	FedexLocTypeStation    = "station"
	FedexLocTypeauthorized = "authorized"
)

type OperationWindow struct {
	DaysOfOperation string `json:"days_of_operation"`
	WindowOpen      string `json:"window_open"`
	WindowClose     string `json:"window_close"`
}
