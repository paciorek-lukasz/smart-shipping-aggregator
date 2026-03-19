package domain

import "time"

type GetQuotesRequest struct {
	Sender    *Party
	Recipient *Party
	Package   *Package
}

type Party struct {
	Name    string
	Address *Address
	Phone   string
	Email   string
}

type Address struct {
	Address    string
	PostalCode string
	City       string
	Country    string
	Longitude  string
	Latitude   string
}

type Package struct {
	Items []*Item
	// Price is provided in cents
	TotalPrice int32
	Currency   string
	Dimensions *Dimensions
}

type Item struct {
	ItemID int32
	Sku    string
	Name   string
	// Price is provided in cents
	Price    int32
	Quantity int32
}

type Dimensions struct {
	LengthCm       int32
	WidthCm        int32
	HeightCm       int32
	TotalWeightG   int32
	TotalVolumeCm3 int32
}

type GetQuotesResponse struct {
	Options []*Option
}

type Option struct {
	OptionId       int32
	CarrierProduct string
	// Price is provided in cents
	Price             int32
	Currency          string
	DeliveryTimeSlots []*DeliveryTimeSlot
	PickupPoints      []*PickupPoint
}

type DeliveryTimeSlot struct {
	Start    time.Time
	End      time.Time
	TimeZone time.Location
}

type PickupPoint struct {
	PickupPointId string
	Name          string
	Address       *Address
	Phone         string
	LocationType  *LocationType
	OpeningHours  []*OpeningHours
	IsOperational bool
}

type LocationType int32

const (
	LocationTypeLocker       LocationType = 0
	LocationTypeServicePoint LocationType = 1
	LocationTypeUnknown      LocationType = 2
)

type OpeningHours struct {
	DayOfWeek string
	Opens     string
	Closes    string
}
