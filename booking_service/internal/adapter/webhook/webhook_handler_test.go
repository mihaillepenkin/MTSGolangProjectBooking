package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/webhook/request"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"gotest.tools/v3/assert"
)

var (
	testServer         *httptest.Server
	testClient         *http.Client
	testWebhookRequest = request.WebhookRequest{
		Price:    100,
		Currency: "USD",
		Info: object.BookingInfo{
			User: userdomain.User{
				Name:  "1",
				Email: "123@mail.com",
				Role:  "client",
				ID:    uuid.New().String(),
			},
			HotelID:    1,
			RoomID:     1,
			CheckIn:    time.Now().UTC(),
			CheckOut:   time.Now().UTC(),
			HotelName:  "1",
			RoomNumber: 1,
		},
		Status: "Success",
	}
	testEndpoint = "/test"
)

type mockSaver struct{}

func (m *mockSaver) BookRoom(ctx context.Context, bookingInfo *object.BookingInfo) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockSaver) DeleteBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	if bookingInfo.HotelName != testWebhookRequest.Info.HotelName {
		return fmt.Errorf("error delete")
	}

	return nil
}

func (m *mockSaver) ConfirmBooking(ctx context.Context, bookingInfo *object.BookingInfo) error {
	if bookingInfo.HotelName != testWebhookRequest.Info.HotelName {
		return fmt.Errorf("error confirm")
	}

	return nil
}

func TestMain(m *testing.M) {
	webhookHandler := NewWebhookHandler(&mockSaver{})

	mux := http.NewServeMux()

	mux.HandleFunc(testEndpoint, webhookHandler.ServeWebhook)

	testServer = httptest.NewServer(mux)

	testClient = testServer.Client()

	code := m.Run()

	testServer.Close()

	os.Exit(code)
}

func TestWebhookHandler_ServeWebhook(t *testing.T) {
	tests := []struct {
		name           string
		webhookRequest request.WebhookRequest
		expectedStatus int
	}{
		{
			name: "internal error for success",
			webhookRequest: request.WebhookRequest{
				Info: object.BookingInfo{
					HotelName: "2",
				},
				Status: "Success",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "internal error for failed",
			webhookRequest: request.WebhookRequest{
				Info: object.BookingInfo{
					HotelName: "2",
				},
				Status: "Failed",
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "OK for success",
			webhookRequest: testWebhookRequest,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OK for failed",
			webhookRequest: testWebhookRequest,
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		jsonBytes, err := json.Marshal(test.webhookRequest)

		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, testServer.URL+testEndpoint, bytes.NewBuffer(jsonBytes))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := testClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, resp.StatusCode, test.expectedStatus)
	}
}

func TestWebhookHandler_ErrForSuccessDeleteBooking(t *testing.T) {

}
