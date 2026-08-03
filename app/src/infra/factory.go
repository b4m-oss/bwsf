package infra

import (
	"fmt"

	"bwsf/src/config"
	"bwsf/src/core"
)

// NewBwClientFromConfig selects a BwClient implementation based on cfg.Backend.
// When cfg is nil or Backend is unset, the Bitwarden CLI (`bw`) adapter is used.
func NewBwClientFromConfig(cfg *config.Config) (core.BwClient, error) {
	backend := config.BackendBW
	if cfg != nil {
		backend = cfg.GetBackend()
	}

	switch backend {
	case config.BackendBW:
		return NewBwClient(), nil
	case config.BackendAPI:
		return NewApiBwClient(), nil
	default:
		return nil, fmt.Errorf("unsupported backend %q: use %q or %q", backend, config.BackendBW, config.BackendAPI)
	}
}
