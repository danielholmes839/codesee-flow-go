// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package verify

import (
	"github.com/dapperlabs/bamboo-node/internal/roles/verify/config"
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
