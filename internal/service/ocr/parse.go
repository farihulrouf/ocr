package ocr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//
// ==========================
// ENUM & STRUCT
// ==========================
//

type LineType string

const (
	LineStore    LineType = "STORE"
	LineItem     LineType = "ITEM"
	LineTax      LineType = "TAX"
	LineSubtotal LineType = "SUBTOTAL"
	LineTotal    LineType = "TOTAL"
	LineDate     LineType = "DATE"
	LineOther    LineType = "OTHER"
)

type DetectedLine struct {
	Index int
	Raw   string
	Type  LineType
}

//
// ==========================
// PUBLIC API
// ==========================
//

// ParseReceipt - FINAL parser struk Jepang
func ParseReceipt(text string) (
	storeName string,
	total int64,
	date *time.Time,
	taxID string,
	isQualified bool,
	subtotal int64,
	tax int64,
	items []string,
) {

	fmt.Println("==================================================")
	fmt.Println("[DEBUG][PARSE] START ParseReceipt")

	lines := detectLines(text)

	// ======================
	// STORE NAME
	// ======================
	for _, l := range lines {
		if l.Type == LineStore {
			storeName = strings.TrimSpace(l.Raw)
			fmt.Println("[DEBUG][PARSE] store_name =", storeName)
			break
		}
	}

	// ======================
	// DATE
	// ======================
	for _, l := range lines {
		if l.Type == LineDate {
			if d := parseDate(l.Raw); d != nil {
				date = d
				fmt.Println("[DEBUG][PARSE] date =", date)
				break
			}
		}
	}

	// ======================
	// ITEMS & SUBTOTAL
	// ======================
	for _, l := range lines {
		switch l.Type {

		case LineItem:
			if amt := parseYen(l.Raw); amt != nil {
				subtotal += *amt
				items = append(items, l.Raw)
				fmt.Printf("[DEBUG][PARSE] item + %d subtotal = %d\n", *amt, subtotal)
			}

		case LineSubtotal:
			if amt := parseYen(l.Raw); amt != nil {
				subtotal = *amt
				fmt.Println("[DEBUG][PARSE] subtotal (override) =", subtotal)
			}

		case LineTax:
			if amt := parseYen(l.Raw); amt != nil {
				tax += *amt
				fmt.Println("[DEBUG][PARSE] tax +", *amt)
			}

		case LineTotal:
			if amt := parseYen(l.Raw); amt != nil {
				if *amt >= subtotal {
					total = *amt
					fmt.Println("[DEBUG][PARSE] total =", total)
				} else {
					fmt.Println("[WARN] total ignored (smaller than subtotal)")
				}
			}
		}
	}

	// ======================
	// FALLBACK TOTAL
	// ======================
	if total == 0 {
		total = subtotal + tax
		fmt.Println("[DEBUG][PARSE] total fallback =", total)
	}

	// ======================
	// QUALIFIED
	// ======================
	isQualified = strings.Contains(text, "適格") ||
		(strings.Contains(text, "登録番号") && !strings.Contains(text, "対象外"))

	fmt.Println("[DEBUG][PARSE] items_count =", len(items))
	fmt.Println("[DEBUG][PARSE] subtotal =", subtotal)
	fmt.Println("[DEBUG][PARSE] tax =", tax)
	fmt.Println("[DEBUG][PARSE] total =", total)
	fmt.Println("[DEBUG][PARSE] qualified =", isQualified)
	fmt.Println("[DEBUG][PARSE] END ParseReceipt")

	return
}

//
// ==========================
// LINE DETECTION
// ==========================
//

func detectLines(text string) []DetectedLine {
	fmt.Println("[DEBUG][LINE] START DetectLines")

	rawLines := strings.Split(text, "\n")
	var result []DetectedLine

	for i, raw := range rawLines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}

		lt := detectLineType(line)

		fmt.Printf("[DEBUG][LINE] %02d | %-9s | %s\n",
			i+1, lt, line)

		result = append(result, DetectedLine{
			Index: i + 1,
			Raw:   line,
			Type:  lt,
		})
	}

	fmt.Println("[DEBUG][LINE] END DetectLines")
	return result
}

func detectLineType(line string) LineType {

	switch {
	case isTotalLine(line):
		return LineTotal
	case isSubtotalLine(line):
		return LineSubtotal
	case isTaxLine(line):
		return LineTax
	case isDateLine(line):
		return LineDate
	case isStoreLine(line):
		return LineStore
	case isItemLine(line):
		return LineItem
	default:
		return LineOther
	}
}

//
// ==========================
// RULES
// ==========================
//

func isItemLine(line string) bool {
	return regexp.MustCompile(`¥\s*[\d,]+`).MatchString(line) &&
		!isTotalLine(line) &&
		!isSubtotalLine(line) &&
		!isTaxLine(line)
}

func isStoreLine(line string) bool {
	if regexp.MustCompile(`¥\s*\d`).MatchString(line) {
		return false
	}
	if strings.Contains(line, "http") {
		return false
	}
	return len([]rune(line)) >= 3 && len([]rune(line)) <= 40
}

func isDateLine(line string) bool {
	return regexp.MustCompile(`\d{4}年.*\d{1,2}月.*\d{1,2}日`).MatchString(line)
}

func isTotalLine(line string) bool {
	keys := []string{"合計", "TOTAL"}
	for _, k := range keys {
		if strings.Contains(line, k) && strings.Contains(line, "¥") {
			return true
		}
	}
	return false
}

func isSubtotalLine(line string) bool {
	keys := []string{"小計", "SUBTOTAL"}
	for _, k := range keys {
		if strings.Contains(line, k) && strings.Contains(line, "¥") {
			return true
		}
	}
	return false
}

func isTaxLine(line string) bool {
	keys := []string{"消費税", "外税", "内税", "税"}
	if !strings.Contains(line, "¥") {
		return false
	}
	for _, k := range keys {
		if strings.Contains(line, k) {
			return true
		}
	}
	return false
}

//
// ==========================
// PARSERS
// ==========================
//

func parseYen(line string) *int64 {
	re := regexp.MustCompile(`¥\s*([\d,]+)`)
	if m := re.FindStringSubmatch(line); len(m) > 1 {
		val, _ := strconv.ParseInt(strings.ReplaceAll(m[1], ",", ""), 10, 64)
		return &val
	}
	return nil
}

func parseDate(line string) *time.Time {
	re := regexp.MustCompile(`(\d{4})年\s*(\d{1,2})月\s*(\d{1,2})日`)
	if m := re.FindStringSubmatch(line); len(m) == 4 {
		y, _ := strconv.Atoi(m[1])
		mo, _ := strconv.Atoi(m[2])
		d, _ := strconv.Atoi(m[3])
		t := time.Date(y, time.Month(mo), d, 0, 0, 0, 0, time.UTC)
		return &t
	}
	return nil
}
