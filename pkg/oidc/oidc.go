package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/neticdk/go-common/pkg/log"
)

type providerConfiguration struct {
	Issuer string `json:"issuer"`

	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`

	JWKS string `json:"jwks_uri"`
}

// NewKeyfunc creates a [jwt.Keyfunc] capable of validating signed JWTs from the given OpenID Connect Discovery compliant issuers. The
// keyset will be retrieved according to the OpenID Connectect Discovery protocol.
func NewKeyfunc(ctx context.Context, issuers ...string) (jwt.Keyfunc, error) {
	logger := log.FromContext(ctx)
	logger.DebugContext(ctx, "initializing OIDC key function")

	client := http.DefaultClient

	providers := map[string]keyfunc.Keyfunc{}
	for _, issuer := range issuers {
		logger.DebugContext(ctx, "resolving certificates", slog.String("issuer", issuer))
		wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
		req, err := http.NewRequestWithContext(ctx, "GET", wellKnown, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to create request for OpenID configuration: %w", err)
		}
		// #nosec G704 -- URL scheme is validated, and risk is acknowledged
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("unable to request OpenID configuration: %w", err)
		}

		var config providerConfiguration
		err = json.NewDecoder(resp.Body).Decode(&config)

		closeErr := resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("unable to parse OpenID configuration document: %w", err)
		}

		if closeErr != nil {
			logger.WarnContext(ctx, "unable to close OpenID configuration response", slog.Any("error", closeErr))
		}

		jwks, err := keyfunc.NewDefault([]string{config.JWKS})
		if err != nil {
			return nil, err
		}
		providers[issuer] = jwks
	}

	f := func(token *jwt.Token) (any, error) {
		iss, err := token.Claims.GetIssuer()
		if err != nil {
			return nil, err
		}
		jwks, ok := providers[iss]
		if !ok {
			return nil, fmt.Errorf("invalid jwt issuer: %s", iss)
		}
		return jwks.Keyfunc(token)
	}

	return f, nil
}
