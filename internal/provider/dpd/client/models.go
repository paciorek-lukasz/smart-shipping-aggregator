package client

type DpdHomeApiRequest struct {
	DpdApiKey        string    `json:"api_key"`
	SenderAddress    *DpdParty `json:"sender_address"`
	RecipientAddress *DpdParty `json:"recipient_address"`
}

type DpdParty struct {
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

type DpdHomeApiResponse struct {
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
	Price    int32  `json:"price"`
	Currency string `json:"currency"`
}

type DpdPickupApiRequest struct {
	DpdApiKey        string            `json:"api_key"`
	SenderAddress    *DpdParty         `json:"sender_address"`
	RecipientAddress *DpdParty         `json:"recipient_address"`
	LocationsLimit   int32             `json:"locations_limit"`
	SearchRadius     int32             `json:"search_radius"`
	LocationsType    []DpdLocationType `json:"locations_type"`
}

type DpdLocationType string

const (
	DPD_LOCATION_TYPE_ALL          DpdLocationType = "all"
	DPD_LOCATION_TYPE_PACKAGE_SHOP DpdLocationType = "packageshop"
	DPD_LOCATION_TYPE_LOCKER       DpdLocationType = "locker"
	DPD_LOCATION_TYPE_PARCEL_BOX   DpdLocationType = "parcelbox"
)

type DpdPickupApiResponse struct {
	Earliest  string         `json:"earliest"`
	Latest    string         `json:"latest"`
	Price     int32          `json:"price"`
	Currency  string         `json:"currency"`
	Locations []*DpdLocation `json:"locations"`
}

type DpdLocation struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	City        string          `json:"city"`
	PostalCode  string          `json:"postal_code"`
	Country     string          `json:"country"`
	Latitude    string          `json:"latitude"`
	Longitude   string          `json:"longitude"`
	AddressLine string          `json:"address_line"`
	Type        string          `json:"type"`
	OpenTimes   []*DpdOpenTimes `json:"open_times"`
	IsAvailable bool            `json:"is_available"`
}

const (
	DpdLocTypeLocker      = "locker"
	DpdLocTypePackageShop = "packageshop"
)

type DpdOpenTimes struct {
	DayOfWeek string `json:"day_of_week"`
	Opens     string `json:"opens"`
	Closes    string `json:"closes"`
}
