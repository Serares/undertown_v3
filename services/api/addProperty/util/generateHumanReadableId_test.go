package util

import (
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository/types"
)

func TestHumanReadableId(t *testing.T) {
	cases := []struct {
		name            string
		transactionType types.TransactionType
		expectedPrefix  string
	}{
		{name: "Sell Type",
			transactionType: types.Sell,
			expectedPrefix:  "SE"},
		{name: "Rent Type",
			transactionType: types.Rent,
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
