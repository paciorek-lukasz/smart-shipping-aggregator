package dhlclient

type DhlHomeApiRequest struct {
	DhlApiKey  string `json:"api_key"`
	Address    string `json:"address"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type DhlHomeApiResponse struct {
}
