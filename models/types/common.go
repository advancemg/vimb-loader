package types

type ItemId struct {
	Id string `json:"id"`
}
type ItemCnl struct {
	Cnl string `json:"Cnl"`
}
type ItemAdtID struct {
	AdtID string `json:"AdtID"`
}
type CommInMpl struct {
	ItemId
	Inventory string `json:"Inventory"`
}
