package webhook

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/webhook/request"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
)

type WebhookHandler struct {
	bookingService booking.Saver
	logger         *slog.Logger
}

func NewWebhookHandler(bookingService booking.Saver) *WebhookHandler {
	return &WebhookHandler{bookingService: bookingService, logger: slog.Default().With("component", "webhook_handler")}
}

func (wh *WebhookHandler) ServeWebhook(w http.ResponseWriter, r *http.Request) {
	wh.logger.Debug("WebhookHandler.ServeWebhook called")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		wh.logger.Error("WebhookHandler.ServeWebhook failed to read body", "error", err)
		http.Error(w, "failed to serve", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var webhookRequest request.WebhookRequest
	err = json.Unmarshal(body, &webhookRequest)
	if err != nil {
		wh.logger.Error("WebhookHandler.ServeWebhook failed to unmarshal body", "error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if webhookRequest.Status == "Success" {
		err = wh.bookingService.ConfirmBooking(r.Context(), &webhookRequest.Info)
		if err != nil {
			wh.logger.Error("WebhookHandler.ServeWebhook failed to confirm postgres", "error", err)
			http.Error(w, "failed to serve", http.StatusInternalServerError)
			return
		}
	} else if webhookRequest.Status == "Failed" {
		err = wh.bookingService.DeleteBooking(r.Context(), &webhookRequest.Info)
		if err != nil {
			wh.logger.Error("WebhookHandler.ServeWebhook failed to delete postgres", "error", err)
			http.Error(w, "failed to serve", http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte("Successfully served webhook"))
	if err != nil {
		wh.logger.Error("WebhookHandler.ServeWebhook failed to write response", "error", err)
		return
	}
}
