package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"ne-go/v1/internal/api"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSSHUserGet(t *testing.T) {
	//given
	resp := api.SSHUserInfoVerbose{}
	if err := readJSONData("./test-fixtures/ne_sshuser_get_resp.json", &resp); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	baseURL := "http://localhost:8888"
	userID := "myTestUser"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s", baseURL, userID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	user, err := c.GetSSHUser(userID)

	//then
	assert.NotNil(t, user, "Returned user is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyUser(t, *user, resp)
}

func TestSSHUserCreate(t *testing.T) {
	//given
	resp := api.SSHUserCreateResponse{}
	if err := readJSONData("./test-fixtures/ne_sshuser_create_resp.json", &resp); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	baseURL := "http://localhost:8888"
	user := SSHUser{
		Username:    "myUser",
		Password:    "myPassword",
		DeviceUUIDs: []string{"deviceOne"},
	}
	req := api.SSHUserCreateRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/services/ssh-user", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(201, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateSSHUser(user.Username, user.Password, user.DeviceUUIDs[0])

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, resp.UUID, uuid, "UUID matches")
	verifyUserRequest(t, user, req)
}

func TestSSHUserUpdate(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	userID := "myTestUser"
	newPass := "myNewPassword"
	newDevices := []string{"nDev1", "nDev2"}
	delDevices := []string{"rDev1", "rDev2"}
	req := api.SSHUserUpdateRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s", baseURL, userID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(201, ""), nil
		},
	)
	for _, dev := range newDevices {
		httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s", baseURL, userID, dev),
			httpmock.NewStringResponder(201, ""))
	}
	for _, dev := range delDevices {
		httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s", baseURL, userID, dev),
			httpmock.NewStringResponder(200, ""))
	}
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	updateReq := c.NewSSHUserUpdateRequest(userID).
		WithNewPassword(newPass).
		WithNewDevices(newDevices).
		WithRemovedDevices(delDevices)
	err := updateReq.Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	verifyUserUpdateRequest(t, updateReq.(*restSSHUserUpdateRequest), req)
	for p, c := range httpmock.GetCallCountInfo() {
		assert.Equal(t, 1, c, "One request received on mock responder %s", p)
	}
}

func TestSSHUserDelete(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	userID := "myTestUser"
	userResp := api.SSHUserInfoVerbose{
		UUID:        userID,
		Username:    "user",
		DeviceUuids: []string{"dev1", "dev2", "dev3"}}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s", baseURL, userID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, userResp)
			return resp, nil
		},
	)
	for _, dev := range userResp.DeviceUuids {
		httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s", baseURL, userID, dev),
			httpmock.NewStringResponder(200, ""))
	}
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteSSHUser(userID)

	//then
	assert.Nil(t, err, "Error is not returned")
	for p, c := range httpmock.GetCallCountInfo() {
		assert.Equal(t, 1, c, "One request received on mock responder %s", p)
	}
}

func verifyUser(t *testing.T, user SSHUser, resp api.SSHUserInfoVerbose) {
	assert.Equal(t, resp.UUID, user.UUID, "UUID matches")
	assert.Equal(t, resp.Username, user.Username, "Username matches")
	assert.ElementsMatch(t, resp.DeviceUuids, user.DeviceUUIDs, "DeviceUUIDs match")
	assert.ElementsMatch(t, resp.Metros, user.MetroCodes, "Metros match")
}

func verifyUserRequest(t *testing.T, user SSHUser, req api.SSHUserCreateRequest) {
	assert.Equal(t, user.Username, *req.Username, "Username matches")
	assert.Equal(t, user.Password, *req.Password, "Password matches")
	assert.Equal(t, user.DeviceUUIDs[0], *req.DeviceUUID, "First DeviceUUID matches")
}

func verifyUserUpdateRequest(t *testing.T, updateReq *restSSHUserUpdateRequest, req api.SSHUserUpdateRequest) {
	assert.Equal(t, updateReq.newPassword, *req.Password, "Password matches")
}
