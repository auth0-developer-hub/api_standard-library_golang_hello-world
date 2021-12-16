package server

import (
	"encoding/json"
	"hello-golang-api/common"
	"hello-golang-api/embed"
	"hello-golang-api/store"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	logger     *zap.SugaredLogger
	authConfig *common.Config
	tenantKeys jwk.Set
	settings   map[string]string
	store      store.MessageStore
}

type message struct {
	Message string `json:"message"`
}

var (
	publicMessage    = &message{"The API doesn't require an access token to share this message."}
	protectedMessage = &message{"The API successfully validated your access token."}
	adminMessage     = &message{"The API successfully recognized you as an admin."}
)

func publicApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, publicMessage)
}

func pingHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=UTF-8")
	rw.Header().Add("X-Frame-Options", "deny")
	rw.Header().Add("X-Content-Type-Options", "nosniff")
	_, err := rw.Write([]byte("pong"))
	if err != nil {
		log.Print("http response write error", err)
	}
}

func protectedApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, protectedMessage)
}

func adminApiHandler(rw http.ResponseWriter, _ *http.Request) {
	sendMessage(rw, adminMessage)
}

func sendMessage(rw http.ResponseWriter, data interface{}) {
	rw.Header().Add("Content-Type", "application/json")
	rw.Header().Add("X-Frame-Options", "deny")
	rw.Header().Add("X-Content-Type-Options", "nosniff")

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Print("json conversion error", err)
		return
	}
	_, err = rw.Write(bytes)
	if err != nil {
		log.Print("http response write error", err)
	}
}

func New(cfg *common.Config, settings map[string]string, dbstore store.MessageStore, logger *zap.SugaredLogger) (*Server, error) {
	tenantKeys, err := fetchTenantKeys(cfg.Domain)
	if err != nil {
		return nil, err
	}
	s := &Server{
		authConfig: cfg,
		tenantKeys: tenantKeys,
		settings:   settings,
		store:      dbstore,
		logger:     logger,
	}
	return s, nil
}

func (s *Server) Handler() http.Handler {

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	// Log Requests - Use appropriate format depending on the encoding
	r.Use(loggerHTTPMiddlewareStackdriver(s.settings["log_requests_body"] == "true"))
	//Add prometheus middleware
	if s.settings["server.prometheus_enabled"] == "true" {
		r.Use(PrometheusMiddleware)
	}
	rateTime, err := time.ParseDuration(s.settings["httprate_limit_time"])
	if err != nil {
		rateTime = time.Second * 10
	}
	limit, err := strconv.Atoi(s.settings["server.httprate_limit"])
	if err != nil {
		limit = 10
	}
	// Enable httprate request limiter
	r.Use(httprate.Limit(
		limit,    // requests
		rateTime, // per duration
		//Rate limit by IP and URL path
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	maxAge, err := strconv.Atoi(s.settings["cors.max_age"])
	if err != nil {
		maxAge = 86400
	}
	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins: strings.Split(s.settings["cros.allowed_origins"], ","),
		AllowedMethods: strings.Split(s.settings["cors.allowed_methods"], ","),
		AllowedHeaders: strings.Split(s.settings["cors.allowed_headers"], ","),
		MaxAge:         maxAge,
	}).Handler)

	if s.settings["prometheus_enabled"] == "true" {
		// Prometheus endpoint
		r.Mount(s.settings["prometheus_path"], promhttp.Handler())
	}

	// Enable profiler
	if s.settings["profiler_enabled"] == "true" {
		zap.S().Debugw("Profiler enabled on API", "path", s.settings["profiler_path"])
		r.Mount(s.settings["profiler_path"], middleware.Profiler())
	}

	r.Route("/api", func(api chi.Router) {
		//sub routes
		api.Route("/messages", func(message chi.Router) {
			//exisiting routes
			message.Get("/public", publicApiHandler)
			message.Get("/protected", s.validateToken(http.HandlerFunc(protectedApiHandler)))
			message.Get("/admin", s.validateToken(hasPermission(http.HandlerFunc(adminApiHandler), "read:admin-messages")))

			message.Post("/", s.validateToken(http.HandlerFunc(s.MessageSave())))
			message.Get("/", s.MessageFind())
			//sub routes
			message.Route("/{id}", func(w chi.Router) {
				//user routes
				w.Put("/", s.validateToken(http.HandlerFunc(s.MessageUpdate())))
				w.Get("/", s.MessageGetByID())
				//only admin has rights
				w.Delete("/", s.validateToken(hasPermission(http.HandlerFunc(s.MessageDeleteByID()), "read:admin-messages")))
			})
		})

	})

	r.Get("/", http.NotFound)
	r.Get("/ping", http.HandlerFunc(pingHandler))
	r.Get("/version", common.GetVersion())

	// Serve api_docs and swagger-ui
	docsFileServer := http.FileServer(http.FS(embed.PublicHTMLFS()))
	r.Mount("/api_docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "no-cache")
		docsFileServer.ServeHTTP(w, r)
	}))
	return r
}

func DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}
