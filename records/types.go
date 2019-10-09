package records

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

//=============================================================================
type HandlerType uint8

//=============================================================================
const (
	ImpressionHandler HandlerType = iota
	ClickHandler
	CompletionHandler
)

//=============================================================================
type AdType int64

//=============================================================================
func (ad AdType) String() string {
	return fmt.Sprintf("%d", ad)
}

//=============================================================================
type Recorder interface {
	Record() *bson.M
}

//=============================================================================
type Request struct {
	Handler HandlerType
	Body    []byte
}

//=============================================================================
type Completion struct {
	DateTime      string `json:"date_time"`
	TransactionID string `json:"transaction_id"`
}

//=============================================================================
func (c *Completion) Record() *bson.M {

	return &bson.M{
		"date-time":      c.DateTime,
		"transaction-id": c.TransactionID,
	}
}

//=============================================================================
type Impression struct {
	DateTime      string `json:"date_time"`
	TransactionID string `json:"transaction_id"`
	Adtype        AdType `json:"ad_type"`
	UserID        string `json:"user_id"`
}

//=============================================================================
func (i *Impression) Record() *bson.M {

	return &bson.M{
		"date-time":      i.DateTime,
		"transaction-id": i.TransactionID,
		"ad-type":        i.Adtype.String(),
		"user-id":        i.UserID,
	}
}

//=============================================================================
type Click struct {
	DateTime      string `json:"date_time"`
	TransactionID string `json:"transaction_id"`
	Adtype        AdType `json:"ad_type"`
	TimeToClick   string `json:"time_to_click"`
	UserId        string `json:"user_id"`
}

//=============================================================================
func (c *Click) Record() *bson.M {

	return &bson.M{
		"date-time":      c.DateTime,
		"transaction-id": c.TransactionID,
		"ad-type":        c.Adtype.String(),
		"time-to-click":  c.TimeToClick,
		"user-id":        c.UserId,
	}
}
