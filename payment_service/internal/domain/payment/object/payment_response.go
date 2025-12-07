package object

type PaymentResponse struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Metadata string  `json:"metadata"`
	Status   string  `json:"status"`
}
