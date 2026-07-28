package money

// Currency menyimpan metadata mata uang.
type Currency struct {
	Code    string // ISO4217
	Name    string
	Symbol  string
	Scale   int
	Country string
}

// SupportedCurrencies berisi seluruh currency yang didukung.
var SupportedCurrencies = map[string]Currency{

	// Asia

	"JPY": {
		Code:    "JPY",
		Name:    "Japanese Yen",
		Symbol:  "¥",
		Scale:   0,
		Country: "Japan",
	},

	"HKD": {
		Code:    "HKD",
		Name:    "Hong Kong Dollar",
		Symbol:  "HK$",
		Scale:   2,
		Country: "Hong Kong",
	},

	"SGD": {
		Code:    "SGD",
		Name:    "Singapore Dollar",
		Symbol:  "S$",
		Scale:   2,
		Country: "Singapore",
	},

	"CNY": {
		Code:    "CNY",
		Name:    "Chinese Yuan",
		Symbol:  "¥",
		Scale:   2,
		Country: "China",
	},

	"TWD": {
		Code:    "TWD",
		Name:    "New Taiwan Dollar",
		Symbol:  "NT$",
		Scale:   2,
		Country: "Taiwan",
	},

	"KRW": {
		Code:    "KRW",
		Name:    "Korean Won",
		Symbol:  "₩",
		Scale:   0,
		Country: "South Korea",
	},

	"THB": {
		Code:    "THB",
		Name:    "Thai Baht",
		Symbol:  "฿",
		Scale:   2,
		Country: "Thailand",
	},

	"IDR": {
		Code:    "IDR",
		Name:    "Indonesian Rupiah",
		Symbol:  "Rp",
		Scale:   0,
		Country: "Indonesia",
	},

	// Europe

	"EUR": {
		Code:    "EUR",
		Name:    "Euro",
		Symbol:  "€",
		Scale:   2,
		Country: "European Union",
	},

	"GBP": {
		Code:    "GBP",
		Name:    "British Pound",
		Symbol:  "£",
		Scale:   2,
		Country: "United Kingdom",
	},

	"CHF": {
		Code:    "CHF",
		Name:    "Swiss Franc",
		Symbol:  "CHF",
		Scale:   2,
		Country: "Switzerland",
	},

	// America

	"USD": {
		Code:    "USD",
		Name:    "US Dollar",
		Symbol:  "$",
		Scale:   2,
		Country: "United States",
	},

	"CAD": {
		Code:    "CAD",
		Name:    "Canadian Dollar",
		Symbol:  "C$",
		Scale:   2,
		Country: "Canada",
	},

	// Oceania

	"AUD": {
		Code:    "AUD",
		Name:    "Australian Dollar",
		Symbol:  "A$",
		Scale:   2,
		Country: "Australia",
	},

	"NZD": {
		Code:    "NZD",
		Name:    "New Zealand Dollar",
		Symbol:  "NZ$",
		Scale:   2,
		Country: "New Zealand",
	},
}

// FindCurrency mencari currency berdasarkan code.
func FindCurrency(code string) (Currency, bool) {

	c, ok := SupportedCurrencies[code]

	return c, ok
}
