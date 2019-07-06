package database

// DeliveryMessage type
type DeliveryMessage struct {
	RequestID       string `json:"requestId" bson:"requestId,omitempty"`
	ChannelID       string `json:"channelId" bson:"channelId,omitempty"`
	LineRequestID   string `json:"lineRequestId" bson:"lineRequestId,omitempty"`
	DeliveryCount   int    `json:"deliveryCount" bson:"deliveryCount,omitempty"`
	UnDeliveryCount int    `json:"unDeliveryCount" bson:"unDeliveryCount"`
	CreatedAt       int64  `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       int64  `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func (d DeliveryMessage) IsDone() bool {
	return d.UnDeliveryCount == 0
}
