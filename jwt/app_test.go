package jwt

import (
	"testing"

	"github.com/nats-io/jwt/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAppSeed = "SAADRH3GLVIAU7KUELAWQ3AFG3XL7T3BB27XC76RGICY2JFWSJ6Q2PXDTQ"

func TestApp(t *testing.T) {
	subjects := Subjects{
		Exports: []*jwt.Export{
			{
				Name:    "subj-foo",
				Subject: "Foo",
				Type:    jwt.Stream,
			},
		},
		Imports: []*jwt.Import{
			{
				Name:    "subj-bar",
				Subject: "Bar",
				Account: "ACEOIUVTS6VQJRUBVV22WCP2M5KTJZK42MJTZQX7JG7CL7WY6DMDLSQZ",
				Type:    jwt.Stream,
			},
		},
	}
	appToken, err := App(testAppSeed, subjects)
	require.NoError(t, err)

	claims, err := jwt.DecodeAccountClaims(appToken)
	require.NoError(t, err)

	assert.EqualValues(t, subjects.Imports, claims.Imports)
	assert.EqualValues(t, subjects.Exports, claims.Exports)
	assert.EqualValues(t, claims.Subject, claims.Issuer)
	assert.True(t, claims.IsSelfSigned())
}
