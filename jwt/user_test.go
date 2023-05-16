package jwt

import (
	"testing"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testUserSeed = "SUAPBVNIX27GUU7KH5CVZT4DSYQZ7JKQWUICMHDDUEQPGYB7HS5BVK6XE4"

func TestAppUser(t *testing.T) {
	userToken, err := AppUser(testAppSeed, testUserSeed)
	require.NoError(t, err)

	claims, err := jwt.DecodeUserClaims(userToken)
	require.NoError(t, err)

	appKeyPair, err := nkeys.FromSeed([]byte(testAppSeed))
	require.NoError(t, err)

	appPubKey, err := appKeyPair.PublicKey()
	require.NoError(t, err)

	assert.EqualValues(t, string(appPubKey), claims.Issuer)
}
