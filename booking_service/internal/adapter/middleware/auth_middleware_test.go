package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"gotest.tools/v3/assert"
)

var (
	testEndpoint = "/test"
	testToken    = "testToken"
	testUser     = userdomain.User{
		ID:    uuid.New().String(),
		Name:  "testUser",
		Email: "123@mail.com",
		Role:  "client",
	}
	testServer *httptest.Server
	testClient *http.Client
)

type mockTokenService struct{}

func (m *mockTokenService) ValidateToken(ctx context.Context, token string) (*userdomain.User, error) {
	if token != testToken {
		return nil, fmt.Errorf("invalid token")
	}

	return &testUser, nil
}

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc(testEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	auth := NewAuthMiddleware(&mockTokenService{})

	testServer = httptest.NewServer(auth.Authorize(mux))
	testClient = testServer.Client()

	log.Println("Running tests...")

	code := m.Run()

	testServer.Close()

	os.Exit(code)
}

func TestAuthMiddleware_OKForAuthorize(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, testServer.URL+testEndpoint, nil)

	if err != nil {
		t.Fatal("Error creating request", err)
	}

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatal("Error sending request", err)
	}
	assert.Assert(t, resp.StatusCode == http.StatusOK)
}

func TestAuthMiddleware_InvalidAuthorize(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, testServer.URL+testEndpoint, nil)
	req.Header.Set("Authorization", "Bearer "+"1")

	if err != nil {
		t.Fatal("Error creating request", err)
	}

	resp, err := testClient.Do(req)

	if err != nil {
		t.Fatal("Error sending request", err)
	}

	assert.Assert(t, resp.StatusCode == http.StatusUnauthorized)
}

func TestAuthMiddleware_Authorize(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, testServer.URL+testEndpoint, nil)

	req.Header.Set("Authorization", "Bearer "+testToken)

	if err != nil {
		t.Fatal("Error creating request", err)
	}

	resp, err := testClient.Do(req)
	if err != nil {
		t.Fatal("Error sending request", err)
	}

	assert.Assert(t, resp.StatusCode == http.StatusOK)
}
