package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
)

type PaymentSenderImpl struct {
	client  *http.Client
	baseURL string
}

func NewPaymentSender(client *http.Client, baseURL string) *PaymentSenderImpl {
	return &PaymentSenderImpl{client: client, baseURL: baseURL}
}

func (p *PaymentSenderImpl) SendPayment(ctx context.Context, info payment.PaymentInfo) (string, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		slog.Error("Error marshalling payment info", "error", err)
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating payment request", "error", err)
		return "", err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		slog.Error("Error executing payment request", "error", err)
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", "error", err)
		return "", err
	}

	var data struct {
		URL string `json:"url"`
	}

	if err = json.Unmarshal(body, &data); err != nil {
		slog.Error("Error unmarshalling response body", "error", err)
		return "", err
	}

	defer resp.Body.Close()

	return data.URL, nil
}
