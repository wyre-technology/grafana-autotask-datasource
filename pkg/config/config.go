package config

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// AutotaskConfig represents the configuration for the Autotask datasource
type AutotaskConfig struct {
	Username        string `json:"username"`
	Zone            string `json:"zone"`
	Secret          string `json:"-"`
	IntegrationCode string `json:"-"`
}

// LoadSettings loads the configuration from Grafana's datasource settings
func LoadSettings(settings backend.DataSourceInstanceSettings) (*AutotaskConfig, error) {
	config := &AutotaskConfig{}

	// Load basic settings
	if err := json.Unmarshal(settings.JSONData, config); err != nil {
		return nil, err
	}

	// Load secure settings
	config.Secret = settings.DecryptedSecureJSONData["secret"]
	config.IntegrationCode = settings.DecryptedSecureJSONData["integrationCode"]

	return config, nil
}

// Validate checks if the configuration is valid
func (c *AutotaskConfig) Validate() error {
	if c.Username == "" {
		return ErrMissingUsername
	}
	if c.Zone == "" {
		return ErrMissingZone
	}
	if c.Secret == "" {
		return ErrMissingSecret
	}
	if c.IntegrationCode == "" {
		return ErrMissingIntegrationCode
	}
	return nil
}
