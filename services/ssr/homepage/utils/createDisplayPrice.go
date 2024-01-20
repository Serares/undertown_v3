package utils

import "strconv"

func CreateDisplayPrice(price int64) string {
	return strconv.Itoa(int(price)) + " €"
}
