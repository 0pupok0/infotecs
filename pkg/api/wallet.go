package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"infotecs/pkg/models"
	"log"
	"net/http"
	"strconv"
)

// WalletCreateHandler - метод API для создания нового кошелька с 100 у.е. на балансе
// Возвращает структуру нового кошелька
func (api *API) WalletCreateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		wallet, err := api.db.CreateWallet()
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// WalletGetHandler - метод API для получения данных о кошельке
//
// Прнимает walletID (в динамическом маршруте) - идентификационный номер кошелька
// Возвращает структуру найденного кошелька
//
// Коды ошибок:
// 400 - если не был передан walletID или он не является чилом
// 404 - если кошелёк с таким id не был найден
func (api *API) WalletGetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := mux.Vars(r)["walletID"]
		if idStr == "" {
			http.Error(w, "Expected wallet id", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Could not convert id to integer", http.StatusBadRequest)
			return
		}
		wallet, err := api.db.GetWalletByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(wallet)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// SendHandler - метод API для перевода денежных средств с одного кошелька на другой
//
// Принимает:
// walletID (в динамическом маршруте) - номер кошелька отправителя
// To (в теле запроса) - номер кошелька получателя
// Amount (в теле запроса) - сумма переводимых денежных средств
//
// Коды ошибок:
// 404 - Исходящий кошелек не найден
// 400 - Целевой кошелек не найден или на исходящем нет нужной суммы
func (api *API) SendHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		senderIDStr := mux.Vars(r)["walletID"]
		if senderIDStr == "" {
			http.Error(w, "Expected wallet id", http.StatusBadRequest)
			return
		}
		var transaction models.Transaction
		err := json.NewDecoder(r.Body).Decode(&transaction)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		transaction.From, err = strconv.Atoi(senderIDStr)
		if err != nil {
			http.Error(w, "Could not convert wallet id to integer", http.StatusBadRequest)
			return
		}
		senderWallet, err := api.db.GetWalletByID(transaction.From)
		if err != nil {
			http.Error(w, "Could not find sender wallet", http.StatusNotFound)
			return
		}
		if senderWallet.Balance < transaction.Amount {
			http.Error(w, "Not enough money to send", http.StatusBadRequest)
			return
		}
		receiverWallet, err := api.db.GetWalletByID(transaction.To)
		if err != nil {
			http.Error(w, "Could not find receiver wallet", http.StatusBadRequest)
			return
		}
		err = api.db.SubmitTransaction(senderWallet, receiverWallet, transaction.Amount)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// HistoryHandler - метод API для получения истории входящих и исходящих транзакций
// Принимает walletID (в динамическом маршруте) - номер кошелька
// Возвращает массив транзакций связанных с этим кошельком
//
// Коды ошибок:
// 400 - если не был передан walletID или он не является чилом
// 404 - если кошелёк с таким id не был найден
func (api *API) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := mux.Vars(r)["walletID"]
		if idStr == "" {
			http.Error(w, "Expected wallet id", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Could not convert wallet id to integer", http.StatusBadRequest)
			return
		}
		history, err := api.db.GetHistory(id)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(history)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
