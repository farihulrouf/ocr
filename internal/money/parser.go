package money

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Parse mengubah string menjadi Money.
//
// Contoh:
//
//	"HK$454.30"
//	"$10.50"
//	"¥540"
//	"Rp15.000"
//	"EUR 12.99"
//	"(HK$54.30)"
//	"-¥300"
func Parse(text string) (Money, error) {

	text = strings.TrimSpace(text)

	if text == "" {
		return Money{}, fmt.Errorf("empty money string")
	}

	// -----------------------------
	// Detect Currency
	// -----------------------------

	currency, ok := DetectCurrency(text)

	if !ok {
		return Money{}, fmt.Errorf("unsupported currency")
	}

	// -----------------------------
	// Remove Currency Symbol
	// -----------------------------

	value := text

	value = strings.ReplaceAll(value, currency.Symbol, "")
	value = strings.ReplaceAll(value, currency.Code, "")

	value = strings.TrimSpace(value)

	negative := false

	// -----------------------------
	// Negative Number
	// -----------------------------

	if strings.HasPrefix(value, "-") {

		negative = true

		value = strings.TrimPrefix(value, "-")
	}

	if strings.HasPrefix(value, "(") &&
		strings.HasSuffix(value, ")") {

		negative = true

		value = strings.TrimPrefix(value, "(")
		value = strings.TrimSuffix(value, ")")
	}

	value = strings.TrimSpace(value)

	// -----------------------------
	// Parse Number
	// -----------------------------

	amount, err := parseNumber(value, currency.Scale)

	if err != nil {
		return Money{}, err
	}

	if negative {
		amount = -amount
	}

	return Money{
		Currency: currency,
		Amount:   amount,
	}, nil
}

// parseNumber
//
// Mengubah:
//
//	454.30
//	15,000
//	15.000
//	1,999.50
//
// menjadi smallest unit.
func parseNumber(text string, scale int) (int64, error) {

	text = strings.TrimSpace(text)

	if text == "" {
		return 0, fmt.Errorf("empty amount")
	}

	// -----------------------------------
	// Currency tanpa decimal
	// -----------------------------------

	if scale == 0 {

		text = strings.ReplaceAll(text, ",", "")
		text = strings.ReplaceAll(text, ".", "")

		value, err := strconv.ParseInt(text, 10, 64)

		if err != nil {
			return 0, err
		}

		return value, nil
	}

	// -----------------------------------
	// Currency decimal
	// -----------------------------------

	// 1,999.50
	// 1999.50

	if strings.Count(text, ".") > 1 {
		return 0, fmt.Errorf("invalid decimal")
	}

	if strings.Contains(text, ".") {

		text = strings.ReplaceAll(text, ",", "")

		f, err := strconv.ParseFloat(text, 64)

		if err != nil {
			return 0, err
		}

		return int64(math.Round(f * math.Pow10(scale))), nil
	}

	// 1,999

	text = strings.ReplaceAll(text, ",", "")

	v, err := strconv.ParseInt(text, 10, 64)

	if err != nil {
		return 0, err
	}

	return v * int64(math.Pow10(scale)), nil
}

// MustParse
//
// Panic jika parsing gagal.
//
// Cocok untuk unit test.
func MustParse(text string) Money {

	m, err := Parse(text)

	if err != nil {
		panic(err)
	}

	return m
}

// ParseAmount
//
// Shortcut jika hanya membutuhkan Amount.
func ParseAmount(text string) (int64, error) {

	m, err := Parse(text)

	if err != nil {
		return 0, err
	}

	return m.Amount, nil
}

// ParseCurrency
//
// Shortcut jika hanya membutuhkan Currency.
func ParseCurrency(text string) (Currency, error) {

	m, err := Parse(text)

	if err != nil {
		return Currency{}, err
	}

	return m.Currency, nil
}
