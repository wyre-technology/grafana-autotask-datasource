package config

import "errors"

var (
	ErrMissingUsername        = errors.New("username is required")
	ErrMissingZone            = errors.New("zone is required")
	ErrMissingSecret          = errors.New("secret is required")
	ErrMissingIntegrationCode = errors.New("integration code is required")
)
