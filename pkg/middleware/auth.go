package middleware

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/auth0-developer-hub/api_standard-library_golang_hello-world/pkg/helpers"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/pkg/errors"
)

const (
	missingJWTErrorMessage = "Requires authentication"
	invalidJWTErrorMessage = "Bad credentials"
)

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
		)
		if err != nil {
			log.Fatalf("Failed to set up the jwt validator")
		}

		errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Encountered error while validating JWT: %v", err)
			if errors.Is(err, jwtmiddleware.ErrJWTMissing) {
				errorMessage := ErrorMessage{Message: missingJWTErrorMessage}
				helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage)
				return
			}
			if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
				errorMessage := ErrorMessage{Message: invalidJWTErrorMessage}
				helpers.WriteJSON(w, http.StatusUnauthorized, errorMessage)
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
