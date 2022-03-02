package auth

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v10"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Provider interface {
	// Scopes decodes a list of scopes from the token.
	Scopes(token string) ([]Scope, error)
}

//
// Scope represents an authorization scope.
type Scope interface {
	// Check whether the scope gives access to the resource with the method.
	Check(resource string, method string) bool
}

//
// NoAuth provider always permits access.
type NoAuth struct{}

//
// Scopes decodes a list of scopes from the token.
// For the NoAuth provider, this just returns a single instance
// of the NoAuthScope.
func (r *NoAuth) Scopes(token string) (scopes []Scope, err error) {
	scopes = append(scopes, &NoAuthScope{})
	return
}

//
// NoAuthScope always permits access.
type NoAuthScope struct{}

//
// Check whether the scope gives access to the resource with the method.
func (r *NoAuthScope) Check(_ string, _ string) (ok bool) {
	ok = true
	return
}

//
// NewKeycloak builds a new Keycloak auth provider.
func NewKeycloak(host, realm, id, secret string) (k Keycloak) {
	client := gocloak.NewClient(host)
	k = Keycloak{
		host:   host,
		realm:  realm,
		id:     id,
		secret: secret,
		client: client,
	}
	return
}

//
// Keycloak auth provider
type Keycloak struct {
	client gocloak.GoCloak
	host   string
	realm  string
	id     string
	secret string
}

//
// Scopes decodes a list of scopes from the token.
func (r *Keycloak) Scopes(token string) (scopes []Scope, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	decoded, _, err := r.client.DecodeAccessToken(ctx, token, r.realm)
	if err != nil {
		err = errors.New("invalid token")
		return
	}
	if !decoded.Valid {
		err = errors.New("invalid token")
		return
	}
	claims, ok := decoded.Claims.(*jwt.MapClaims)
	if !ok || claims == nil {
		err = errors.New("invalid token")
		return
	}
	rawClaimScopes, ok := (*claims)["scope"].(string)
	if !ok {
		err = errors.New("invalid token")
		return
	}
	claimScopes := strings.Split(rawClaimScopes, " ")
	for _, s := range claimScopes {
		scope := NewKeycloakScope(s)
		scopes = append(scopes, &scope)
	}
	return
}

//
// NewKeycloakScope builds a Scope object from a string.
func NewKeycloakScope(s string) (scope KeycloakScope) {
	if strings.Contains(s, ":") {
		segments := strings.Split(s, ":")
		scope.resource = segments[0]
		scope.method = segments[1]
	} else {
		scope.resource = s
	}
	return
}

//
// KeycloakScope is a scope decoded from a Keycloak token.
type KeycloakScope struct {
	resource string
	method   string
}

//
// Check whether the scope gives access to the resource with the method.
func (r *KeycloakScope) Check(resource string, method string) (ok bool) {
	ok = r.resource == resource && r.method == method
	return
}
