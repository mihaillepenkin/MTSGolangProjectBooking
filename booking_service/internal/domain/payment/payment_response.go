package payment

type PaymentResponse struct {
	PaymentID string `json:"payment_id"`
	URL       string `json:"url"`
}
