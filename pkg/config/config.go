package config

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// AutotaskConfig represents the configuration for the Autotask datasource
type AutotaskConfig struct {
	Username        string `json:"username"`
	URL             string `json:"url"`
	Secret          string `json:"-"`
	IntegrationCode string `json:"-"`
}

// LoadSettings loads the configuration from Grafana's datasource settings
func LoadSettings(settings backend.DataSourceInstanceSettings) (*AutotaskConfig, error) {
	config := &AutotaskConfig{}

	if err := json.Unmarshal(settings.JSONData, config); err != nil {
		return nil, err
	}

	// Load secure settings from Grafana's encrypted store
	config.Secret = settings.DecryptedSecureJSONData["secret"]
	config.IntegrationCode = settings.DecryptedSecureJSONData["integrationCode"]

	// Fall back to datasource URL if not in JSONData
	if config.URL == "" {
		config.URL = settings.URL
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *AutotaskConfig) Validate() error {
	if c.Username == "" {
		return ErrMissingUsername
	}
	if c.Secret == "" {
		return ErrMissingSecret
	}
	if c.IntegrationCode == "" {
		return ErrMissingIntegrationCode
	}
	if c.URL == "" {
		return ErrMissingURL
	}
	return nil
}
