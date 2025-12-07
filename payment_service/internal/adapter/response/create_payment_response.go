package response

type CreatePaymentResponse struct {
	URL       string `json:"url"`
	PaymentID string `json:"payment_id"`
}
