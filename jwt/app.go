package jwt

import (
	"fmt"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type Subjects struct {
	Imports []*jwt.Import
	Exports []*jwt.Export
}

// App creates new self-signed app with necessary subjects in it.
func App(seed string, subj Subjects) (string, error) {
	keyPair, err := nkeys.FromSeed([]byte(seed))
	if err != nil {
		return "", fmt.Errorf("failed to get nkeys: %w", err)
	}

	pubKey, err := keyPair.PublicKey()
	if err != nil {
		return "", fmt.Errorf("failed to get pub key: %w", err)
	}

	claims := jwt.NewAccountClaims(pubKey)
	claims.Name = "developer"

	if len(subj.Exports) == 0 && len(subj.Imports) == 0 {
		return "", ErrNoImportAndExport
	}

	for i := range subj.Imports {
		claims.Imports = append(claims.Imports, subj.Imports[i])
		claims.DefaultPermissions.Sub.Allow.Add(string(subj.Imports[i].Subject))
	}

	for i := range subj.Exports {
		claims.Exports = append(claims.Exports, subj.Exports[i])
		claims.DefaultPermissions.Pub.Allow.Add(string(subj.Exports[i].Subject))
	}

	token, err := claims.Encode(keyPair)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	return token, nil
}
