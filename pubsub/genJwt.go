package pubsub

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nkeys"
	"strconv"

	"math/rand"
	"time"
)

type Payload struct {
	Jti  string                 `json:"jti"`
	Iat  int64                  `json:"iat"`
	Iss  string                 `json:"iss"`
	Name string                 `json:"name"`
	Sub  string                 `json:"sub"`
	Nats map[string]interface{} `json:"nats"`
}

func generateJti() string {
	timestamp := time.Now().Unix()
	randomNumber := strconv.FormatFloat(rand.Float64(), 'f', -1, 64)[2:]
	return fmt.Sprintf("%d%s", timestamp, randomNumber)
}

func generateIat() int64 {
	return time.Now().Unix()
}

func getNatsConfig() map[string]interface{} {
	return map[string]interface{}{
		"pub":     map[string]interface{}{},
		"sub":     map[string]interface{}{},
		"subs":    -1,
		"data":    -1,
		"payload": -1,
		"type":    "user",
		"version": 2,
	}
}

func signJwt(payload Payload, account nkeys.KeyPair) (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "ed25519-nkey",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerBytes)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadBytes)

	jwtbase := headerEncoded + "." + payloadEncoded

	signatureBytes, err := account.Sign([]byte(jwtbase))
	if err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(signatureBytes)

	return jwtbase + "." + signature, nil
}

func CreateAppJwt(seed string) (string, error) {
	account, err := nkeys.FromSeed([]byte(seed))
	if err != nil {
		return "", err
	}

	accPubKey, err := account.PublicKey()
	if err != nil {
		return "", err
	}

	payload := Payload{
		Jti:  generateJti(),
		Iat:  generateIat(),
		Iss:  accPubKey,
		Name: "developer",
		Sub:  accPubKey,
		Nats: getNatsConfig(),
	}

	jwt, err := signJwt(payload, account)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
