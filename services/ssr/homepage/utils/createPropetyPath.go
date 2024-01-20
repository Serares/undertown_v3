package utils

import "strings"

// ❗ The GET_PROPERTY backend should search for the propertyId query string
func CreatePropertyPath(title, humanReadableId string) string {
	replacedString := strings.ReplaceAll(title, " ", string('%'))
	return strings.ToLower(replacedString) + strings.Join([]string{"?propertyId=", humanReadableId}, "")
}
