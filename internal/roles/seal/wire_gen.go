// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package seal

import (
	"github.com/dapperlabs/bamboo-node/internal/roles/seal/config"
)

// Injectors from wire.go:

func InitializeServer() (*Server, error) {
	configConfig := config.New()
	controller := NewController()
	server, err := NewServer(configConfig, controller)
	if err != nil {
		return nil, err
	}
	return server, nil
}
