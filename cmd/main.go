package main

import (
	"infotecs/pkg/config"
	"infotecs/pkg/server"
	"log"
)

func main() {
	log.Println("Initializing server...")
	cfg := config.ServerConfig{
		Host:   "localhost",
		Port:   "8081",
		DbName: "postgres",
		DbHost: "localhost",
		DbPort: "5433",
		DbUser: "kl",
		DbPass: "password",
		Cert:   "localhost.crt",
		Key:    "localhost.key",
	}
	err := server.Start(cfg)
	if err != nil {
		log.Println("Error - Starting server:", err.Error())
	}
}
