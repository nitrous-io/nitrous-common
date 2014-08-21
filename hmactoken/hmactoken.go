package hmactoken

import (
	"errors"
	"time"

	"github.com/nitrous-io/nitrous-common/msgverifier"
)

type Token struct {
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

func Generate(permissions []string, secret string) (string, error) {
	createdAt := time.Now().UTC()
	expiresAt := createdAt.Add(DefaultValidityPeriod)
	return msgverifier.New(secret).Generate(&Token{
		Permissions: permissions,
		CreatedAt:   &createdAt,
		ExpiresAt:   &expiresAt,
	})
}

func Verify(permissions []string, token string, secret string) error {
	var t *Token
	err := msgverifier.New(secret).Verify(token, &t)
	if err != nil {
		return ErrInvalidToken
	}
	if t.ExpiresAt.Before(time.Now()) {
		return ErrTokenExpired
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
			return ErrPermissionDenied
		}
	}
	return nil
}
