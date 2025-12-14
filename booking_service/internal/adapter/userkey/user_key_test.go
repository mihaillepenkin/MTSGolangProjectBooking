package userkey

import (
	"context"
	"net/http"
	"testing"

	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestExtractUserFromReq(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/test", nil)

	if err != nil {
		t.Fatal(err)
	}

	testUser := &userdomain.User{Email: "123@mail.com",
		Name: "1", Role: "admin"}

	ctx := context.WithValue(req.Context(), UserKey{}, testUser)
	newReq := req.WithContext(ctx)

	user, err := ExtractUserFromReq(newReq)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)

	_, err = ExtractUserFromReq(req)
	assert.Error(t, err)
}
