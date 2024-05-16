package service

import "github.com/tanush-128/openzo_backend/online_order/internal/pb"

var acceptedNotification = pb.Notification{
	Title:     "Order Accepted",
	Body:      "Your order has been accepted!",
	ImageURL:  "",
	ActionURL: "",
	Token:     "",
}

var rejectedNotification = pb.Notification{
	Title:     "Order Rejected",
	Body:      "Your order has been rejected!",
	ImageURL:  "",
	ActionURL: "",
	Token:     "",
}

var outForDeliveryNotification = pb.Notification{
	Title:     "Order Out For Delivery",
	Body:      "Your order is out for delivery!",
	ImageURL:  "",
	ActionURL: "",
	Token:     "",
}

var cancelledNotification = pb.Notification{
	Title:     "Order Cancelled",
	Body:      "Your order has been cancelled!",
	ImageURL:  "",
	ActionURL: "",
	Token:     "",
}

var deliveredNotification = pb.Notification{
	Title:     "Order Delivered",
	Body:      "Your order has been delivered!",
	ImageURL:  "",
	ActionURL: "",
	Token:     "",
}
