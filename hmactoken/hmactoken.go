package hmactoken

import (
	"errors"
	"time"

	"github.com/nitrous-io/nitrous-common/msgverifier"
)

type Token struct {
	UserID      int64      `json:"user_id"`
	Permissions []string   `json:"permissions"`
	CreatedAt   *time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// configurable
var DefaultValidityPeriod time.Duration = 2 * time.Minute

var (
	ErrInvalidToken     = errors.New("hmactoken: invalid token")
	ErrPermissionDenied = errors.New("hmactoken: permission denied")
	ErrTokenExpired     = errors.New("hmactoken: token expired")
)

func Generate(userID int64, permissions []string, secret string) (string, error) {
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(DefaultValidityPeriod)
	return msgverifier.New(secret).Generate(&Token{
		UserID:      userID,
		Permissions: permissions,
		CreatedAt:   &createdAt,
		ExpiresAt:   &expiresAt,
	})
}

func Verify(permissions []string, token string, secret string) (*Token, error) {
	var t *Token
	err := msgverifier.New(secret).Verify(token, &t)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if t.ExpiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	for _, p := range permissions {
		found := false
		for _, q := range t.Permissions {
			if p == q {
				found = true
				break
			}
		}
		if !found {
			return nil, ErrPermissionDenied
		}
	}
	return t, nil
}
