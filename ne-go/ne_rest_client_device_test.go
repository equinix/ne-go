package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"ne-go/internal/api"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testDevice = Device{
	AdditionalBandwidth: 100,
	DeviceTypeCode:      "PA-VM",
	HostName:            "myhostSRmy",
	LicenseType:         "BYOL",
	LicenseToken:        "I3372903",
	MetroCode:           "SV",
	Notifications:       []string{"test1@example.com", "test2@example.com"},
	PackageCode:         "VM100",
	TermLength:          24,
	Throughput:          1,
	ThroughputUnit:      "Gbps",
	Name:                "PaloAltoSRmy",
	ACL:                 []string{"192.168.1.1/32"},
	AccountNumber:       "1777643"}

func TestCreateDevice(t *testing.T) {
	//given
	resp := api.VirtualDeviceCreateResponse{}
	if err := readJSONData("./test-fixtures/ne_device_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannont read test response")
	}
	baseURL := "http://localhost:8888"
	device := testDevice
	req := api.VirtualDeviceRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/device", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(202, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateDevice(device)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, uuid, resp.UUID, "UUID matches")
	verifyDeviceRequest(t, device, req)
}

func TestCreateRedundantDevice(t *testing.T) {
	//given
	resp := api.VirtualDeviceCreateResponseDto{}
	if err := readJSONData("./test-fixtures/ne_device_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannont read test response")
	}
	baseURL := "http://localhost:8888"
	req := api.VirtualDeviceRequest{}
	primary := testDevice
	secondary := Device{
		AccountNumber:       "222222",
		ACL:                 primary.ACL,
		AdditionalBandwidth: 500,
		DeviceTypeCode:      "PA-SEC",
		HostName:            "secondaryHost",
		LicenseFileID:       "",
		LicenseKey:          "licKey",
		LicenseType:         "",
		LicenseSecret:       "licSecret",
		LicenseToken:        "licToken",
		MetroCode:           "DC",
		Notifications:       []string{"sec@sec.com", "secTwo@sec.com"},
		PackageCode:         "VM222",
		SiteID:              "secSiteId",
		SystemIPAddress:     "192.168.1.1",
		Throughput:          5,
		ThroughputUnit:      "Gbps",
		Name:                "secondary"}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/device", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(202, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	pUUID, sUUID, err := c.CreateRedundantDevice(primary, secondary)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, resp.UUID, pUUID, "Primary device UUID matches")
	assert.Equal(t, resp.SecondaryUUID, sUUID, "Secondary device UUID matches")
	verifyRedundantDeviceRequest(t, primary, secondary, req)
}

func TestGetDevice(t *testing.T) {
	//given
	resp := api.VirtualDeviceDetailsResponse{}
	if err := readJSONData("./test-fixtures/ne_device_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannont read test response")
	}
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/device/%s", baseURL, devID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	dev, err := c.GetDevice(devID)

	//then
	assert.NotNil(t, dev, "Returned device is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyDevice(t, *dev, resp)
}

func TestUpdateDeviceBasicFields(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	newName := "myNewName"
	newNotifications := []string{"new@new.com", "new2@new.com"}
	newTermLength := 24
	req := api.VirtualDeviceInternalPatchRequestDto{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/device/%s", baseURL, devID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(200, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.NewDeviceUpdateRequest(devID).WithDeviceName(newName).
		WithNotifications(newNotifications).WithTermLength(newTermLength).Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, newName, req.VirtualDeviceName, "DeviceName matches")
	assert.ElementsMatch(t, newNotifications, req.Notifications, "Notifications match")
	assert.Equal(t, int64(newTermLength), req.TermLength, "TermLength match")
}

func TestUpdateDeviceACL(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	newACLs := []string{"127.0.0.1/32", "192.168.0.0/24"}
	testHc := &http.Client{}
	req := make([]string, 0)
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/device/%s/acl", baseURL, devID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.NewDeviceUpdateRequest(devID).WithACLs(newACLs).Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.ElementsMatch(t, newACLs, req, "ACLs match")
}

func TestUpdateDeviceAdditionalBandwidth(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	newBandwidth := 1000
	testHc := &http.Client{}
	req := api.AdditionalBandwidthUpdateRequest{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/device/additionalbandwidth/%s", baseURL, devID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.NewDeviceUpdateRequest(devID).WithAdditionalBandwidth(newBandwidth).Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, int32(newBandwidth), *req.AdditionalBandwidth, "AdditionalBandwidth match")
}

func TestDeleteDevice(t *testing.T) {
	//given
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/device/%s", baseURL, devID),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteDevice(devID)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifyDeviceRequest(t *testing.T, dev Device, req api.VirtualDeviceRequest) {
	assert.Equal(t, req.AccountNumber, dev.AccountNumber, "AccountNumber matches")
	assert.ElementsMatch(t, req.FqdnACL, mapACLsDomainToAPI(dev.ACL), "ACL matches")
	assert.Equal(t, req.AdditionalBandwidth, int32(dev.AdditionalBandwidth), "AdditionalBandwidth matches")
	if dev.DeviceTypeCode != "" {
		assert.Equal(t, *req.DeviceTypeCode, dev.DeviceTypeCode, "DeviceTypeCode matches")
	}
	if dev.HostName != "" {
		assert.Equal(t, *req.HostNamePrefix, dev.HostName, "HostNamePrefix matches")
	}
	assert.Equal(t, req.LicenseFileID, dev.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, req.LicenseKey, dev.LicenseKey, "LicenseKey matches")
	if dev.LicenseType != "" {
		assert.Equal(t, *req.LicenseMode, dev.LicenseType, "LicenseMode matches")
	}
	assert.Equal(t, req.LicenseSecret, dev.LicenseSecret, "LicenseSecret matches")
	assert.Equal(t, req.LicenseToken, dev.LicenseToken, "LicenseToken matches")
	if dev.MetroCode != "" {
		assert.Equal(t, *req.MetroCode, dev.MetroCode, "MetroCode matches")
	}
	assert.ElementsMatch(t, req.Notifications, dev.Notifications, "Notifications matches")
	assert.Equal(t, req.PackageCode, dev.PackageCode, "PackageCode matches")
	assert.Equal(t, req.SiteID, dev.SiteID, "SiteID matches")
	assert.Equal(t, req.SystemIPAddress, dev.SystemIPAddress, "SystemIPAddress matches")
	assert.Equal(t, req.Throughput, int32(dev.Throughput), "Throughput matches")
	assert.Equal(t, req.ThroughputUnit, dev.ThroughputUnit, "ThroughputUnit matches")
	if dev.Name != "" {
		assert.Equal(t, *req.VirtualDeviceName, dev.Name, "VirtualDeviceName matches")
	}
	assert.Equal(t, req.Version, dev.Version, "Version matches")
}

func verifyRedundantDeviceRequest(t *testing.T, primary Device, secondary Device, req api.VirtualDeviceRequest) {
	verifyDeviceRequest(t, primary, req)
	assert.Equal(t, secondary.AccountNumber, req.Secondary.AccountNumber, "Account number matches")
	assert.ElementsMatch(t, secondary.ACL, req.Secondary.ACL, "ACL matches")
	assert.Equal(t, int32(secondary.AdditionalBandwidth), req.Secondary.AdditionalBandwidth, "AdditionalBandwidth matches")
	assert.Equal(t, secondary.LicenseFileID, req.Secondary.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, secondary.LicenseKey, req.Secondary.LicenseKey, "LicenseKey matches")
	assert.Equal(t, secondary.LicenseSecret, req.Secondary.LicenseSecret, "LicenseSecret matches")
	assert.Equal(t, secondary.LicenseToken, req.Secondary.LicenseToken, "LicenseToken matches")
	if secondary.MetroCode != "" {
		assert.Equal(t, secondary.MetroCode, *req.Secondary.MetroCode, "MetroCode matches")
	}
	assert.ElementsMatch(t, secondary.Notifications, req.Secondary.Notifications, "Notifications match")
	assert.Equal(t, secondary.SiteID, req.Secondary.SiteID, "SiteID matches")
	assert.Equal(t, secondary.SystemIPAddress, req.Secondary.SystemIPAddress, "SystemIPAddress matches")
	if secondary.Name != "" {
		assert.Equal(t, secondary.Name, *req.Secondary.VirtualDeviceName, "VirtualDeviceName matches")
	}
	assert.Equal(t, secondary.HostName, req.Secondary.HostNamePrefix, "HostName matches")
}

func verifyDevice(t *testing.T, dev Device, resp api.VirtualDeviceDetailsResponse) {
	assert.Equal(t, resp.AccountNumber, dev.AccountNumber, "AccountNumber matches")
	assert.ElementsMatch(t, resp.ACL, dev.ACL, "ACL matches")
	assert.Equal(t, resp.AdditionalBandwidth, int32(dev.AdditionalBandwidth), "AdditionalBandwidth matches")
	assert.Equal(t, resp.Controller1, dev.Controller1, "Controller1 matches")
	assert.Equal(t, resp.Controller2, dev.Controller2, "Controller2 matches")
	assert.Equal(t, resp.DeviceSerialNo, dev.DeviceSerialNo, "DeviceSerialNo matches")
	assert.Equal(t, resp.DeviceTypeCategory, dev.DeviceTypeCategory, "DeviceTypeCategory matches")
	assert.Equal(t, resp.DeviceTypeCode, dev.DeviceTypeCode, "DeviceTypeCode matches")
	assert.Equal(t, resp.DeviceTypeName, dev.DeviceTypeName, "DeviceTypeName matches")
	assert.Equal(t, resp.DeviceTypeVendor, dev.DeviceTypeVendor, "DeviceTypeVendor matches")
	assert.Equal(t, resp.Expiry, dev.Expiry, "Expiry matches")
	assert.Equal(t, resp.HostName, dev.HostName, "HostName matches")
	assert.Equal(t, resp.LicenseFileID, dev.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, resp.LicenseKey, dev.LicenseKey, "LicenseKey matches")
	assert.Equal(t, resp.LicenseName, dev.LicenseName, "LicenseName matches")
	assert.Equal(t, resp.LicenseSecret, dev.LicenseSecret, "LicenseSecret matches")
	assert.Equal(t, resp.LicenseStatus, dev.LicenseStatus, "LicenseStatus matches")
	assert.Equal(t, resp.LicenseType, dev.LicenseType, "LicenseType matches")
	assert.Equal(t, resp.LocalID, dev.LocalID, "LocalID matches")
	assert.Equal(t, resp.ManagementGatewayIP, dev.ManagementGatewayIP, "ManagementGatewayIP matches")
	assert.Equal(t, resp.ManagementIP, dev.ManagementIP, "ManagementIP matches")
	assert.Equal(t, resp.MetroCode, dev.MetroCode, "MetroCode matches")
	assert.Equal(t, resp.MetroName, dev.MetroName, "MetroName matches")
	assert.Equal(t, resp.Name, dev.Name, "Name matches")
	assert.ElementsMatch(t, resp.Notifications, dev.Notifications, "Notifications matches")
	assert.Equal(t, resp.PackageCode, dev.PackageCode, "PackageCode matches")
	assert.Equal(t, resp.PackageName, dev.PackageName, "PackageName matches")
	assert.Equal(t, resp.PrimaryDNSName, dev.PrimaryDNSName, "PrimaryDNSName matches")
	assert.Equal(t, resp.PublicGatewayIP, dev.PublicGatewayIP, "PublicGatewayIP matches")
	assert.Equal(t, resp.PublicIP, dev.PublicIP, "PublicIP matches")
	assert.Equal(t, resp.PublicIP, dev.PurchaseOrderNumber, "PublicIP matches")
	assert.Equal(t, resp.RedundancyType, dev.RedundancyType, "RedundancyType matches")
	assert.Equal(t, resp.RedundantUUID, dev.RedundantUUID, "RedundantUUID matches")
	assert.Equal(t, resp.Region, dev.Region, "Region matches")
	assert.Equal(t, resp.RemoteID, dev.RemoteID, "RemoteID matches")
	assert.Equal(t, resp.SecondaryDNSName, dev.SecondaryDNSName, "SecondaryDNSName matches")
	assert.Equal(t, resp.SerialNumber, dev.SerialNumber, "SerialNumber matches")
	assert.Equal(t, resp.SiteID, dev.SiteID, "SiteID matches")
	assert.Equal(t, resp.SSHIPAddress, dev.SSHIPAddress, "SSHIPAddress matches")
	assert.Equal(t, resp.SSHIPFqdn, dev.SSHIPFqdn, "SSHIPFqdn matches")
	assert.Equal(t, resp.Status, dev.Status, "Status matches")
	assert.Equal(t, resp.SystemIPAddress, dev.SystemIPAddress, "SystemIPAddress matches")
	assert.Equal(t, resp.TermLength, int32(dev.TermLength), "TermLength matches")
	assert.Equal(t, resp.Throughput, fmt.Sprintf("%d", dev.Throughput), "Throughput matches")
	assert.Equal(t, resp.ThroughputUnit, dev.ThroughputUnit, "ThroughputUnit matches")
	assert.Equal(t, resp.UUID, dev.UUID, "UUID matches")
	assert.Equal(t, resp.Version, dev.Version, "Version matches")
}
