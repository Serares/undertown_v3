package utils

import (
	"testing"
)

func TestHumanReadableId(t *testing.T) {
	cases := []struct {
		name            string
		transactionType TransactionType
		expectedPrefix  string
	}{
		{name: "Sell Type",
			transactionType: Sell,
			expectedPrefix:  "SE"},
		{name: "Rent Type",
			transactionType: Rent,
			expectedPrefix:  "RE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id := HumanReadableId(tc.transactionType)
			if id[:2] != tc.expectedPrefix {
				t.Errorf("expected the prefix %s and got %s", tc.expectedPrefix, id[:2])
			}
		})
	}
}
