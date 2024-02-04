package server

import (
	"infotecs/pkg/api"
	"infotecs/pkg/config"
)

// Start - функция для запуска сервера
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
