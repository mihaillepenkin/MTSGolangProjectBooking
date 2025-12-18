package object

type PaymentResponse struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Metadata any     `json:"metadata"`
	Status   string  `json:"status"`
}
