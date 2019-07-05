package database

type DeliveryMessage struct {
	RequestID       string `json:"requestId" bson:"requestId,omitempty"`
	DeliveryCount   int    `json:"deliveryCount" bson:"deliveryCount,omitempty"`
	UnDeliveryCount int    `json:"unDeliveryCount" bson:"unDeliveryCount,omitempty"`
	CreatedAt       int64  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       int64  `json:"updatedAt" bson:"updatedAt,omitempty"`
}
