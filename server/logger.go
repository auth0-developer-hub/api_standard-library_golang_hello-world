package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Returns a middleware function for logging requests
func loggerHTTPMiddlewareStackdriver(logRequestBody bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//start timer
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				ww.Tee(response)
			}

			//serve HTTP
			next.ServeHTTP(ww, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("X-Forwarded-For")
			if remoteIP == "" {
				remoteIP = r.RemoteAddr
			}

			fields := []zapcore.Field{
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: r.Method,
					RequestURL:    r.RequestURI,
					RequestSize:   strconv.FormatInt(r.ContentLength, 10),
					Status:        ww.Status(),
					ResponseSize:  strconv.Itoa(ww.BytesWritten()),
					UserAgent:     r.UserAgent(),
					RemoteIP:      remoteIP,
					Referer:       r.Referer(),
					Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
					Protocol:      r.Proto,
				}),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}

			zap.L().Info("HTTP Request", fields...)
		})
	}
}
