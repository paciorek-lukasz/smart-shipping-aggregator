package rpc

import (
	"errors"

	pb "github.com/dzwiedz90/smart-shipping-aggregator/api/shipping"
	"github.com/dzwiedz90/smart-shipping-aggregator/internal/domain"
)

func mapProtoToDomain(pbReq *pb.GetQuotesRequest) *domain.GetQuotesRequest {
	if pbReq == nil {
		return nil
	}

	pbItems := pbReq.GetPackage().GetItems()
	items := make([]*domain.Item, 0, len(pbItems))
	for _, item := range pbItems {
		items = append(items, &domain.Item{
			ItemID:   item.GetItemId(),
			Sku:      item.GetSku(),
			Name:     item.GetName(),
			Price:    item.GetPrice(),
			Quantity: item.GetQuantity(),
		})
	}

	dt := domain.DeliveryType(pbReq.GetDeliveryType().String())

	return &domain.GetQuotesRequest{
		Sender: &domain.Party{
			Name: pbReq.GetSender().GetName(),
			Address: &domain.Address{
				Address:    pbReq.GetSender().GetAddress().GetAddress(),
				PostalCode: pbReq.GetSender().GetAddress().GetPostalCode(),
				City:       pbReq.GetSender().GetAddress().GetCity(),
				Country:    pbReq.GetSender().GetAddress().GetCountry(),
				Longitude:  pbReq.GetSender().GetAddress().GetLongitude(),
				Latitude:   pbReq.GetSender().GetAddress().GetLatitude(),
			},
			Phone: pbReq.GetSender().GetPhone(),
			Email: pbReq.GetSender().GetEmail(),
		},
		Recipient: &domain.Party{
			Name: pbReq.GetRecipient().GetName(),
			Address: &domain.Address{
				Address:    pbReq.GetRecipient().GetAddress().GetAddress(),
				PostalCode: pbReq.GetRecipient().GetAddress().GetPostalCode(),
				City:       pbReq.GetRecipient().GetAddress().GetCity(),
				Country:    pbReq.GetRecipient().GetAddress().GetCountry(),
				Longitude:  pbReq.GetRecipient().GetAddress().GetLongitude(),
				Latitude:   pbReq.GetRecipient().GetAddress().GetLatitude(),
			},
			Phone: pbReq.GetRecipient().GetPhone(),
			Email: pbReq.GetRecipient().GetEmail(),
		},
		Package: &domain.Package{
			Items:      items,
			TotalPrice: pbReq.GetPackage().GetTotalPrice(),
			Currency:   pbReq.GetPackage().GetCurrency(),
			Dimensions: &domain.Dimensions{
				LengthCm:       pbReq.GetPackage().GetDimensions().GetLengthCm(),
				WidthCm:        pbReq.GetPackage().GetDimensions().GetWidthCm(),
				HeightCm:       pbReq.GetPackage().GetDimensions().GetHeightCm(),
				TotalWeightG:   pbReq.GetPackage().GetDimensions().GetTotalWeightG(),
				TotalVolumeCm3: pbReq.GetPackage().GetDimensions().GetTotalVolumeCm3(),
			},
		},
		DeliveryType:  dt,
		LocationTypes: parseLocationTypes(pbReq.GetLocationTypes()),
	}
}

func mapDomainToProto(resp *domain.GetQuotesResponse) (*pb.GetQuotesResponse, error) {
	if resp == nil {
		return nil, errors.New("empty response")
	}

	options := make([]*pb.Option, 0, len(resp.Options))
	for _, option := range resp.Options {
		if option == nil {
			continue
		}

		options = append(options, &pb.Option{
			OptionId:          option.OptionId,
			CarrierProduct:    option.CarrierProduct,
			Price:             option.Price,
			Currency:          option.Currency,
			DeliveryTimeSlots: parseDeliveryTimeslots(option.DeliveryTimeSlots),
			PickupPoints:      parsePickupPoints(option.PickupPoints),
			DeliveryType:      mapDeliveryType(option.DeliveryType),
		})
	}

	return &pb.GetQuotesResponse{
		Options: options,
	}, nil
}

func parseDeliveryTimeslots(deliveryTimeSlots []*domain.DeliveryTimeSlot) []*pb.DeliveryTimeSlot {
	if len(deliveryTimeSlots) == 0 {
		return nil
	}

	timeslots := make([]*pb.DeliveryTimeSlot, 0, len(deliveryTimeSlots))
	for _, ts := range deliveryTimeSlots {
		if ts == nil {
			continue
		}
		timeslots = append(timeslots, &pb.DeliveryTimeSlot{
			Start:    ts.Start.String(),
			End:      ts.End.String(),
			TimeZone: ts.TimeZone.String(),
		})
	}

	return timeslots
}

func parsePickupPoints(pickupPoints []*domain.PickupPoint) []*pb.PickupPoint {
	if len(pickupPoints) == 0 {
		return nil
	}

	pbPickupPoints := make([]*pb.PickupPoint, 0, len(pickupPoints))
	for _, pp := range pickupPoints {
		if pp == nil {
			continue
		}

		var pbAddr *pb.Address
		if pp.Address != nil {
			pbAddr = &pb.Address{
				Address:    pp.Address.Address,
				PostalCode: pp.Address.PostalCode,
				City:       pp.Address.City,
				Country:    pp.Address.Country,
				Longitude:  pp.Address.Longitude,
				Latitude:   pp.Address.Latitude,
			}
		}

		var locType pb.LocationType
		if pp.LocationType != 0 {
			locType = pb.LocationType(pp.LocationType)
		}

		pbPickupPoints = append(pbPickupPoints, &pb.PickupPoint{
			PickupPointId: pp.PickupPointId,
			Name:          pp.Name,
			Address:       pbAddr,
			Phone:         pp.Phone,
			LocationType:  locType,
			OpeningHours:  parseOpeningHours(pp.OpeningHours),
			IsOperational: pp.IsOperational,
		})
	}

	return pbPickupPoints
}

func parseOpeningHours(openingHours []*domain.OpeningHours) []*pb.OpeningHour {
	if len(openingHours) == 0 {
		return nil
	}

	pbOpeningHours := make([]*pb.OpeningHour, 0, len(openingHours))
	for _, oh := range openingHours {
		if oh == nil {
			continue
		}

		pbOpeningHours = append(pbOpeningHours, &pb.OpeningHour{
			DayOfWeek: oh.DayOfWeek,
			Opens:     oh.Opens,
			Closes:    oh.Closes,
		})
	}

	return pbOpeningHours
}

func mapDeliveryType(dt domain.DeliveryType) pb.DeliveryType {
	switch dt {
	case domain.DELIVERY_TYPE_HOME_DELIVERY:
		return pb.DeliveryType_DELIVERY_TYPE_HOME_DELIVERY
	case domain.DELIVERY_TYPE_PICKUP:
		return pb.DeliveryType_DELIVERY_TYPE_PICKUP
	default:
		return 0
	}
}

func parseLocationTypes(types []pb.LocationType) []domain.LocationType {
	var result []domain.LocationType

	for _, t := range types {
		if t == pb.LocationType_TYPE_LOCKER {
			result = append(result, domain.LOCATION_TYPE_LOCKER)
			continue
		}
		if t == pb.LocationType_TYPE_SERVICE_POINT {
			result = append(result, domain.LOCATION_TYPE_SERVICE_POINT)
			continue
		}
	}

	return result
}
