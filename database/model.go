package database

type DeliveryMessage struct {
	RequestID        string `json:"requestId" bson:"requestId"`
	DeliveryCount    int    `json:"deliveryCount" bson:"deliveryCount,omitempty"`
	UnDeliveryCount  int    `json:"unDeliveryCount" bson:"unDeliveryCount,omitempty"`
	CallServiceCount int    `json:"callServiceCount" bson:"callServiceCount,omitempty"`
	CreatedAt        int64  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt        int64  `json:"updatedAt" bson:"updatedAt,omitempty"`
}
