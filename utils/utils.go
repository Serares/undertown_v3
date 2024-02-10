package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Serares/undertown_v3/utils/constants"
)

func ReplyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	http.Error(w, message, status)
}

func ReplySuccess(w http.ResponseWriter, r *http.Request, status int, message interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(message)
	if err != nil {

		w.WriteHeader(500)
		return fmt.Errorf("error marshalling JSON: %s", err)
	}
	w.WriteHeader(status)
	w.Write(dat)
	return nil
}

func CheckIfStructIsEmpty(s interface{}) bool {
	// Get the value of the struct
	val := reflect.ValueOf(s)

	// If the struct is a pointer, find the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Iterate over all fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			return false // Found a non-zero field
		}
	}
	return true
}

func CreateDisplayPrice(price int64) string {
	return strconv.Itoa(int(price)) + " €"
}

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

func UrlEncodeString(s string) string {
	return url.PathEscape(s)
}

func UrlDecodeString(s string) (string, error) {
	return url.PathUnescape(s)
}

func BoolToInt(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func AddParamToUrl(baseUrl, param, value string) (string, error) {
	if baseUrl == "" {
		return "", fmt.Errorf("the base url is empty")
	}
	parsedURL, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("error parsing the url %v; %s; %s", err, param, value)
	}

	// Prepare query values
	parameters := url.Values{}
	parameters.Add(param, value)
	// Add more parameters as needed

	// Attach query parameters to the URL
	parsedURL.RawQuery = parameters.Encode()

	// The final URL with query parameters
	finalURL := parsedURL.String()

	return finalURL, nil
}

// TODO the base path should be a variable?
func CreateImagePath(imageName string) string {
	return "/assets/uploads/" + imageName
}

func CreateImagePathList(imageNames []string) []string {
	listOfPaths := make([]string, 0)

	for _, imgName := range imageNames {
		imagePath := CreateImagePath(imgName)
		listOfPaths = append(listOfPaths, imagePath)
	}

	return listOfPaths
}

// this is used for homepage single property
func CreateSinglePropertyPath(transactionType, title, humanReadableId string) (string, error) {
	underscoredTitle := ReplaceWhiteSpaceWithUnderscore(title)
	translatedTransactionType, err := TranslatePropertyTransactionType(transactionType)
	if err != nil {
		return "", err
	}
	firstPart := fmt.Sprintf("/%s/%s", translatedTransactionType, underscoredTitle)

	url, err := AddParamToUrl(firstPart, constants.HumanReadableIdQueryKey, humanReadableId)
	if err != nil {
		return "", err
	}
	return url, nil
}

// used to just attach the title and HumanReadableId query string
func CreatePropertyPath(baseUrl, title, humanReadableId string) (string, error) {
	underscoredTitle := ReplaceWhiteSpaceWithUnderscore(title)
	firstPart := fmt.Sprintf("%s/%s", baseUrl, underscoredTitle)

	url, err := AddParamToUrl(firstPart, constants.HumanReadableIdQueryKey, humanReadableId)
	if err != nil {
		return "", err
	}

	return url, nil
}

// replaces the whitespace with an underscore
func ReplaceWhiteSpaceWithUnderscore(s string) string {
	newString := strings.Split(s, " ")

	joinedString := strings.Join(newString, "_")

	return joinedString
}

// Return a translated transaction type
// See ssr/admin/types/PropertyTransactions
func TranslatePropertyTransactionType(transactionType string) (string, error) {
	if transactionType == "" {
		return "", fmt.Errorf("transaction type not provided")
	}

	if strings.EqualFold(transactionType, "SELL") {
		return constants.TranslatedTransactionSell, nil
	}

	return constants.TranslatedTransactionRent, nil
}
