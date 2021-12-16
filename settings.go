package main

import (
	_ "embed"
	"fmt"
	"net/http"
)

//go:embed dev.yaml
var Env []byte

var Settings = map[string]string{
	//server settings
	"server.host": "",
	"server.port": "4040",

	// Logger Defaults
	"logger.level":            "info",
	"logger.dev_mode":         "true",
	"logger.disable_caller":   "false",
	"logger.outputpath":       "logs/oauth.log",
	"logger.error_outputpath": "logs/internal_error.log",

	//profiler settings
	"profiler.enabled": "true",
	"profiler.host":    "localhost",
	"profiler.port":    "6060",
	"profiler_path":    "/debug",

	//prometheus settings
	"prometheus_enabled": "true",
	"prometheus_path":    "/prometheus",

	//request settings
	"httprate_limit":      "10",
	"httprate_limit_time": "10s",

	//cors settings
	"cros.allowed_origins": "http://localhost:4040",
	"cors.allowed_methods": fmt.Sprintf("%s,%s,%s,%s,%s,%s", http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete),
	"cors.allowed_headers": "Authorization,Content-Type",
	"cors.max_age":         "86400",

	// Database Settings
	"database.file":            "helloworld.db",
	"database.auto_create":     "true",
	"database.max_connections": "40",
}
