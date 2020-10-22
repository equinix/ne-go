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

func TestSSHUserGet(t *testing.T) {
	//given
	resp := api.SSHUser{}
	if err := readJSONData("./test-fixtures/ne_sshuser_get_resp.json", &resp); err != nil {
		assert.Failf(t, "Cannot read test response due to %s", err.Error())
	}
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

func TestSSHUsersGet(t *testing.T) {
	//Given
	var respBody api.SSHUsersResponse
	if err := readJSONData("./test-fixtures/ne_sshusers_get.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	pageSize := 100
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/services/ssh-user?pageSize=%d&verbose=true", baseURL, pageSize),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = pageSize
	users, err := c.GetSSHUsers()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, users, "Client should return a response")
	assert.Equal(t, respBody.TotalCount, len(users))
	for i := range users {
		verifyUser(t, users[i], respBody.List[i])
	}
}

func TestSSHUserCreate(t *testing.T) {
	//given
	resp := api.SSHUserRequestResponse{}
	if err := readJSONData("./test-fixtures/ne_sshuser_create_resp.json", &resp); err != nil {
		assert.Failf(t, "Cannot read test response due to %s", err.Error())
	}
	user := SSHUser{
		Username:    "myUser",
		Password:    "myPassword",
		DeviceUUIDs: []string{"deviceOne"},
	}
	req := api.SSHUserRequest{}
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
	userID := "myTestUser"
	newPass := "myNewPassword"
	oldDevices := []string{"Dev1", "Dev2"}
	newDevices := []string{"Dev3", "Dev4"}
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
	removed, added := diffStringSlices(oldDevices, newDevices)
	for _, dev := range added {
		httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s", baseURL, userID, dev),
			httpmock.NewStringResponder(201, ""))
	}
	for _, dev := range removed {
		httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s", baseURL, userID, dev),
			httpmock.NewStringResponder(200, ""))
	}
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	updateReq := c.NewSSHUserUpdateRequest(userID).
		WithNewPassword(newPass).
		WithDeviceChange(oldDevices, newDevices)
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
	userID := "myTestUser"
	userResp := api.SSHUser{
		UUID:        userID,
		Username:    "user",
		DeviceUUIDs: []string{"dev1", "dev2", "dev3"}}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/services/ssh-user/%s", baseURL, userID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, userResp)
			return resp, nil
		},
	)
	for _, dev := range userResp.DeviceUUIDs {
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

func verifyUser(t *testing.T, user SSHUser, resp api.SSHUser) {
	assert.Equal(t, resp.UUID, user.UUID, "UUID matches")
	assert.Equal(t, resp.Username, user.Username, "Username matches")
	assert.ElementsMatch(t, resp.DeviceUUIDs, user.DeviceUUIDs, "DeviceUUIDs match")
}

func verifyUserRequest(t *testing.T, user SSHUser, req api.SSHUserRequest) {
	assert.Equal(t, user.Username, req.Username, "Username matches")
	assert.Equal(t, user.Password, req.Password, "Password matches")
	assert.Equal(t, user.DeviceUUIDs[0], req.DeviceUUID, "First DeviceUUID matches")
}

func verifyUserUpdateRequest(t *testing.T, updateReq *restSSHUserUpdateRequest, req api.SSHUserUpdateRequest) {
	assert.Equal(t, updateReq.newPassword, req.Password, "Password matches")
}
