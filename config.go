package gosparkclient

import (
	"errors"
	"time"
)

const (
	defaultTimeout  = 30 * time.Second
	defaultUID      = "12345"
	defaultAuditing = "default"
)

// Config holds all configuration options for SparkClient
type Config struct {
	AppID     string
	ApiSecret string
	ApiKey    string
	HostURL   string
	EMBURL    string
	Domain    string
	Timeout   time.Duration
	UID       string
	Auditing  string
}

// ConfigOption defines a function type for setting config options
type ConfigOption func(*Config)

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Timeout:  defaultTimeout,
		UID:      defaultUID,
		Auditing: defaultAuditing,
	}
}

// validateConfig checks if the configuration is valid
func validateConfig(c *Config) error {
	if c.AppID == "" {
		return errors.New("AppID is required")
	}
	if c.ApiSecret == "" {
		return errors.New("ApiSecret is required")
	}
	if c.ApiKey == "" {
		return errors.New("ApiKey is required")
	}
	if c.HostURL == "" {
		return errors.New("HostURL is required")
	}
	return nil
}

// WithTimeout sets the timeout for the client
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithUID sets the UID for the client
func WithUID(uid string) ConfigOption {
	return func(c *Config) {
		c.UID = uid
	}
}

// WithAuditing sets the auditing mode
func WithAuditing(auditing string) ConfigOption {
	return func(c *Config) {
		c.Auditing = auditing
	}
}

// WithCredentials sets the credentials for the client
func WithCredentials(appID, apiKey, apiSecret string) ConfigOption {
	return func(c *Config) {
		c.AppID = appID
		c.ApiKey = apiKey
		c.ApiSecret = apiSecret
	}
}

// WithURLs sets the URLs for the client
func WithURLs(hostURL, embURL string) ConfigOption {
	return func(c *Config) {
		c.HostURL = hostURL
		c.EMBURL = embURL
	}
}

// WithDomain sets the domain for the client
func WithDomain(domain string) ConfigOption {
	return func(c *Config) {
		c.Domain = domain
	}
}

// WithConfig sets the entire configuration
func WithConfig(config *Config) ConfigOption {
	return func(c *Config) {
		*c = *config
	}
}
