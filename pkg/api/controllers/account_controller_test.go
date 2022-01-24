package controllers_test

import (
	"github.com/numary/ledger/pkg/api"
	"github.com/numary/ledger/pkg/api/internal"
	"github.com/numary/ledger/pkg/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetAccounts(t *testing.T) {

	internal.RunTest(t, func(h *api.API) {

		rsp := internal.PostTransaction(t, h, core.Transaction{
			Postings: core.Postings{
				{
					Source:      "world",
					Destination: "alice",
					Amount:      100,
					Asset:       "USD",
				},
			},
		})
		assert.Equal(t, http.StatusOK, rsp.Result().StatusCode)

		rsp = internal.PostTransaction(t, h, core.Transaction{
			Postings: core.Postings{
				{
					Source:      "world",
					Destination: "bob",
					Amount:      100,
					Asset:       "USD",
				},
			},
		})
		assert.Equal(t, http.StatusOK, rsp.Result().StatusCode)

		rsp = internal.GetAccounts(h)
		assert.Equal(t, http.StatusOK, rsp.Result().StatusCode)

		cursor := internal.DecodeCursorResponse(t, rsp.Body, core.Account{})
		assert.EqualValues(t, 3, cursor.Total)
		assert.Len(t, cursor.Data, 3)
	})

}
