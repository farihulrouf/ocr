package money

import (
	"errors"
	"fmt"
)

// Error yang digunakan validator.
var (
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrUnsupported     = errors.New("unsupported currency")
)

// Validate memastikan Money valid.
func Validate(m Money) error {

	// Currency wajib ada
	if m.Currency.Code == "" {
		return ErrInvalidCurrency
	}

	// Currency harus didukung
	if !IsSupportedCurrency(m.Currency.Code) {
		return fmt.Errorf("%w: %s", ErrUnsupported, m.Currency.Code)
	}

	// Scale tidak boleh negatif
	if m.Currency.Scale < 0 {
		return fmt.Errorf("invalid currency scale")
	}

	return nil
}

// ValidateCurrencyCode
//
// Contoh:
//
//	USD -> nil
//	XXX -> error
func ValidateCurrencyCode(code string) error {

	if code == "" {
		return ErrInvalidCurrency
	}

	if !IsSupportedCurrency(code) {
		return fmt.Errorf("%w: %s", ErrUnsupported, code)
	}

	return nil
}

// ValidateAmount
//
// Saat ini hanya memastikan amount masih dalam batas int64.
// Tempat yang tepat jika nanti ingin menambah business rule.
func ValidateAmount(amount int64) error {

	// int64 selalu valid.
	// Tambahkan rule lain jika diperlukan.

	return nil
}

// MustValidate
//
// Panic jika tidak valid.
// Berguna untuk unit test.
func MustValidate(m Money) {

	if err := Validate(m); err != nil {
		panic(err)
	}
}

// IsZero
func IsZero(m Money) bool {
	return m.Amount == 0
}

// IsPositive
func IsPositive(m Money) bool {
	return m.Amount > 0
}

// IsNegative
func IsNegative(m Money) bool {
	return m.Amount < 0
}

// SameCurrency
//
// Memastikan dua Money memiliki currency yang sama.
func SameCurrency(a, b Money) bool {
	return a.Currency.Code == b.Currency.Code
}

// CanAdd
//
// Sebelum Add() atau Subtract()
// pastikan currency sama.
func CanAdd(a, b Money) error {

	if !SameCurrency(a, b) {
		return fmt.Errorf(
			"currency mismatch: %s != %s",
			a.Currency.Code,
			b.Currency.Code,
		)
	}

	return nil
}
