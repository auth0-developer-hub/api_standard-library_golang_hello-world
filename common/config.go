package common

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Audience     string
	Domain       string
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

func (c *Config) Validate() error {
	if c.Audience == "" {
		return errors.New("auth0 API identifier (as audience) missing")
	}
	if c.Domain == "" {
		return errors.New("auth0 API tenant domain missing")
	}
	if c.ClientID == "" {
		return errors.New("auth0 API Client-iD missing")
	}
	if c.ClientSecret == "" {
		return errors.New("auth0 API Client-Secret missing")
	}
	if c.CallbackURL == "" {
		return errors.New("auth0 API Callback-URL missing")
	}
	return nil
}

func (c *Config) Overrides(config *Config) {
	if config == nil {
		return
	}

	if config.Audience != "" {
		c.Audience = config.Audience
	}
	if config.Domain != "" {
		c.Domain = config.Domain
	}
	if config.ClientID != "" {
		c.ClientID = config.ClientID
	}
	if config.ClientSecret != "" {
		c.ClientSecret = config.ClientSecret
	}
	if config.CallbackURL != "" {
		c.CallbackURL = config.CallbackURL
	}
}

func loadAuth0YAML(env []byte) (*Config, error) {
	cfg := map[string]string{}
	err := yaml.Unmarshal(env, cfg)
	if err != nil {
		return nil, fmt.Errorf("yaml unmarshal error: %v", err)
	}
	authConfig := &Config{}
	if val, found := cfg[domainYamlKey]; found {
		authConfig.Domain = val
	}
	if val, found := cfg[audienceYamlKey]; found {
		authConfig.Audience = val
	}
	if val, found := cfg[clientIDYamlKey]; found {
		authConfig.ClientID = val
	}
	if val, found := cfg[clientSecretYamlKey]; found {
		authConfig.ClientSecret = val
	}
	if val, found := cfg[callbackURLYamlKey]; found {
		authConfig.CallbackURL = val
	}
	return authConfig, nil
}

func parseArgs() *Config {
	authConfig := &Config{}
	flag.StringVar(&authConfig.Audience, "a",
		os.Getenv("AUTH0_AUDIENCE"), "Auth0 API identifier, as audience")
	flag.StringVar(&authConfig.Domain, "d",
		os.Getenv("AUTH0_DOMAIN"), "Auth0 API tenant domain")
	flag.StringVar(&authConfig.ClientID, "i",
		os.Getenv("AUTH0_CLIENTID"), "Auth0 API client-ID")
	flag.StringVar(&authConfig.ClientSecret, "p",
		os.Getenv("AUTH0_CLIENTSECRET"), "Auth0 API client-Secret")
	flag.StringVar(&authConfig.CallbackURL, "c",
		os.Getenv("AUTH0_CALLBACK_URL"), "Auth0 API callback-URL")

	flag.Parse()
	return authConfig
}

func InitConfig(env []byte) (*Config, error) {
	//load cfg from YAML file
	config, err := loadAuth0YAML(env)
	if err != nil {
		return nil, err
	}
	//get cfg from Arguments and Environment variables and Overrides YAML cfg values
	cfg := parseArgs()
	config.Overrides(cfg)

	//validate the final config
	err = config.Validate()
	if err != nil {
		return nil, err
	}
	return config, nil
}
