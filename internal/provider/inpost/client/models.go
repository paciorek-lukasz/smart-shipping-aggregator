package client

type InpostHomeApiRequest struct {
	InpostToken     string        `json:"inpost_token"`
	SenderDetails   *PartyDetails `json:"sender_details"`
	ReceiverDetails *PartyDetails `json:"receiver_details"`
}

type PartyDetails struct {
	StreetAddress string `json:"street_address"`
	ZipCode       string `json:"zip_code"`
	Locality      string `json:"locality"`
	CountryIso    string `json:"country_iso"`
}

type InpostHomeApiResponse struct {
	DeliveryStart string `json:"delivery_start"`
	DeliveryEnd   string `json:"delivery_end"`
	ShipmentCost  int32  `json:"shipment_cost"`
	IsoCurrency   string `json:"iso_currency"`
}

type InpostPickupApiRequest struct {
	InpostToken     string           `json:"inpost_token"`
	SenderDetails   *PartyDetails    `json:"sender_details"`
	ReceiverDetails *PartyDetails    `json:"receiver_details"`
	MachinesLimit   int32            `json:"machines_limit"`
	RadiusKm        int32            `json:"radius_km"`
	MachineCategory []InpostCategory `json:"machine_category"`
}

type InpostCategory string

const (
	INPOST_CATEGORY_ALL     InpostCategory = "all"
	INPOST_CATEGORY_LOCKER  InpostCategory = "locker"
	INPOST_CATEGORY_PARTNER InpostCategory = "partner"
)

type InpostPickupApiResponse struct {
	DeliveryStart string          `json:"delivery_start"`
	DeliveryEnd   string          `json:"delivery_end"`
	ShipmentCost  int32           `json:"shipment_cost"`
	IsoCurrency   string          `json:"iso_currency"`
	Machines      []*ParcelLocker `json:"machines"`
}

type ParcelLocker struct {
	MachineId  string     `json:"machine_id"`
	Name       string     `json:"name"`
	Locality   string     `json:"locality"`
	ZipCode    string     `json:"zip_code"`
	CountryIso string     `json:"country_iso"`
	Latitude   string     `json:"latitude"`
	Longitude  string     `json:"longitude"`
	Address    string     `json:"address"`
	Status     string     `json:"status"`
	Hours      []*Opening `json:"hours"`
	Accepts    []string   `json:"accepts"`
}

const (
	InpostStatusActive   = "active"
	InpostStatusInactive = "inactive"
	InpostAcceptPack     = "pack"
	InpostAcceptDoc      = "document"
)

type Opening struct {
	DayIndex  string `json:"day_index"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}
