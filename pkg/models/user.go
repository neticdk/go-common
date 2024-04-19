package models

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type User struct {
	Name            string
	Email           string
	Username        string
	Roles           []string
	Scopes          []string
	IsAuthenticated bool
	IsAdmin         bool
	OAuth2Token     *oauth2.Token
	JWTToken        *jwt.Token
}

func (u User) Validate() (errs []error) {
	if u.Name == "" {
		errs = append(errs, errors.New("name is empty"))
	}
	if u.Email == "" {
		errs = append(errs, errors.New("email is empty"))
	}
	if u.Username == "" {
		errs = append(errs, errors.New("username is empty"))
	}
	return
}

func (u User) HasRole(role string) bool {
	return u.HasRolesOr([]string{role})
}

func (u User) HasRolesAnd(roles []string) bool {
	for _, neededRole := range roles {
		found := false

		for _, role := range u.Roles {
			if role == string(neededRole) {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func (u User) HasRolesOr(roles []string) bool {
	for _, neededRole := range roles {
		for _, role := range u.Roles {
			if role == string(neededRole) {
				return true
			}
		}
	}

	return false
}

// UserFromJWTToken builds a user model from the claims of a JWT token
// clientID is used to find the client from the resource_access claim where the roles are found
func UserFromJWTToken(token *jwt.Token, clientID string) (User, error) {
	user := User{
		Roles:           make([]string, 0),
		IsAdmin:         false,
		IsAuthenticated: true,
		JWTToken:        token,
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("unable to use token claims as jwt.MapClaims")
	}

	user.Name = stringClaim(claims, "name")
	user.Email = stringClaim(claims, "email")
	user.Username = stringClaim(claims, "preferred_username")

	resourceAccess, ok := claims["resource_access"].(map[string]interface{})
	if !ok {
		resourceAccess = map[string]interface{}{}
	}

	if clientResources, ok := resourceAccess[clientID].(map[string]interface{}); ok {
		if roles, ok := clientResources["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					user.Roles = append(user.Roles, roleStr)
				}
			}
		}
	}

	user.Scopes = strings.Split(stringClaim(claims, "scope"), " ")
	user.IsAdmin = user.HasRole("admin")

	return user, nil
}

func stringClaim(claims map[string]interface{}, key string) string {
	if value, ok := claims[key].(string); ok {
		return value
	}
	return ""
}
