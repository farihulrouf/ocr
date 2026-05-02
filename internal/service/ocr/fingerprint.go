package ocr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func GenerateFingerprint(store string, date time.Time, total int64) string {
	raw := fmt.Sprintf("%s|%s|%d",
		strings.ToLower(strings.TrimSpace(store)),
		date.Format("2006-01-02"),
		total,
	)

	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}
