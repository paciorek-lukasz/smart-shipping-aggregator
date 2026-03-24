package dhlclient

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
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	AddressLine string     `json:"address_line"`
	Type        string     `json:"type"`
	OpenTimes   *OpenTimes `json:"open_times"`
	IsAvailable bool       `json:"is_available"`
}

type OpenTimes struct {
	Monday    string `json:"Monday"`
	Tuesday   string `json:"Tuesday"`
	Wednesday string `json:"Wednesday"`
	Thursday  string `json:"Thursday"`
	Friday    string `json:"Friday"`
	Saturday  string `json:"Saturday"`
	Sunday    string `json:"Sunday"`
}
