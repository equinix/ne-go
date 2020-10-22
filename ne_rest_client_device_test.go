package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testDevice = Device{
	AdditionalBandwidth: 100,
	TypeCode:            "PA-VM",
	HostName:            "myhostSRmy",
	IsBYOL:              true,
	LicenseToken:        "somelicensetokenaaaaazzzzz",
	MetroCode:           "SV",
	Notifications:       []string{"test1@example.com", "test2@example.com"},
	PackageCode:         "VM100",
	TermLength:          24,
	Throughput:          1,
	ThroughputUnit:      "Gbps",
	Name:                "PaloAltoSRmy",
	ACLs:                []string{"192.168.1.1/32"},
	AccountNumber:       "1777643",
	OrderReference:      "orderRef",
	PurchaseOrderNumber: "PO123456789",
	InterfaceCount:      10,
	CoreCount:           2,
	Version:             "10.09.05",
	IsSelfManaged:       true,
	VendorConfiguration: map[string]string{
		"serialNumber": "12312312",
		"controller1":  "1.1.1.1",
	},
}

func TestCreateDevice(t *testing.T) {
	//given
	resp := api.DeviceRequestResponse{}
	if err := readJSONData("./test-fixtures/ne_device_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	device := testDevice
	req := api.DeviceRequest{}
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
	resp := api.DeviceRequestResponse{}
	if err := readJSONData("./test-fixtures/ne_device_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	req := api.DeviceRequest{}
	primary := testDevice
	secondary := Device{
		MetroCode:           "DC",
		Name:                "secondary",
		Notifications:       []string{"secondary@secondary.com"},
		HostName:            "secondaryHostname",
		AccountNumber:       "99999",
		AdditionalBandwidth: 200,
		ACLs:                []string{"2.2.2.2/32"},
		VendorConfiguration: map[string]string{
			"serialNumber": "2222222",
			"controller1":  "2.2.2.2",
		}}
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
	resp := api.Device{}
	if err := readJSONData("./test-fixtures/ne_device_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
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

func TestGetDevices(t *testing.T) {
	//Given
	var respBody api.DevicesResponse
	if err := readJSONData("./test-fixtures/ne_devices_get.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	pageSize := respBody.PageSize
	statuses := []string{"INITIALIZING", "PROVISIONING"}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/device?size=%d&status=%s", baseURL, pageSize, url.QueryEscape("INITIALIZING,PROVISIONING")),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = pageSize
	devices, err := c.GetDevices(statuses)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, devices, "Client should return a response")
	assert.Equal(t, len(respBody.Content), len(devices), "Number of objects matches")
	for i := range respBody.Content {
		verifyDevice(t, devices[i], respBody.Content[i])
	}
}

func TestUpdateDeviceBasicFields(t *testing.T) {
	//given
	devID := "myDevice"
	newName := "myNewName"
	newNotifications := []string{"new@new.com", "new2@new.com"}
	newTermLength := 24
	req := api.DeviceUpdateRequest{}
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
	assert.Equal(t, newTermLength, req.TermLength, "TermLength match")
}

func TestUpdateDeviceACL(t *testing.T) {
	//given
	devID := "myDevice"
	newACLs := []string{"127.0.0.1/32", "192.168.0.0/24"}
	testHc := &http.Client{}
	req := make([]api.DeviceFqdnACL, 0)
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/device/%s/fqdn-acl", baseURL, devID),
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
	assert.ElementsMatch(t, req, mapDeviceACLsToFQDNACLs(newACLs), "ACL matches")
}

func TestUpdateDeviceAdditionalBandwidth(t *testing.T) {
	//given
	devID := "myDevice"
	newBandwidth := 1000
	testHc := &http.Client{}
	req := api.DeviceAdditionalBandwidthUpdateRequest{}
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
	assert.Equal(t, newBandwidth, req.AdditionalBandwidth, "AdditionalBandwidth match")
}

func TestDeleteDevice(t *testing.T) {
	//given
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

func TestGetFqdnACLs(t *testing.T) {
	//given
	resp := api.DeviceFqdnACLResponse{}
	if err := readJSONData("./test-fixtures/ne_device_fqdnAcls_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	baseURL := "http://localhost:8888"
	devID := "myDevice"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/device/%s/fqdn-acl", baseURL, devID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	acls, err := c.GetDeviceACLs(devID)

	//then
	assert.NotNil(t, acls, "Returned device is not nil")
	assert.Nil(t, err, "Error is not returned")

	cnt := 0
	for i := range resp.FqdnACLs {
		cnt += len(resp.FqdnACLs[i].CIDRs)
	}
	assert.Equal(t, cnt, len(acls.ACLs), "Number of CIDRS and ACLs matches")
	assert.Equal(t, resp.Status, acls.Status, "Status matches")
}

func verifyDevice(t *testing.T, device Device, resp api.Device) {
	assert.Equal(t, resp.UUID, device.UUID, "UUID matches")
	assert.Equal(t, resp.Name, device.Name, "Name matches")
	assert.Equal(t, resp.DeviceTypeCode, device.TypeCode, "DeviceTypeCode matches")
	assert.Equal(t, resp.Status, device.Status, "Status matches")
	assert.Equal(t, resp.LicenseStatus, device.LicenseStatus, "LicenseStatus matches")
	assert.Equal(t, resp.MetroCode, device.MetroCode, "MetroCode matches")
	assert.Equal(t, resp.IBX, device.IBX, "IBX matches")
	assert.Equal(t, resp.Region, device.Region, "Region matches")
	assert.Equal(t, resp.Throughput, device.Throughput, "Throughput matches")
	assert.Equal(t, resp.ThroughputUnit, device.ThroughputUnit, "ThroughputUnit matches")
	assert.Equal(t, resp.HostName, device.HostName, "HostName matches")
	assert.Equal(t, resp.PackageCode, device.PackageCode, "PackageCode matches")
	assert.Equal(t, resp.Version, device.Version, "Version matches")
	if resp.LicenseType == deviceLicenseModeSubscription {
		assert.False(t, device.IsBYOL, "LicenseType matches")
	} else {
		assert.True(t, device.IsBYOL, "LicenseType matches")
	}
	assert.ElementsMatch(t, resp.ACL, device.ACLs, "ACLs matches")
	assert.Equal(t, resp.SSHIPAddress, device.SSHIPAddress, "SSHIPAddress matches")
	assert.Equal(t, resp.SSHIPFqdn, device.SSHIPFqdn, "SSHIPFqdn matches")
	assert.Equal(t, resp.AccountNumber, device.AccountNumber, "AccountNumber matches")
	assert.ElementsMatch(t, resp.Notifications, device.Notifications, "Notifications matches")
	assert.Equal(t, resp.PurchaseOrderNumber, device.PurchaseOrderNumber, "PurchaseOrderNumber matches")
	assert.Equal(t, resp.RedundancyType, device.RedundancyType, "RedundancyType matches")
	assert.Equal(t, resp.RedundantUUID, device.RedundantUUID, "RedundantUUID matches")
	assert.Equal(t, resp.TermLength, device.TermLength, "TermLength matches")
	assert.Equal(t, resp.AdditionalBandwidth, device.AdditionalBandwidth, "AdditionalBandwidth matches")
	assert.Equal(t, resp.OrderReference, device.OrderReference, "OrderReference matches")
	assert.Equal(t, resp.InterfaceCount, device.InterfaceCount, "InterfaceCount matches")
	assert.Equal(t, resp.Core.Core, device.CoreCount, "Core.Core matches")
	if resp.DeviceManagementType == deviceManagementTypeEquinix {
		assert.False(t, device.IsSelfManaged, "DeviceManagementType matches")
	} else {
		assert.True(t, device.IsSelfManaged, "DeviceManagementType matches")
	}
	assert.Equal(t, len(resp.Interfaces), len(device.Interfaces), "Number of interfaces matches")
	for i := range resp.Interfaces {
		verifyDeviceInterface(t, device.Interfaces[i], resp.Interfaces[i])
	}
	assert.Equal(t, resp.VendorConfig, device.VendorConfiguration, "VendorConfigurations match")
}

func verifyDeviceInterface(t *testing.T, inf DeviceInterface, apiInf api.DeviceInterface) {
	assert.Equal(t, apiInf.ID, inf.ID, "ID matches")
	assert.Equal(t, apiInf.Name, inf.Name, "Name matches")
	assert.Equal(t, apiInf.Status, inf.Status, "Status matches")
	assert.Equal(t, apiInf.OperationalStatus, inf.OperationalStatus, "OperationalStatus matches")
	assert.Equal(t, apiInf.MACAddress, inf.MACAddress, "MACAddress matches")
	assert.Equal(t, apiInf.IPAddress, inf.IPAddress, "IPAddress matches")
	assert.Equal(t, apiInf.AssignedType, inf.AssignedType, "AssignedType matches")
	assert.Equal(t, apiInf.Type, inf.Type, "Type matches")
}

func verifyDeviceRequest(t *testing.T, device Device, req api.DeviceRequest) {
	assert.Equal(t, device.Throughput, req.Throughput, "Throughput matches")
	assert.Equal(t, device.ThroughputUnit, req.ThroughputUnit, "ThroughputUnit matches")
	assert.Equal(t, device.MetroCode, req.MetroCode, "MetroCode matches")
	assert.Equal(t, device.TypeCode, req.DeviceTypeCode, "TypeCode matches")
	assert.Equal(t, strconv.Itoa(device.TermLength), req.TermLength, "TermLength matches")
	if device.IsBYOL {
		assert.Equal(t, deviceLicenseModeBYOL, req.LicenseMode, "LicenseMode matches")
	} else {
		assert.Equal(t, deviceLicenseModeSubscription, req.LicenseMode, "LicenseMode matches")
	}
	assert.Equal(t, device.LicenseToken, req.LicenseToken, "LicenseToken matches")
	assert.Equal(t, device.PackageCode, req.PackageCode, "PackageCode matches")
	assert.Equal(t, device.Name, req.VirtualDeviceName, "Name matches")
	assert.ElementsMatch(t, device.Notifications, req.Notifications, "Notifications matches")
	assert.Equal(t, device.HostName, req.HostNamePrefix, "HostName matches")
	assert.Equal(t, device.OrderReference, req.OrderReference, "OrderReference matches")
	assert.Equal(t, device.PurchaseOrderNumber, req.PurchaseOrderNumber, "PurchaseOrderNumber matches")
	assert.Equal(t, device.AccountNumber, req.AccountNumber, "AccountNumber matches")
	assert.Equal(t, device.Version, req.Version, "Version matches")
	assert.Equal(t, device.InterfaceCount, req.InterfaceCount, "InterfaceCount matches")
	if device.IsSelfManaged {
		assert.Equal(t, deviceManagementTypeSelf, req.DeviceManagementType, "DeviceManagementType matches")
	} else {
		assert.Equal(t, deviceManagementTypeEquinix, req.DeviceManagementType, "DeviceManagementType matches")
	}
	assert.Equal(t, device.CoreCount, req.Core, "Core matches")
	assert.Equal(t, device.AdditionalBandwidth, req.AdditionalBandwidth, "AdditionalBandwidth matches")
	assert.ElementsMatch(t, mapDeviceACLsToFQDNACLs(device.ACLs), req.FqdnACL, "ACLs matches")
	assert.Equal(t, device.VendorConfiguration, req.VendorConfig, "VendorConfigurations match")
}

func verifyRedundantDeviceRequest(t *testing.T, primary, secondary Device, req api.DeviceRequest) {
	verifyDeviceRequest(t, primary, req)
	assert.Equal(t, secondary.MetroCode, req.Secondary.MetroCode, "Secondary MetroCode matches")
	assert.Equal(t, secondary.Name, req.Secondary.VirtualDeviceName, "Secondary Name matches")
	assert.ElementsMatch(t, secondary.Notifications, req.Secondary.Notifications, "Secondary Notifications matches")
	assert.Equal(t, secondary.HostName, req.Secondary.HostNamePrefix, "Secondary HostName matches")
	assert.Equal(t, secondary.AccountNumber, req.Secondary.AccountNumber, "Secondary AccountNumber matches")
	assert.Equal(t, secondary.AdditionalBandwidth, req.Secondary.AdditionalBandwidth, "Secondary AdditionalBandwidth matches")
	assert.ElementsMatch(t, mapDeviceACLsToFQDNACLs(secondary.ACLs), req.Secondary.FqdnACL, "Secondary ACLs matches")
	assert.Equal(t, secondary.VendorConfiguration, req.Secondary.VendorConfig, "Secondary VendorConfigurations match")
}
