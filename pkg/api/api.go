package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"infotecs/pkg/config"
	"infotecs/pkg/storage"
	"net/http"
)

type API struct {
	router        *mux.Router
	db            *storage.Storage
	cert          string
	key           string
	serverAddress string
}

// New - метод создания экземпляра API
func New(cfg config.ServerConfig) (*API, error) {
	db, err := storage.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbName,
	))
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter()

	return &API{
		router:        router,
		db:            db,
		cert:          cfg.Cert,
		key:           cfg.Key,
		serverAddress: fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
	}, nil
}

// Setup - назначение хэндлеров для эндпоинтов
func (api *API) Setup() error {
	err := api.db.CreateTables()
	if err != nil {
		return err
	}
	api.router.HandleFunc("/api/v1/ping", api.PingHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/api/v1/wallet", api.WalletCreateHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/wallet/{walletID}", api.WalletGetHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/api/v1/wallet/{walletID}/send", api.SendHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.HandleFunc("/api/v1/wallet/{walletID}/history", api.HistoryHandler).Methods(http.MethodGet, http.MethodOptions)
	return nil
}

// Serve - запуск API
func (api *API) Serve() error {
	if api.cert == "" || api.key == "" {
		return http.ListenAndServe(api.serverAddress, api.router)
	}
	return http.ListenAndServeTLS(api.serverAddress, api.cert, api.key, api.router)
}
