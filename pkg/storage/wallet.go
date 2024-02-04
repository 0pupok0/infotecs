package storage

import (
	"context"
	"infotecs/pkg/models"
	"time"
)

func (s *Storage) CreateWallet() (models.Wallet, error) {
	var id int
	err := s.pool.QueryRow(context.Background(), `INSERT INTO wallets DEFAULT VALUES RETURNING id`).
		Scan(&id)
	return models.Wallet{ID: id, Balance: 100}, err
}

func (s *Storage) GetWalletByID(id int) (models.Wallet, error) {
	var wallet models.Wallet
	err := s.pool.QueryRow(context.Background(),
		`SELECT id, balance FROM wallets WHERE id = $1`, id).
		Scan(&wallet.ID, &wallet.Balance)
	return wallet, err
}

func (s *Storage) SubmitTransaction(sender, receiver models.Wallet, amount float64) error {
	_, err := s.pool.Query(context.Background(), `BEGIN TRANSACTION`)
	if err != nil {
		return err
	}
	_, err = s.pool.Query(context.Background(),
		`UPDATE wallets SET balance = $2 WHERE id = $1`, sender.ID, sender.Balance-amount)
	if err != nil {
		_, err2 := s.pool.Query(context.Background(), `ROLLBACK TRANSACTION`)
		if err2 != nil {
			return err2
		}
		return err
	}
	_, err = s.pool.Query(context.Background(),
		`UPDATE wallets SET balance = $2 WHERE id = $1`, receiver.ID, receiver.Balance+amount)
	if err != nil {
		_, err2 := s.pool.Query(context.Background(), `ROLLBACK TRANSACTION`)
		if err2 != nil {
			return err2
		}
		return err
	}
	_, err = s.pool.Query(context.Background(),
		`INSERT INTO transactions_history(time, from_id, to_id, amount) VALUES ($1, $2, $3, $4)`,
		time.Now().UTC().Format(time.RFC3339), sender.ID, receiver.ID, amount)
	if err != nil {
		_, err2 := s.pool.Query(context.Background(), `ROLLBACK TRANSACTION`)
		if err2 != nil {
			return err2
		}
		return err
	}
	_, err = s.pool.Query(context.Background(), `COMMIT TRANSACTION`)
	return err
}

func (s *Storage) GetHistory(id int) ([]models.Transaction, error) {
	var res []models.Transaction
	rows, err := s.pool.Query(context.Background(),
		`SELECT to_char(time, 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"'), from_id, to_id, amount FROM transactions_history WHERE from_id = $1 OR to_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err = rows.Scan(&transaction.Time, &transaction.From, &transaction.To, &transaction.Amount)
		if err != nil {
			return nil, err
		}
		res = append(res, transaction)
	}
	return res, nil
}
