package userkey

import (
	"errors"
	"net/http"

	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type UserKey struct{}

func ExtractUserFromReq(r *http.Request) (*userdomain.User, error) {
	value := r.Context().Value(UserKey{})
	if value == nil {
		return nil, errors.New("user key not found in context")
	}
	user, ok := value.(*userdomain.User)
	if !ok {
		return nil, errors.New("user is not valid in context")
	}

	return user, nil
}
