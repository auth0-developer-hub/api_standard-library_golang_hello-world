package config

import (
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

func SecureOptions() secure.Options {
	return secure.Options{
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ForceSTSHeader:        true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		CustomBrowserXssValue: "0",
		ContentSecurityPolicy: "default-src 'self', frame-ancestors 'none'",
	}
}

func CorsOptions(clientOriginUrl string) cors.Options {
	return cors.Options{
		AllowedOrigins: []string{clientOriginUrl},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:         86400,
	}
}
