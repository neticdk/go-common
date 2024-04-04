package oidc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-jose/go-jose/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestKeySet(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(t, err)

	jwk := jose.JSONWebKey{
		Key:   &key.PublicKey,
		KeyID: "kid",
	}
	keys := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{jwk},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/openid-configuration" {
			w.WriteHeader(http.StatusOK)
			cfg := fmt.Sprintf(`{"jwks_uri":"http://%s/keyset"}`, r.Host)
			w.Write([]byte(cfg))
		} else if r.URL.Path == "/keyset" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(keys)
		}
	}))
	defer srv.Close()

	kf, err := NewKeyfunc(context.Background(), srv.URL)
	assert.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Issuer: srv.URL,
	})
	token.Header["kid"] = jwk.KeyID
	jws, err := token.SignedString(key)
	assert.NoError(t, err)

	_, err = jwt.Parse(jws, kf)
	assert.NoError(t, err)
}
