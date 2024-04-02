package middleware

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/helpers"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/pkg/errors"
)

const (
	missingJWTErrorMessage       = "Requires authentication"
	invalidJWTErrorMessage       = "Bad credentials"
	permissionDeniedErrorMessage = "Permission denied"
)

type CustomClaims struct {
	Permissions []string `json:"permissions"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasPermissions(expectedClaims []string) bool {
	if len(expectedClaims) == 0 {
		return false
	}
	for _, scope := range expectedClaims {
		if !helpers.Contains(c.Permissions, scope) {
			return false
		}
	}
	return true
}

func ValidatePermissions(expectedClaims []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasPermissions(expectedClaims) {
			errorMessage := ErrorMessage{Message: permissionDeniedErrorMessage}
			if err := helpers.WriteJSON(w, http.StatusForbidden, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ValidateJWT(audience, domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		issuerURL, err := url.Parse("https://" + domain + "/")
		if err != nil {
			log.Fatalf("Failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{audience},
			validator.WithCustomClaims(func() validator.CustomClaims {
				return new(CustomClaims)
			}),
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}

		if authHeaderParts := strings.Fields(r.Header.Get("Authorization")); len(authHeaderParts) > 0 && strings.ToLower(authHeaderParts[0]) != "bearer" {
			errorMessage := ErrorMessage{Message: invalidJWTErrorMessage}
			if err := helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
				log.Printf("Failed to write error message: %v", err)
			}
			return
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)
			if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
				errorMessage := ErrorMessage{Message: missingJWTErrorMessage}
				if err := helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
				errorMessage := ErrorMessage{Message: invalidJWTErrorMessage}
				if err := helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage); err != nil {
					log.Printf("Failed to write error message: %v", err)
				}
				return
			}
			ServerError(w, err)
		}

		middleware := jwtmiddleware.New(
			jwtValidator.ValidateToken,
			jwtmiddleware.WithErrorHandler(errorHandler),
		)

		middleware.CheckJWT(next).ServeHTTP(w, r)
	})
}
