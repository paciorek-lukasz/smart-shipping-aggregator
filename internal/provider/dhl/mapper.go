package dhl

import "github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"

func toDhlRequest(req *domain.GetQuotesRequest) *DhlApiRequest {
	return &DhlApiRequest{}
}
