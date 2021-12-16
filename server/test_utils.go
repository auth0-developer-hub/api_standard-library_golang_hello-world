package server

import (
	"encoding/json"
	"fmt"
	"hello-golang-api/common"
	"io/ioutil"
	"net/http"
	"strings"
)

var settings = map[string]string{
	//server settings
	"server.host": "",
	"server.port": "4040",

	// Logger Defaults
	"logger.level":            "info",
	"logger.dev_mode":         "false",
	"logger.disable_caller":   "false",
	"logger.outputpath":       "logs/oauth_test.log",
	"logger.error_outputpath": "logs/internal_error_test.log",

	//profiler settings
	"profiler.enabled": "false",
	"profiler.host":    "localhost",
	"profiler.port":    "6060",
	"profiler_path":    "/debug",

	//prometheus settings
	"prometheus_enabled": "false",
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
	"database.file":            "helloworld_test.db",
	"database.auto_create":     "true",
	"database.max_connections": "40",
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

func GetToken(cfg *common.Config) (*Token, error) {

	url := fmt.Sprintf("https://%s/oauth/token", cfg.Domain)

	tokReq := &TokenRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Audience:     cfg.Audience,
		GrantType:    "client_credentials",
	}
	tokBytes, err := json.Marshal(tokReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(tokBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := &Token{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}
