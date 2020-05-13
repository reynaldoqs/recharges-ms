package entity

// RechargeMsg it needs to be in a package
type RechargeMsg struct {
	IDRecharge  string `json:"idRecharge"`
	PhoneNumber uint   `json:"phoneNumber" validate:"required,gte=9999999,lte=100000000"`
	Company     string `json:"company" validate:"required,oneof=entel viva tigo"`
	CardNumber  uint   `json:"cardNumber" validate:"required,min=99999999999999,max=10000000000000000"`
	CreatedAt   int64  `json:"createdAt"`
	IDResolver  string `json:"idResolver"`
	ResolvedAt  int64  `json:"resolvedAt"`
}
