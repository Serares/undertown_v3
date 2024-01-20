package utils

import (
	"fmt"
	"time"
)

func CreateDisplayCreatedAt(createdAt time.Time) string {
	passedTime := int(time.Since(createdAt).Hours() / 24)
	if passedTime < 0 {
		return "Adaugat recent"
	}

	if passedTime > 10 {
		return ""
	}

	return fmt.Sprintf("Adaugat cu %d zile in urma", passedTime)
}
