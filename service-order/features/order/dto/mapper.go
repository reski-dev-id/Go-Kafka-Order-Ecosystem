package dto

import "order-service/features/order/entity"

func ToOrderResponse(
	order *entity.Order,
) OrderResponse {
	return OrderResponse{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		ProductName:  order.ProductName,
		Quantity:     order.Quantity,
		Amount:       order.Amount,
		Status:       order.Status,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
	}
}

func ToOrderResponses(
	orders []entity.Order,
) []OrderResponse {
	results := make([]OrderResponse, 0)

	for _, order := range orders {
		results = append(
			results,
			ToOrderResponse(&order),
		)
	}

	return results
}
