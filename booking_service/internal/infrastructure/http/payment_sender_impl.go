package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
)

type PaymentSenderImpl struct {
	client  *http.Client
	baseURL string
	logger  *slog.Logger
}

func NewPaymentSender(client *http.Client, baseURL string) *PaymentSenderImpl {
	return &PaymentSenderImpl{client: client, baseURL: baseURL, logger: slog.Default().With("component", "payment_sender")}
}

func (p *PaymentSenderImpl) SendPayment(ctx context.Context, info *payment.PaymentInfo) (*payment.PaymentResponse, error) {
	jsonData, err := json.Marshal(info)
	if err != nil {
		p.logger.Error("Error marshalling http info", "error", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		p.logger.Error("Error creating http request", "error", err)
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		p.logger.Error("Error executing http request", "error", err)
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		p.logger.Error("Payment request failed with status",
			"status", resp.StatusCode)

		return nil, fmt.Errorf("http failed with status %d",
			resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		p.logger.Error("Error reading response body", "error", err)
		return nil, err
	}

	response := &payment.PaymentResponse{}

	if err = json.Unmarshal(body, response); err != nil {
		p.logger.Error("Error unmarshalling response body", "error", err)
		return nil, err
	}

	defer resp.Body.Close()

	return response, nil
}
