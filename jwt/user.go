package jwt

import (
	"fmt"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

// AppUser creates jwt for user of passed app.
func AppUser(appSeed, userSeed string) (string, error) {
	userKeyPair, err := nkeys.FromSeed([]byte(userSeed))
	if err != nil {
		return "", fmt.Errorf("failed to get user nkeys: %w", err)
	}

	userPubKey, err := userKeyPair.PublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to get user pub key: %w", err)
	}

	appKeyPair, err := nkeys.FromSeed([]byte(appSeed))
	if err != nil {
		return "", fmt.Errorf("failed to get app pub key: %w", err)
	}

	claims := jwt.NewUserClaims(userPubKey)
	claims.Name = "app-user"

	token, err := claims.Encode(appKeyPair)
	if err != nil {
		return "", fmt.Errorf("failed to create user token: %w", err)
	}
	return token, nil
}
