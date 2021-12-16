package server

import (
	"context"
	"errors"
	"fmt"
	"hello-golang-api/common"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

// validateToken middleware verifies a valid Auth0 JWT token being present in the request.
func (s *Server) validateToken(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		token, err := extractToken(req, s.authConfig.Audience, s.tenantKeys)
		if err != nil {
			s.logger.Errorf("failed to parshasPermissione payload: %s\n", err)
			rw.WriteHeader(http.StatusUnauthorized)
			sendMessage(rw, &message{err.Error()})
			return
		}
		ctxWithToken := context.WithValue(req.Context(), common.CtxTokenKey, token)
		next.ServeHTTP(rw, req.WithContext(ctxWithToken))
	})
}

// extractToken parses the Authorization HTTP header for valid JWT token and
// validates it with AUTH0 JWK keys. Also verifies if the audience present in
// the token matches with the designated audience as per current configuration.
func extractToken(req *http.Request, auth0Audience string, tenantKeys jwk.Set) (jwt.Token, error) {
	authorization := req.Header.Get(common.AuthHeader)
	if authorization == "" {
		return nil, errors.New("authorization header missing")
	}
	bearerAndToken := strings.Split(authorization, " ")
	if len(bearerAndToken) < 2 {
		return nil, errors.New("malformed authorization header: " + authorization)
	}
	token, err := jwt.Parse([]byte(bearerAndToken[1]), jwt.WithKeySet(tenantKeys),
		jwt.WithValidate(true), jwt.WithAudience(auth0Audience))
	if err != nil {
		return nil, err
	}
	return token, nil
}

func hasPermission(next http.Handler, permission string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		token := req.Context().Value(common.CtxTokenKey).(jwt.Token)
		if token == nil {
			fmt.Printf("failed to find token in context\n")
			rw.WriteHeader(http.StatusForbidden)
			sendMessage(rw, &message{http.StatusText(http.StatusForbidden)})
			return
		}
		if !tokenHasPermission(token, permission) {
			fmt.Printf("permission check failed\n")
			rw.WriteHeader(http.StatusForbidden)
			sendMessage(rw, &message{http.StatusText(http.StatusForbidden)})
			return
		}
		next.ServeHTTP(rw, req)
	})
}

func tokenHasPermission(token jwt.Token, permission string) bool {
	claims := token.PrivateClaims()
	tkPermissions, ok := claims[common.PermClaim]
	if !ok {
		return false
	}
	tkPermList, ok := tkPermissions.([]interface{})
	if !ok {
		return false
	}
	for _, perm := range tkPermList {
		if perm == permission {
			return true
		}
	}
	return false
}

// fetchTenantKeys fetch and parse the tenant JSON Web Keys (JWK). The keys
// are used for JWT token validation during requests authorization.
func fetchTenantKeys(auth0Domain string) (jwk.Set, error) {
	set, err := jwk.Fetch(context.Background(),
		fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain))
	if err != nil {
		return nil, fmt.Errorf("failed to parse tenant json web keys: %s\n", err)
	}
	return set, nil

}
