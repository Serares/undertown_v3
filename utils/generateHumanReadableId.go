package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ❔
// create a maximu length of 8
// use the TransactionType enum from the repository models
// define the first two letters as the transaction type (SE || RE)
// follow by 3+3 random generated numbers
func HumanReadableId(propertyType string) string {
	var randomBlock []string

	for i := 0; i < 2; i++ {
		randomBlock = append(randomBlock, generateRandomBlock(5, int32(i+1)))
	}

	transactionTypeTrimmedCapitalized := strings.ToUpper(propertyType[:2])
	generatedId := fmt.Sprintf("%s%s%s", transactionTypeTrimmedCapitalized, randomBlock[0], randomBlock[1])

	return generatedId
}

func generateRandomInt() int {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)

	return random.Int()
}

func generateRandomBlock(numOfDigits, n int32) string {
	for {
		randomInt := generateRandomInt()
		randomInt = randomInt * int(n)
		intToString := fmt.Sprintf("%d", randomInt)
		if len(intToString) > int(numOfDigits) {
			return intToString[1 : numOfDigits+1]
		}
	}
}
