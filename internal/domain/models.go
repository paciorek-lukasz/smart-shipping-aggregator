package domain

import "time"

type GetQuotesRequest struct {
	Sender       *Party
	Recipient    *Party
	Package      *Package
	DeliveryType *DeliveryType
}

type DeliveryType string

const (
	DELIVERY_TYPE_UNKNOWN       DeliveryType = "DELIVERY_TYPE_UNKNOWN"
	DELIVERY_TYPE_HOME_DELIVERY DeliveryType = "DELIVERY_TYPE_HOME_DELIVERY"
	DELIVERY_TYPE_PICKUP        DeliveryType = "DELIVERY_TYPE_PICKUP"
)

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
	DeliveryType      *DeliveryType
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
	LOCATION_TYPE_UNKNOWN       LocationType = 0
	LOCATION_TYPE_LOCKER        LocationType = 1
	LOCATION_TYPE_SERVICE_POINT LocationType = 2
)

type OpeningHours struct {
	DayOfWeek string
	Opens     string
	Closes    string
}
