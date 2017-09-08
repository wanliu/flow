package graphqlClient

import (
	"time"
)

// {
// 	"data": {
// 		"createOrder": {
// 			"__typename": "OrderPayload",
// 			"order": {
// 				"address": "发明家广场",
// 				"deliveryTime": "0001-01-01T00:00:00Z",
// 				"gifts": [],
// 				"id": "T3JkZXI6Mzc=",
// 				"items": [
// 					{
// 						"price": 0,
// 						"product": {
// 							"name": "250ml伊利纯牛奶",
// 							"picUrl": "",
// 							"price": 0
// 						},
// 						"quantity": 1
// 					}
// 				],
// 				"saler": {
// 					"id": "VXNlcjow",
// 					"name": ""
// 				}
// 			}
// 		}
// 	}
// }

type CreateOrderResponse struct {
	Data   data  `json:"data"`
	Errors []err `json:"errors"`
}

type data struct {
	CreateOrder createOrder `json:"createOrder"`
}

type createOrder struct {
	Typename string `json:"__typename"`
	Order    order  `json:"order"`
}

type order struct {
	Status uint   `json:"status";`
	No     string `json:"no"`

	CustomerId uint `json:"customerId"`
	// Customer   People

	SalerId uint `json:"salerId"`
	// Saler   User

	// OriginAmount float64   `json:"originAmount"`
	// Amount       float64   `json:"amount"`
	Address      string    `json:"address"`
	DeliveryTime time.Time `json:"deliveryTime"`

	// OrderItems []OrderItem `json:"items"`
	// GiftItems  []GiftItem  `json:"gifts"`
}

type err struct {
	Message   string   `json:"message"`
	Locations []string `json:"locations"`
}

func (res CreateOrderResponse) OrderNo() string {
	// if len(res.Errors) == 0 {
	return res.Data.CreateOrder.Order.No
	// } else {
	// 	return ""
	// }
}
