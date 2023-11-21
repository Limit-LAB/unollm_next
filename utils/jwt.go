package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"crypto/hmac"
	"crypto/sha256"
)

func CreateJWTToken(apiKey string, expire time.Duration) (string, error) {
	sp := strings.Split(apiKey, ".")
	if len(sp) < 2 {
		return "", errors.New("invalid API key format")
	}
	key := sp[0]
	secret := sp[len(sp)-1]

	now := time.Now().Unix()
	exp := now + int64(expire.Seconds())

	header := map[string]interface{}{
		"alg":       "HS256",
		"sign_type": "SIGN",
		"typ":       "JWT",
	}

	claims := map[string]interface{}{
		"api_key":   key,
		"exp":       exp,
		"timestamp": now,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerBase64 := base64.StdEncoding.EncodeToString(headerJSON)
	claimsBase64 := base64.StdEncoding.EncodeToString(claimsJSON)

	kkey := []byte(secret)
	h := hmac.New(sha256.New, kkey)
	h.Write([]byte(headerBase64))
	h.Write([]byte("."))
	h.Write([]byte(claimsBase64))
	signature := h.Sum(nil)

	tokenString := headerBase64 + "." + claimsBase64 + "." + base64.StdEncoding.EncodeToString(signature)
	return tokenString, nil
}
