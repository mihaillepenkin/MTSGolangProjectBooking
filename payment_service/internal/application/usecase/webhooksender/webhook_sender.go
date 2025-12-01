package webhooksender

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/payment_service/internal/domain/payment/object"
)

type WebhookSender struct {
	client *http.Client
}

func NewWebhookSender(client *http.Client) *WebhookSender {
	return &WebhookSender{client: client}
}

func (w *WebhookSender) Send(ctx context.Context, URL string, response *object.PaymentResponse) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		slog.Error("Error marshalling response", "error", err)
		return err
	}
	_, err = http.NewRequestWithContext(ctx, "POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating request", "error", err)
		return err
	}

	/*resp, err := w.client.Do(req)
	if err != nil {
		slog.Error("Error executing request", "error", err)
		return err
	}*/

	//defer resp.Body.Close()
	return nil
}
