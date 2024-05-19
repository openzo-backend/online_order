package models

// type OnlineOrder struct {
// 	ID string `json:"id" gorm:"primaryKey"` // The ID field is also the UUID for the store

// }

type OrderStatus string

const (
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
	OrderTime   string            `json:"order_time"`
	OrderStatus OrderStatus       `json:"order_status"`
	TotalAmount float64           `json:"total_amount"`
}

type OnlineOrderItem struct {
	ID            int    `json:"id"`
	ProductID     string `json:"product_id"`
	OnlineOrderId string `json:"sale_id"`
	Quantity      int    `json:"quantity"`
}

type OnlineCustomer struct {
	ID            int    `json:"id"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	UserDataId    string `json:"user_data_id"`
	OnlineOrderID string `json:"sale_id"`
}
