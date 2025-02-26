package jwttest

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

var SampleSecretKey = []byte("SharedSecret")

// GenerateJwt is a test function to generate a JWT structure similar the one received from Keycloak
func GenerateJwt(t *testing.T, iss string, roles map[string][]string) (*jwt.Token, string) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "aae9782e-fee4-4423-bcbd-2252397683fb"
	claims["iss"] = iss
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["azp"] = "inventory.k8s.netic.dk"
	access := make(map[string]any)
	for c, r := range roles {
		access[c] = map[string]any{
			"roles": r,
		}
	}
	claims["resource_access"] = access
	tokenString, err := token.SignedString(SampleSecretKey)
	token.Raw = tokenString
	require.NoError(t, err)
	return token, tokenString
}
