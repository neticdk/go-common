package jwttest

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

var SampleSecretKey = []byte("SharedSecret")

// GenerateJwt is a test function to generate a JWT structure similar the one received from Keycloak
func GenerateJwt(t *testing.T, roles map[string][]string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "aae9782e-fee4-4423-bcbd-2252397683fb"
	claims["iss"] = "https://keycloak.netic.dk/auth/realms/mcs"
	claims["exp"] = time.Now().Add(10 * time.Minute)
	access := make(map[string]interface{})
	for c, r := range roles {
		access[c] = map[string]interface{}{
			"roles": r,
		}
	}
	claims["resource_access"] = access
	tokenString, err := token.SignedString(SampleSecretKey)
	require.NoError(t, err)
	return tokenString
}
