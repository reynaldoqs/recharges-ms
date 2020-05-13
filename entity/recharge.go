package entity

type rechargeStatus int

const (
	// Pending Recharge request doesn't taken"
	Pending rechargeStatus = 1

	// Taken Recharge request was taken for a CPR"
	Taken rechargeStatus = 2

	// Resolved Recharge request was resolved by a CPR"
	Resolved rechargeStatus = 3
)

// Recharge represents a basic model for each recharge which clients sends
type Recharge struct {
	ID          string         `json:"id" bson:"_id"`
	PhoneNumber uint           `json:"phoneNumber" validate:"required,gte=9999999,lte=100000000"`
	Company     string         `json:"company" validate:"required,oneof=entel viva tigo"`
	CardNumber  uint           `json:"cardNumber" validate:"required,min=99999999999999,max=10000000000000000"`
	Status      rechargeStatus `json:"status"`
	IDResolver  string         `json:"idResolver"`
	Mount       int            `json:"mount" validate:"gte=10"`
	CreatedAt   int64          `json:"createdAt"`
	ResolvedAt  int64          `json:"resolvedAt"`
}
