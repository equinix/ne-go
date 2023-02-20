package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testSSHPublicKey = SSHPublicKey{
	Name:  String("testKey"),
	Value: String("keyyyyyyyyyyyyyyyyyyyyyyyyy"),
	ProjectId: String("00d0a000-b000-0000-0000-00000c0f0000"),
}

func TestGetSSHPublicKeys(t *testing.T) {
	//given
	resp := make([]api.SSHPublicKey, 0)
	if err := readJSONData("./test-fixtures/ne_sshpubkeys_get.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/publicKeys", baseURL), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	keys, err := c.GetSSHPublicKeys()

	//then
	assert.NotNil(t, keys, "Returned list of keys is not nil")
	assert.Nil(t, err, "Returned error is nil")
	assert.Equal(t, len(resp), len(keys), "Number of keys matches")
	for i := range keys {
		verifySSHPublicKey(t, resp[i], keys[i])
	}
}

func TestGetSSHPublicKey(t *testing.T) {
	//given
	resp := api.SSHPublicKey{}
	if err := readJSONData("./test-fixtures/ne_sshpubkey_get.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	keyUUID := "keyID"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/publicKeys/%s", baseURL, keyUUID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	key, err := c.GetSSHPublicKey(keyUUID)

	//then
	assert.NotNil(t, key, "Returned key is not nil")
	assert.Nil(t, err, "Returned error is nil")
	verifySSHPublicKey(t, resp, *key)
}

func TestCreateSSHPublicKey(t *testing.T) {
	//given
	newUUID := "e3e5d6ce-5238-4fea-a454-e4a74a7bd060"
	key := testSSHPublicKey
	req := api.SSHPublicKey{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/publicKeys", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp := httpmock.NewStringResponse(201, "")
			resp.Header.Add("Location", "/ne/v1/publicKeys/"+newUUID)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateSSHPublicKey(key)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, newUUID, *uuid, "UUID matches")
	verifySSHPublicKey(t, req, key)
}

func TestDeleteSSHPublicKey(t *testing.T) {
	//given
	keyUUID := "keyID"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/publicKeys/%s", baseURL, keyUUID),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteSSHPublicKey(keyUUID)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifySSHPublicKey(t *testing.T, apiKey api.SSHPublicKey, key SSHPublicKey) {
	assert.Equal(t, apiKey.UUID, key.UUID, "UUID matches")
	assert.Equal(t, apiKey.KeyName, key.Name, "Name matches")
	assert.Equal(t, apiKey.KeyValue, key.Value, "Value matches")
}

