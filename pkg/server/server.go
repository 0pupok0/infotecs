package server

import (
	"infotecs/pkg/api"
	"infotecs/pkg/config"
)

func Start(serverConfig config.ServerConfig) error {
	apiKeeper, err := api.New(serverConfig)
	if err != nil {
		return err
	}

	err = apiKeeper.Setup()
	if err != nil {
		return err
	}
	return apiKeeper.Serve()
}
