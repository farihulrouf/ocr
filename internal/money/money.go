package money

// Money merepresentasikan nilai uang dalam smallest unit.
//
// Contoh:
//
//	JPY ¥500
//	Amount = 500
//
//	HKD 454.30
//	Amount = 45430
//
//	USD 12.99
//	Amount = 1299
type Money struct {
	Currency Currency `json:"currency"`

	// Amount disimpan dalam smallest unit.
	//
	// USD 12.50 -> 1250
	// HKD 99.90 -> 9990
	// JPY 500 -> 500
	Amount int64 `json:"amount"`
}

// IsZero mengembalikan true jika amount == 0.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsNegative mengembalikan true jika amount < 0.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// IsPositive mengembalikan true jika amount > 0.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// Clone membuat salinan Money.
func (m Money) Clone() Money {
	return Money{
		Currency: m.Currency,
		Amount:   m.Amount,
	}
}

// Equal membandingkan dua Money.
func (m Money) Equal(other Money) bool {
	return m.Currency.Code == other.Currency.Code &&
		m.Amount == other.Amount
}
