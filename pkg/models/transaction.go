package models

// Transaction - модель данных, описывающая транзакцию
type Transaction struct {
	// From - номер кошелька отправителя
	From int
	// To - номер кошелька получателя
	To int
	// Amount - сумма транзакции
	Amount float64
	// Time - время проыедения транзакции
	Time string
}
