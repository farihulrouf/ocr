package money

import (
	"fmt"
	"strconv"
	"strings"
)

// Format mengubah Money menjadi string.
//
// Contoh:
//
//	HKD 45430 -> HK$454.30
//	USD 1299  -> $12.99
//	JPY 540   -> ¥540
//	IDR 15000 -> Rp15.000
func Format(m Money) string {

	if m.Currency.Code == "" {
		return ""
	}

	scale := m.Currency.Scale

	negative := false
	amount := m.Amount

	if amount < 0 {
		negative = true
		amount = -amount
	}

	// Currency tanpa decimal
	if scale == 0 {

		value := formatThousands(amount)

		if negative {
			return "-" + m.Currency.Symbol + value
		}

		return m.Currency.Symbol + value
	}

	// Currency dengan decimal
	divisor := int64(1)

	for i := 0; i < scale; i++ {
		divisor *= 10
	}

	whole := amount / divisor
	fraction := amount % divisor

	value := fmt.Sprintf(
		"%s.%0*d",
		formatThousands(whole),
		scale,
		fraction,
	)

	if negative {
		return "-" + m.Currency.Symbol + value
	}

	return m.Currency.Symbol + value
}

// FormatAmount
//
// Menghasilkan hanya angka.
//
// Contoh:
//
//	45430 -> 454.30
func FormatAmount(m Money) string {

	scale := m.Currency.Scale

	if scale == 0 {
		return formatThousands(m.Amount)
	}

	divisor := int64(1)

	for i := 0; i < scale; i++ {
		divisor *= 10
	}

	whole := m.Amount / divisor
	fraction := m.Amount % divisor

	return fmt.Sprintf(
		"%s.%0*d",
		formatThousands(whole),
		scale,
		fraction,
	)
}

// FormatCode
//
// Contoh:
//
//	HKD 45430
func FormatCode(m Money) string {

	return fmt.Sprintf(
		"%s %s",
		m.Currency.Code,
		FormatAmount(m),
	)
}

// formatThousands
//
// 1234567
//
// ->
//
// 1,234,567
func formatThousands(value int64) string {

	s := strconv.FormatInt(value, 10)

	n := len(s)

	if n <= 3 {
		return s
	}

	var b strings.Builder

	first := n % 3

	if first == 0 {
		first = 3
	}

	b.WriteString(s[:first])

	for i := first; i < n; i += 3 {

		b.WriteString(",")

		b.WriteString(s[i : i+3])
	}

	return b.String()
}
