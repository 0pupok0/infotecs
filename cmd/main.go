package main

import (
	"infotecs/pkg/config"
	"infotecs/pkg/server"
	"log"
)

func main() {
	log.Println("Initializing server...")
	// Конфигурация сервера. Данные указаны для примера
	cfg := config.ServerConfig{
		// Хост и порт на котором будет запущен сервер
		Host: "localhost",
		Port: "8081",
		// Данные для подключения к БД
		DbName: "postgres",
		DbHost: "localhost",
		DbPort: "5433",
		DbUser: "kl",
		DbPass: "password",
		// Сертификат SSL для HTTPS подключения (оставить пустыми для HTTP)
		Cert: "localhost.crt",
		Key:  "localhost.key",
	}
	err := server.Start(cfg)
	if err != nil {
		log.Println("Error - Starting server:", err.Error())
	}
}
