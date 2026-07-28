// internal/money/detector.go

package money

import "strings"

// symbolPriorities digunakan agar simbol yang lebih panjang
// dicek lebih dahulu.
//
// Contoh:
// HK$ harus dideteksi sebelum $
var symbolPriorities = []string{
	"HK$",
	"NT$",
	"NZ$",
	"A$",
	"C$",
	"S$",
	"CHF",

	"Rp",

	"¥",
	"￥",

	"€",
	"£",
	"$",
	"₩",
	"₫",
	"฿",
}

// symbolCurrencyMap
//
// Mapping simbol -> Currency Code
//
// NOTE:
//
// ¥ diprioritaskan menjadi JPY.
//
// Jika nanti ingin membedakan Jepang vs China,
// bisa dibuat detector berdasarkan bahasa OCR.
var symbolCurrencyMap = map[string]string{

	"HK$": "HKD",
	"NT$": "TWD",
	"NZ$": "NZD",
	"A$":  "AUD",
	"C$":  "CAD",
	"S$":  "SGD",

	"CHF": "CHF",

	"Rp": "IDR",

	"¥": "JPY",
	"￥": "JPY",

	"€": "EUR",
	"£": "GBP",
	"$": "USD",

	"₩": "KRW",
	"₫": "VND",
	"฿": "THB",
}

// DetectCurrency
//
// Mendeteksi currency dari OCR text.
//
// Prioritas:
//
// 1. Symbol
// 2. Currency Code
func DetectCurrency(text string) (Currency, bool) {

	text = strings.TrimSpace(text)

	if text == "" {
		return Currency{}, false
	}

	//-----------------------------------
	// STEP 1
	// Cari berdasarkan simbol
	//-----------------------------------

	for _, symbol := range symbolPriorities {

		if strings.Contains(text, symbol) {

			code := symbolCurrencyMap[symbol]

			currency, ok := FindCurrency(code)

			if ok {
				return currency, true
			}
		}
	}

	//-----------------------------------
	// STEP 2
	// Cari berdasarkan currency code
	//-----------------------------------

	upper := strings.ToUpper(text)

	for code, currency := range SupportedCurrencies {

		if strings.Contains(upper, code) {
			return currency, true
		}
	}

	return Currency{}, false
}

// DetectCurrencyCode
//
// Shortcut jika hanya membutuhkan code.
//
// Contoh:
//
// "HK$100"
// -> HKD
func DetectCurrencyCode(text string) string {

	currency, ok := DetectCurrency(text)

	if !ok {
		return ""
	}

	return currency.Code
}

// DetectCurrencySymbol
//
// Contoh:
//
// "HK$100"
//
// -> HK$
func DetectCurrencySymbol(text string) string {

	currency, ok := DetectCurrency(text)

	if !ok {
		return ""
	}

	return currency.Symbol
}

// DetectCurrencyScale
//
// Contoh:
//
// USD -> 2
//
// JPY -> 0
func DetectCurrencyScale(text string) int {

	currency, ok := DetectCurrency(text)

	if !ok {
		return 2
	}

	return currency.Scale
}

// IsSupportedCurrency
//
// Mengecek apakah currency didukung.
func IsSupportedCurrency(code string) bool {

	_, ok := SupportedCurrencies[strings.ToUpper(code)]

	return ok
}
