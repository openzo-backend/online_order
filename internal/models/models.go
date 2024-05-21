package models

import "time"

// type OnlineOrder struct {
// 	ID string `json:"id" gorm:"primaryKey"` // The ID field is also the UUID for the store

// }

type OrderStatus string

const (
	OrderNotPlaced OrderStatus = "not_placed"
	OrderPlaced    OrderStatus = "placed"
	OrderAccepted  OrderStatus = "accepted"
	OrderRejected  OrderStatus = "rejected"
	OrderOutForDel OrderStatus = "out_for_delivery"
	OrderCancelled OrderStatus = "cancelled"
	OrderDelivered OrderStatus = "delivered"
)

type OnlineOrder struct {
	ID          string            `json:"id" gorm:"primaryKey"`
	OrderItems  []OnlineOrderItem `json:"order_items"`
	StoreID     string            `json:"store_id"`
	Customer    OnlineCustomer    `json:"customer"`
	OrderTime   time.Time         `json:"order_time" gorm:"autoCreateTime"`
	OrderStatus OrderStatus       `json:"order_status"`
	TotalAmount float64           `json:"total_amount"`
}

type OnlineOrderItem struct {
	ID            int    `json:"id"`
	ProductID     string `json:"product_id" gorm:"size:36"`
	OnlineOrderId string `json:"sale_id" gorm:"size:36"`
	Quantity      int    `json:"quantity"`
}

type OnlineCustomer struct {
	ID            int    `json:"id"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	UserDataId    string `json:"user_data_id" gorm:"size:36"`
	OnlineOrderID string `json:"sale_id" gorm:"size:36"`
}
