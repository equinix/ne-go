package ne

import (
	"fmt"
	"ne-go/v1/tmp/api"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	associateDevice   = "ADD"
	unassociateDevice = "DELETE"
)

type restSSHUserUpdateRequest struct {
	uuid           string
	newPassword    string
	newDevices     []string
	removedDevices []string
	c              RestClient
}

func (c RestClient) CreateSSHUser(username string, password string, device string) (string, error) {
	u := c.baseURL + "/ne/v1/services/ssh-user"
	reqBody := api.SSHUserCreateRequest{
		Username:   &username,
		Password:   &password,
		DeviceUUID: &device,
	}
	respBody := api.SSHUserCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.execute(req, resty.MethodPost, u); err != nil {
		return "", err
	}
	return respBody.UUID, nil
}

func (c RestClient) GetSSHUser(uuid string) (*SSHUser, error) {
	url := c.baseURL + "/ne/v1/services/ssh-user/" + url.PathEscape(uuid)
	respBody := api.SSHUserInfoVerbose{}
	req := c.R().SetResult(&respBody)
	if err := c.execute(req, resty.MethodGet, url); err != nil {
		return nil, err
	}
	return mapSSHUserAPIToDomain(respBody), nil
}

func (client RestClient) NewSSHUserUpdateRequest(uuid string) SSHUserUpdateRequest {
	return &restSSHUserUpdateRequest{
		uuid: uuid,
		c:    client}
}

func (req *restSSHUserUpdateRequest) WithNewPassword(password string) SSHUserUpdateRequest {
	req.newPassword = password
	return req
}

func (req *restSSHUserUpdateRequest) WithNewDevices(uuids []string) SSHUserUpdateRequest {
	req.newDevices = uuids
	return req
}

func (req *restSSHUserUpdateRequest) WithRemovedDevices(uuids []string) SSHUserUpdateRequest {
	req.removedDevices = uuids
	return req
}

func (req *restSSHUserUpdateRequest) Execute() error {
	updateErr := UpdateError{}
	if req.newPassword != "" {
		if err := req.c.changeUserPassword(req.uuid, req.newPassword); err != nil {
			updateErr.failed = append(updateErr.failed, ChangeError{
				Type:   ChangeTypeUpdate,
				Target: "password",
				Value:  req.newPassword,
				Cause:  err})
		}
	}
	for _, dev := range req.newDevices {
		if err := req.c.changeDeviceAssociation(associateDevice, req.uuid, dev); err != nil {
			updateErr.failed = append(updateErr.failed, ChangeError{
				Type:   ChangeTypeCreate,
				Target: "devices",
				Value:  dev,
				Cause:  err})
		}
	}
	for _, dev := range req.removedDevices {
		if err := req.c.changeDeviceAssociation(unassociateDevice, req.uuid, dev); err != nil {
			updateErr.failed = append(updateErr.failed, ChangeError{
				Type:   ChangeTypeDelete,
				Target: "devices",
				Value:  dev,
				Cause:  err})
		}
	}
	if len(updateErr.failed) > 0 {
		return updateErr
	}
	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func (c RestClient) changeUserPassword(userID string, newPassword string) error {
	url := fmt.Sprintf("%s/ne/v1/services/ssh-user/%s",
		c.baseURL, url.PathEscape(userID))
	reqBody := api.SSHUserUpdateRequest{Password: &newPassword}
	req := c.R().SetBody(&reqBody)
	if err := c.execute(req, resty.MethodPut, url); err != nil {
		return err
	}
	return nil
}

func (c RestClient) changeDeviceAssociation(changeType string, userID string, deviceID string) error {
	url := fmt.Sprintf("%s/ne/v1/services/ssh-user/%s/association?deviceUuid=%s",
		c.baseURL, url.PathEscape(userID), url.PathEscape(deviceID))
	var method string
	switch changeType {
	case associateDevice:
		method = resty.MethodPatch
	case unassociateDevice:
		method = resty.MethodDelete
	default:
		return fmt.Errorf("unsupported association change type")
	}
	req := c.R().
		//due to bug in NE API that requires content type and content len = 0 altough there is no content needed in any case
		SetHeader("Content-Type", "application/json").
		SetBody("{}")
		//SetContentLength(true).
		//SetHeader("Content-Length", "0")
	if err := c.execute(req, method, url); err != nil {
		return err
	}
	return nil
}

func mapSSHUserAPIToDomain(apiUser api.SSHUserInfoVerbose) *SSHUser {
	return &SSHUser{
		UUID:        apiUser.UUID,
		Username:    apiUser.Username,
		DeviceUUIDs: apiUser.DeviceUuids,
		MetroCodes:  apiUser.Metros}
}
