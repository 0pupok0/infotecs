package config

type ServerConfig struct {
	// Host - хост на котором будет размещён сервер
	Host string
	// Port - порт на котором будет размещён сервер
	Port string
	// DbName - имя БД к которой будет подключаться сервер
	DbName string
	// DbHost - хост БД к которой будет подключаться сервер
	DbHost string
	// DbPort - порт БД к которой будет подключаться сервер
	DbPort string
	// DbUser - пользователь БД под которым будет авторизован сервер
	DbUser string
	// DbPass - пароль для авторизации под пользователем БД
	DbPass string
	// Cert - файл с сертификатом SSL для HTTPS (пустая строка для HTTP)
	Cert string
	// Cert - файл с ключом SSL для HTTPS (пустая строка для HTTP)
	Key string
}
