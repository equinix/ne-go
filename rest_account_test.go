package ne

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccounts(t *testing.T) {
	//given
	resp := api.AccountResponse{}
	if err := readJSONData("./test-fixtures/ne_accounts.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	metro := "SV"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/accounts/%s", baseURL, metro), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	accounts, err := c.GetAccounts(metro)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.NotNil(t, accounts, "Returned accounts slice is not nil")
	assert.Equal(t, len(resp.Accounts), len(accounts), "Number of accounts matches")
	for i := range resp.Accounts {
		verifyAccount(t, resp.Accounts[i], accounts[i])
	}
}

func verifyAccount(t *testing.T, apiAccount api.Account, account Account) {
	assert.Equal(t, apiAccount.Name, account.Name, "Account Name matches")
	assert.Equal(t, apiAccount.Number, account.Number, "Account Number matches")
	assert.Equal(t, apiAccount.Status, account.Status, "Account Status matches")
	assert.Equal(t, apiAccount.UCMID, account.UCMID, "Account UCMID matches")
}
