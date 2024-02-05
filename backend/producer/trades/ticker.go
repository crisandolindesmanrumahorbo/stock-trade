package trades

type Ticker struct {
	Symbol string `json:"s"`
	Price float64 `json:"p"`
	Volume uint64 `json:"v"`
	Time uint64 `json:"t"`
}