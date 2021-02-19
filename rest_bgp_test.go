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

var testBGPConfiguration = BGPConfiguration{
	ConnectionUUID:    String("e8b2e48e-2eba-4412-bc0b-c88dadb48050"),
	LocalIPAddress:    String("10.0.0.1/30"),
	LocalASN:          Int(10012),
	RemoteIPAddress:   String("10.0.0.2"),
	RemoteASN:         Int(10013),
	AuthenticationKey: String("authKey"),
}

func TestCreateBGPConfiguration(t *testing.T) {
	//given
	resp := api.BGPConfigurationCreateResponse{}
	if err := readJSONData("./test-fixtures/ne_bgp_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	bgpConfig := testBGPConfiguration
	reqBody := api.BGPConfiguration{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/bgp", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(202, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateBGPConfiguration(bgpConfig)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, uuid, resp.UUID, "UUID matches")
	verifyBGPConfig(t, bgpConfig, reqBody)
}

func TestGetBGPConfiguration(t *testing.T) {
	//given
	resp := api.BGPConfiguration{}
	if err := readJSONData("./test-fixtures/ne_bgp_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	bgpConfID := "e8b2e48e-2eba-4412-bc0b-c88dadb48050"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/bgp/%s", baseURL, bgpConfID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	bgpConf, err := c.GetBGPConfiguration(bgpConfID)

	//then
	assert.NotNil(t, bgpConf, "Returned device is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyBGPConfig(t, *bgpConf, resp)
}

func TestGetBGPConfigurationForConnection(t *testing.T) {
	//given
	resp := api.BGPConfiguration{}
	if err := readJSONData("./test-fixtures/ne_bgp_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	connID := "e8b2e48e-2eba-4412-bc0b-c88dadb48050"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/bgp/connection/%s", baseURL, connID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	bgpConf, err := c.GetBGPConfigurationForConnection(connID)

	//then
	assert.NotNil(t, bgpConf, "Returned device is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyBGPConfig(t, *bgpConf, resp)
}

func TestUpdateBGPConfiguration(t *testing.T) {
	//given
	resp := api.BGPConfigurationCreateResponse{}
	if err := readJSONData("./test-fixtures/ne_bgp_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	bgpConfID := "e8b2e48e-2eba-4412-bc0b-c88dadb48050"
	testHc := &http.Client{}
	reqBody := api.BGPConfiguration{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/bgp/%s", baseURL, bgpConfID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(202, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()
	newLocalASN := 12345
	newRemoteASN := 22241
	newLocalIPAddress := "1.1.1.1/30"
	newRemoteIPAddress := "1.1.1.2"
	newAuthenticationKey := "authKey"

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.NewBGPConfigurationUpdateRequest(bgpConfID).
		WithLocalASN(newLocalASN).
		WithRemoteASN(newRemoteASN).
		WithLocalIPAddress(newLocalIPAddress).
		WithRemoteIPAddress(newRemoteIPAddress).
		WithAuthenticationKey(newAuthenticationKey).
		Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	verifyBGPConfig(t, BGPConfiguration{
		LocalASN:          Int(newLocalASN),
		RemoteASN:         Int(newRemoteASN),
		LocalIPAddress:    String(newLocalIPAddress),
		RemoteIPAddress:   String(newRemoteIPAddress),
		AuthenticationKey: String(newAuthenticationKey),
	}, reqBody)
}

func verifyBGPConfig(t *testing.T, config BGPConfiguration, apiConfig api.BGPConfiguration) {
	assert.Equal(t, config.ConnectionUUID, apiConfig.ConnectionUUID, "ConnectionUUID matches")
	assert.Equal(t, config.LocalIPAddress, apiConfig.LocalIPAddress, "LocalIPAddress matches")
	assert.Equal(t, config.LocalASN, apiConfig.LocalASN, "LocalASN matches")
	assert.Equal(t, config.RemoteIPAddress, apiConfig.RemoteIPAddress, "RemoteIPAddress matches")
	assert.Equal(t, config.RemoteASN, apiConfig.RemoteASN, "RemoteASN matches")
	assert.Equal(t, config.AuthenticationKey, apiConfig.AuthenticationKey, "AuthenticationKey matches")
	assert.Equal(t, config.DeviceUUID, apiConfig.VirtualDeviceUUID, "DeviceUUID matches")
	assert.Equal(t, config.State, apiConfig.State, "State matches")
	assert.Equal(t, config.ProvisioningStatus, apiConfig.ProvisioningStatus, "ProvisioningStatus matches")
}
