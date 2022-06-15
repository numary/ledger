package core

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReverseTransaction(t *testing.T) {
	tx := &Transaction{
		TransactionData: TransactionData{
			Postings: Postings{
				{
					Source:      WorldAccount,
					Destination: "users:001",
					Amount:      100,
					Asset:       "COIN",
				},
				{
					Source:      "users:001",
					Destination: "payments:001",
					Amount:      100,
					Asset:       "COIN",
				},
			},
			Reference: "foo",
		},
	}

	expected := TransactionData{
		Postings: Postings{
			{
				Source:      "payments:001",
				Destination: "users:001",
				Amount:      100,
				Asset:       "COIN",
			},
			{
				Source:      "users:001",
				Destination: WorldAccount,
				Amount:      100,
				Asset:       "COIN",
			},
		},
		Reference: "revert_foo",
	}

	if diff := cmp.Diff(expected, tx.Reverse()); diff != "" {
		t.Errorf("Reverse() mismatch (-want +got):\n%s", diff)
	}
}
