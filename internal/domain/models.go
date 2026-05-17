package domain

import "time"

type GetQuotesRequest struct {
	Sender        *Party         `json:"sender"`
	Recipient     *Party         `json:"recipient"`
	Package       *Package       `json:"package"`
	DeliveryType  DeliveryType   `json:"delivery_type"`
	LocationTypes []LocationType `json:"location_types"`
}

type DeliveryType string

const (
	DELIVERY_TYPE_UNKNOWN       DeliveryType = "DELIVERY_TYPE_UNKNOWN"
	DELIVERY_TYPE_HOME_DELIVERY DeliveryType = "DELIVERY_TYPE_HOME_DELIVERY"
	DELIVERY_TYPE_PICKUP        DeliveryType = "DELIVERY_TYPE_PICKUP"
)

type Party struct {
	Name    string   `json:"name"`
	Address *Address `json:"address"`
	Phone   string   `json:"phone"`
	Email   string   `json:"email"`
}

type Address struct {
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Longitude  string `json:"longitude"`
	Latitude   string `json:"latitude"`
}

type Package struct {
	Items []*Item `json:"items"`
	// Price is provided in cents
	TotalPrice int32       `json:"total_price"`
	Currency   string      `json:"currency"`
	Dimensions *Dimensions `json:"dimensions"`
}

type Item struct {
	ItemID int32  `json:"item_id"`
	Sku    string `json:"sku"`
	Name   string `json:"name"`
	// Price is provided in cents
	Price    int32 `json:"price"`
	Quantity int32 `json:"quantity"`
}

type Dimensions struct {
	LengthCm       int32 `json:"length"`
	WidthCm        int32 `json:"width"`
	HeightCm       int32 `json:"height"`
	TotalWeightG   int32 `json:"weight"`
	TotalVolumeCm3 int32 `json:"volume"`
}

type GetOptionsResponse struct {
	Options []*Option `json:"options"`
}

type GetQuotesResponse struct {
	Options *Option `json:"option"`
}

type Option struct {
	OptionId       int32  `json:"option_id"`
	CarrierProduct string `json:"carrier_product"`
	// Price is provided in cents
	Price             int32               `json:"price"`
	Currency          string              `json:"currency"`
	DeliveryTimeSlots []*DeliveryTimeSlot `json:"delivery_timeslots"`
	PickupPoints      []*PickupPoint      `json:"pickup_point"`
	DeliveryType      DeliveryType        `json:"delivery_type"`
}

type DeliveryTimeSlot struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	TimeZone time.Location `json:"time_zone"`
}

type PickupPoint struct {
	PickupPointId string          `json:"pickup_point_id"`
	Name          string          `json:"name"`
	Address       *Address        `json:"address"`
	Phone         string          `json:"phone"`
	LocationType  LocationType    `json:"location_type"`
	OpeningHours  []*OpeningHours `json:"opening_hours"`
	IsOperational bool            `json:"is_operational"`
}

type LocationType int32

const (
	LOCATION_TYPE_UNKNOWN       LocationType = 0
	LOCATION_TYPE_LOCKER        LocationType = 1
	LOCATION_TYPE_SERVICE_POINT LocationType = 2
	LOCATION_TYPE_POSTOFFICE    LocationType = 3
)

type OpeningHours struct {
	DayOfWeek string `json:"day_of_week"`
	Opens     string `json:"opens"`
	Closes    string `json:"closes"`
}
