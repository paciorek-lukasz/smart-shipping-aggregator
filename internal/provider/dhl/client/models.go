package client

type DhlHomeApiRequest struct {
	DhlApiKey        string `json:"api_key"`
	SenderAddress    *Party `json:"sender_address"`
	RecipientAddress *Party `json:"recipient_address"`
}

type Party struct {
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Country    string `json:"count"`
}

type DhlHomeApiResponse struct {
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
	Price    int32  `json:"price"`
	Currency string `json:"currency"`
}

type DhlPickupApiRequest struct {
	DhlApiKey        string            `json:"api_key"`
	SenderAddress    *Party            `json:"sender_address"`
	RecipientAddress *Party            `json:"recipient_address"`
	LocationsLimit   int32             `json:"locations_limit"`
	SearchRadius     int32             `json:"search_radius"`
	LocationsType    []DhlLocationType `json:"locations_type"`
}

type DhlLocationType string

const (
	DHL_LOCATION_TYPE_ALL           DhlLocationType = "all"
	DHL_LOCATION_TYPE_POSTOFFICE    DhlLocationType = "postoffice"
	DHL_LOCATION_TYPE_LOCKER        DhlLocationType = "locker"
	DHL_LOCATION_TYPE_SERVICE_POINT DhlLocationType = "servicepoint"
)

type DhlPickupApiResponse struct {
	Earliest  string      `json:"earliest"`
	Latest    string      `json:"latest"`
	Price     int32       `json:"price"`
	Currency  string      `json:"currency"`
	Locations []*Location `json:"locations"`
}

type Location struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	City        string       `json:"city"`
	PostalCode  string       `json:"postal_code"`
	Country     string       `json:"country"`
	Latitude    string       `json:"latitude"`
	Longitude   string       `json:"longitude"`
	AddressLine string       `json:"address_line"`
	Type        string       `json:"type"`
	OpenTimes   []*OpenTimes `json:"open_times"`
	IsAvailable bool         `json:"is_available"`
}

const (
	DhlLocTypeLocker       = "locker"
	DhlLocTypeServicePoint = "service-point"
)

type OpenTimes struct {
	DayOfWeek string `json:"day_of_week"`
	Opens     string `json:"opens"`
	Closes    string `json:"closes"`
}
