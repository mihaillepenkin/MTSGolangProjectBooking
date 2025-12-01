package payment

import error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/error"

type Payment struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	URL      string  `json:"url"`
	Metadata string  `json:"metadata"`
}

func ValidatePayment(payment *Payment) error {
	if payment.Price <= 0 || payment.Currency == "" {
		return error2.ErrPaymentIsInvalid
	}
	return nil
}
