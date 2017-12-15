package graphqlClient

import (
	"strconv"
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

type Product struct {
	PicUrl string `json:"picUrl"`
	Name   string `json:"name"`
}

type OrderItem struct {
	ProductId   uint    `json:"productId"`
	OriginPrice float64 `json:"originPrice"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	Product     Product `json:"product"`
	ProductName string  `json:"productName"`
}

type GiftItem struct {
	ProductId   uint    `json:"productId"`
	OriginPrice float64 `json:"originPrice"`
	Quantity    int     `json:"quantity"`
	Product     Product `json:"product"`
	ProductName string  `json:"productName"`
}

type order struct {
	Status uint   `json:"status";`
	No     string `json:"no"`
	Note   string `json:"note"`

	CustomerId uint `json:"customerId"`
	// Customer   People

	SalerId uint `json:"salerId"`
	// Saler   User

	// OriginAmount float64   `json:"originAmount"`
	// Amount       float64   `json:"amount"`
	Address      string    `json:"address"`
	DeliveryTime time.Time `json:"deliveryTime"`

	OrderItems []OrderItem `json:"items"`
	GiftItems  []GiftItem  `json:"gifts"`
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

func (res CreateOrderResponse) Items() []OrderItem {
	items := res.Data.CreateOrder.Order.OrderItems
	result := make([]OrderItem, 0, len(items))

	for _, item := range items {
		item.ProductName = item.Product.Name
		result = append(result, item)
	}

	return result
}

func (res CreateOrderResponse) Gifts() []GiftItem {
	gifts := res.Data.CreateOrder.Order.GiftItems
	result := make([]GiftItem, 0, len(gifts))

	for _, gift := range gifts {
		gift.ProductName = gift.Product.Name
		result = append(result, gift)
	}

	return result
}

func (res CreateOrderResponse) Note() string {
	return res.Data.CreateOrder.Order.Note
}

func (res CreateOrderResponse) AnswerBody() string {
	desc := ""

	for _, item := range res.Items() {
		desc = desc + item.ProductName + " " + strconv.Itoa(item.Quantity) + "件\n"
	}

	gifts := res.Gifts()
	if len(gifts) > 0 {
		desc = desc + "申请的赠品:\n"

		for _, g := range gifts {
			desc = desc + g.ProductName + " " + strconv.Itoa(g.Quantity) + "件\n"
		}
	}

	desc = desc + "时间:" + res.Data.CreateOrder.Order.DeliveryTime.Format("2006年01月02日") + "\n"

	if res.Note() != "" {
		desc = desc + "备注：" + res.Note() + "\n"
	}

	return desc
}
