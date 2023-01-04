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
	AdditionalBandwidth: Int(100),
	TypeCode:            String("PA-VM"),
	HostName:            String("myhostSRmy"),
	IsBYOL:              Bool(true),
	LicenseToken:        String("somelicensetokenaaaaazzzzz"),
	LicenseFileID:       String("8d180057-8309-4c59-b645-f630f010ad43"),
	CloudInitFileID:     String("9318885d-4b8c-48a5-9aa4-24387834ebae"),
	MetroCode:           String("SV"),
	Notifications:       []string{"test1@example.com", "test2@example.com"},
	PackageCode:         String("VM100"),
	TermLength:          Int(24),
	Throughput:          Int(1),
	ThroughputUnit:      String("Gbps"),
	Name:                String("PaloAltoSRmy"),
	ACLTemplateUUID:     String("4792d9ab-b8aa-49cc-8fe2-b56ced6c9c2f"),
	AccountNumber:       String("1777643"),
	OrderReference:      String("orderRef"),
	PurchaseOrderNumber: String("PO123456789"),
	InterfaceCount:      Int(10),
	CoreCount:           Int(2),
	Version:             String("10.09.05"),
	IsSelfManaged:       Bool(true),
	VendorConfiguration: map[string]string{
		"serialNumber": "12312312",
		"controller1":  "1.1.1.1",
	},
	UserPublicKey: &DeviceUserPublicKey{
		Username: String("testUserName"),
		KeyName:  String("testKey"),
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
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/devices", baseURL),
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
		MetroCode:           String("DC"),
		LicenseToken:        String("licenseToken"),
		LicenseFileID:       String("5a1102c6-d556-4498-b7ca-a10e902ef783"),
		Name:                String("secondary"),
		Notifications:       []string{"secondary@secondary.com"},
		HostName:            String("secondaryHostname"),
		AccountNumber:       String("99999"),
		AdditionalBandwidth: Int(200),
		ACLTemplateUUID:     String("4972e8d2-417f-4821-91a8-f4a61a6dcdc3"),
		MgmtAclTemplateUuid: String("4972e8d2-417f-4821-91a8-f4a61a6dcdc3"),
		VendorConfiguration: map[string]string{
			"serialNumber": "2222222",
			"controller1":  "2.2.2.2",
		},
		UserPublicKey: &DeviceUserPublicKey{
			Username: String("testUserSec"),
			KeyName:  String("testKeySec"),
		}}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/devices", baseURL),
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

func TestCreateClusterDevice(t *testing.T) {
	//given
	resp := api.DeviceRequestResponse{}
	if err := readJSONData("./test-fixtures/ne_device_create_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	device := Device{
		Name:                String("PANW-cluster"),
		MetroCode:           String("SV"),
		TypeCode:            String("PA-VM"),
		IsSelfManaged:       Bool(true),
		IsBYOL:              Bool(true),
		PackageCode:         String("VM100"),
		Notifications:       []string{"test1@example.com", "test2@example.com"},
		HostName:            String("panwHostName"),
		TermLength:          Int(24),
		AccountNumber:       String("177643"),
		Version:             String("10.1.3"),
		InterfaceCount:      Int(10),
		CoreCount:           Int(2),
		ACLTemplateUUID:     String("4972e8d2-417f-4821-91a8-f4a61a6dcdc3"),
		MgmtAclTemplateUuid: String("4972e8d2-417f-4821-91a8-f4a61a6dcdc3"),
		UserPublicKey: &DeviceUserPublicKey{
			Username: String("testUserName"),
			KeyName:  String("testKey"),
		},
		ClusterDetails: &ClusterDetails{
			ClusterName: String("clusterName"),
			Node0: &ClusterNodeDetail{
				VendorConfiguration: map[string]string{
					"hostname": "panw-host0",
				},
				LicenseFileId: String("8d180057-8309-4c59-b645-f630f010ad43"),
				LicenseToken:  String("licenseToken"),
			},
			Node1: &ClusterNodeDetail{
				VendorConfiguration: map[string]string{
					"hostname": "panw-host1",
				},
				LicenseFileId: String("8d180057-8309-4c59-b645-f630f010ad43"),
				LicenseToken:  String("licenseToken"),
			},
		},
	}
	req := api.DeviceRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/devices", baseURL),
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
	verifyClusterDeviceRequest(t, device, req)
}

func TestGetDevice(t *testing.T) {
	//given
	resp := api.Device{}
	if err := readJSONData("./test-fixtures/ne_device_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	devID := "myDevice"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/devices/%s", baseURL, devID), 200, resp)
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
	limit := respBody.Pagination.Limit
	statuses := []string{"INITIALIZING", "PROVISIONING"}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/devices?limit=%d&status=%s", baseURL, limit, url.QueryEscape("INITIALIZING,PROVISIONING")),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	devices, err := c.GetDevices(statuses)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, devices, "Client should return a response")
	assert.Equal(t, len(respBody.Data), len(devices), "Number of objects matches")
	for i := range respBody.Data {
		verifyDevice(t, devices[i], respBody.Data[i])
	}
}

func TestGetDeviceAdditionalBandwidthDetails(t *testing.T) {
	//given
	resp := api.DeviceAdditionalBandwidthResponse{}
	if err := readJSONData("./test-fixtures/ne_device_additionalbandwidth_get.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	devID := "myDevice"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/devices/%s/additionalBandwidths", baseURL, devID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	details, err := c.GetDeviceAdditionalBandwidthDetails(devID)

	//then
	assert.NotNil(t, details, "Returned additional bandwidth details struct is not nil")
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, resp.AdditionalBandwidth, details.AdditionalBandwidth, "AdditionalBandwidth matches")
	assert.Equal(t, resp.Status, details.Status, "Status matches")
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
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/devices/%s", baseURL, devID),
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
	assert.Equal(t, &newName, req.VirtualDeviceName, "DeviceName matches")
	assert.ElementsMatch(t, newNotifications, req.Notifications, "Notifications match")
	assert.Equal(t, &newTermLength, req.TermLength, "TermLength match")
}

func TestUpdateDeviceACLTemplate(t *testing.T) {
	//given
	devID := "myDevice"
	newACLTemplateID := "0647398e-2827-43cb-8fee-e6a9010ba78d"
	testHc := &http.Client{}
	req := api.DeviceACLTemplateRequest{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/devices/%s/acl", baseURL, devID),
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
	err := c.NewDeviceUpdateRequest(devID).WithACLTemplate(newACLTemplateID).Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, &newACLTemplateID, req.TemplateUUID, "ACLTemplateUUID matches")
}

func TestUpdateDeviceAdditionalBandwidth(t *testing.T) {
	//given
	devID := "myDevice"
	newBandwidth := 1000
	testHc := &http.Client{}
	req := api.DeviceAdditionalBandwidthUpdateRequest{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/devices/%s/additionalBandwidths", baseURL, devID),
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
	assert.Equal(t, &newBandwidth, req.AdditionalBandwidth, "AdditionalBandwidth match")
}

func TestDeleteDevice(t *testing.T) {
	//given
	devID := "myDevice"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/devices/%s", baseURL, devID),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteDevice(devID)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifyDevice(t *testing.T, device Device, resp api.Device) {
	assert.Equal(t, resp.UUID, device.UUID, "UUID matches")
	assert.Equal(t, resp.Name, device.Name, "Name matches")
	assert.Equal(t, resp.DeviceTypeCode, device.TypeCode, "DeviceTypeCode matches")
	assert.Equal(t, resp.Status, device.Status, "Status matches")
	assert.Equal(t, resp.LicenseStatus, device.LicenseStatus, "LicenseStatus matches")
	assert.Equal(t, resp.LicenseToken, device.LicenseToken, "LicenseToken matches")
	assert.Equal(t, resp.LicenseFileID, device.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, resp.MetroCode, device.MetroCode, "MetroCode matches")
	assert.Equal(t, resp.IBX, device.IBX, "IBX matches")
	assert.Equal(t, resp.Region, device.Region, "Region matches")
	assert.Equal(t, resp.Throughput, device.Throughput, "Throughput matches")
	assert.Equal(t, resp.ThroughputUnit, device.ThroughputUnit, "ThroughputUnit matches")
	assert.Equal(t, resp.HostName, device.HostName, "HostName matches")
	assert.Equal(t, resp.PackageCode, device.PackageCode, "PackageCode matches")
	assert.Equal(t, resp.Version, device.Version, "Version matches")
	if *resp.LicenseType == DeviceLicenseModeSubscription {
		assert.False(t, *device.IsBYOL, "LicenseType matches")
	} else {
		assert.True(t, *device.IsBYOL, "LicenseType matches")
	}
	assert.Equal(t, resp.ACLTemplateUUID, device.ACLTemplateUUID, "ACLTemplateUUID matches")
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
	if *resp.DeviceManagementType == DeviceManagementTypeEquinix {
		assert.False(t, *device.IsSelfManaged, "DeviceManagementType matches")
	} else {
		assert.True(t, *device.IsSelfManaged, "DeviceManagementType matches")
	}
	assert.Equal(t, len(resp.Interfaces), len(device.Interfaces), "Number of interfaces matches")
	for i := range resp.Interfaces {
		verifyDeviceInterface(t, device.Interfaces[i], resp.Interfaces[i])
	}
	assert.Equal(t, resp.VendorConfig, device.VendorConfiguration, "VendorConfigurations match")
	assert.NotNil(t, device.UserPublicKey, "UserPublicKey is not nil")
	verifyDeviceUserPublicKey(t, *device.UserPublicKey, *resp.UserPublicKey)
	assert.Equal(t, resp.ASN, device.ASN, "ASN matches")
	assert.Equal(t, resp.ZoneCode, device.ZoneCode, "ZoneCode matches")
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
	termLengthStr := strconv.Itoa(*device.TermLength)
	assert.Equal(t, &termLengthStr, req.TermLength, "TermLength matches")
	if *device.IsBYOL {
		assert.Equal(t, DeviceLicenseModeBYOL, StringValue(req.LicenseMode), "LicenseMode matches")
	} else {
		assert.Equal(t, DeviceLicenseModeSubscription, StringValue(req.LicenseMode), "LicenseMode matches")
	}
	assert.Equal(t, device.LicenseToken, req.LicenseToken, "LicenseToken matches")
	assert.Equal(t, device.LicenseFileID, req.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, device.CloudInitFileID, req.CloudInitFileID, "CloudInitFileID matches")
	assert.Equal(t, device.PackageCode, req.PackageCode, "PackageCode matches")
	assert.Equal(t, device.Name, req.VirtualDeviceName, "Name matches")
	assert.ElementsMatch(t, device.Notifications, req.Notifications, "Notifications matches")
	assert.Equal(t, device.HostName, req.HostNamePrefix, "HostName matches")
	assert.Equal(t, device.OrderReference, req.OrderReference, "OrderReference matches")
	assert.Equal(t, device.PurchaseOrderNumber, req.PurchaseOrderNumber, "PurchaseOrderNumber matches")
	assert.Equal(t, device.AccountNumber, req.AccountNumber, "AccountNumber matches")
	assert.Equal(t, device.Version, req.Version, "Version matches")
	assert.Equal(t, device.InterfaceCount, req.InterfaceCount, "InterfaceCount matches")
	if *device.IsSelfManaged {
		assert.Equal(t, DeviceManagementTypeSelf, StringValue(req.DeviceManagementType), "DeviceManagementType matches")
	} else {
		assert.Equal(t, DeviceManagementTypeEquinix, StringValue(req.DeviceManagementType), "DeviceManagementType matches")
	}
	assert.Equal(t, device.CoreCount, req.Core, "Core matches")
	assert.Equal(t, device.AdditionalBandwidth, req.AdditionalBandwidth, "AdditionalBandwidth matches")
	assert.Equal(t, device.ACLTemplateUUID, req.ACLTemplateUUID, "ACLTemplateUUID matches")
	assert.Equal(t, device.VendorConfiguration, req.VendorConfig, "VendorConfigurations match")
	assert.NotNil(t, req.UserPublicKey, "UserPublicKey is not nil")
	verifyDeviceUserPublicKeyRequest(t, *device.UserPublicKey, *req.UserPublicKey)
}

func verifyRedundantDeviceRequest(t *testing.T, primary, secondary Device, req api.DeviceRequest) {
	verifyDeviceRequest(t, primary, req)
	assert.Equal(t, secondary.MetroCode, req.Secondary.MetroCode, "Secondary MetroCode matches")
	assert.Equal(t, secondary.LicenseToken, req.Secondary.LicenseToken, "LicenseFileID matches")
	assert.Equal(t, secondary.LicenseFileID, req.Secondary.LicenseFileID, "LicenseFileID matches")
	assert.Equal(t, secondary.Name, req.Secondary.VirtualDeviceName, "Secondary Name matches")
	assert.ElementsMatch(t, secondary.Notifications, req.Secondary.Notifications, "Secondary Notifications matches")
	assert.Equal(t, secondary.HostName, req.Secondary.HostNamePrefix, "Secondary HostName matches")
	assert.Equal(t, secondary.AccountNumber, req.Secondary.AccountNumber, "Secondary AccountNumber matches")
	assert.Equal(t, secondary.AdditionalBandwidth, req.Secondary.AdditionalBandwidth, "Secondary AdditionalBandwidth matches")
	assert.Equal(t, secondary.ACLTemplateUUID, req.Secondary.ACLTemplateUUID, "Secondary ACLTemplateUUID matches")
	assert.Equal(t, secondary.MgmtAclTemplateUuid, req.Secondary.MgmtAclTemplateUUID, "Secondary MgmtAclTemplateUuid matches")
	assert.Equal(t, secondary.VendorConfiguration, req.Secondary.VendorConfig, "Secondary VendorConfigurations match")
	assert.NotNil(t, req.Secondary.UserPublicKey, "UserPublicKey is not nil")
	verifyDeviceUserPublicKeyRequest(t, *secondary.UserPublicKey, *req.Secondary.UserPublicKey)
}

func verifyClusterDeviceRequest(t *testing.T, device Device, req api.DeviceRequest) {
	assert.Equal(t, device.Name, req.VirtualDeviceName, "Name matches")
	assert.Equal(t, device.MetroCode, req.MetroCode, "MetroCode matches")
	assert.Equal(t, device.TypeCode, req.DeviceTypeCode, "TypeCode matches")
	if *device.IsSelfManaged {
		assert.Equal(t, DeviceManagementTypeSelf, StringValue(req.DeviceManagementType), "DeviceManagementType matches")
	} else {
		assert.Equal(t, DeviceManagementTypeEquinix, StringValue(req.DeviceManagementType), "DeviceManagementType matches")
	}
	if *device.IsBYOL {
		assert.Equal(t, DeviceLicenseModeBYOL, StringValue(req.LicenseMode), "LicenseMode matches")
	} else {
		assert.Equal(t, DeviceLicenseModeSubscription, StringValue(req.LicenseMode), "LicenseMode matches")
	}
	assert.Equal(t, device.PackageCode, req.PackageCode, "PackageCode matches")
	assert.ElementsMatch(t, device.Notifications, req.Notifications, "Notifications matches")
	assert.Equal(t, device.HostName, req.HostNamePrefix, "HostName matches")
	termLengthStr := strconv.Itoa(*device.TermLength)
	assert.Equal(t, &termLengthStr, req.TermLength, "TermLength matches")
	assert.Equal(t, device.AccountNumber, req.AccountNumber, "AccountNumber matches")
	assert.Equal(t, device.Version, req.Version, "Version matches")
	assert.Equal(t, device.InterfaceCount, req.InterfaceCount, "InterfaceCount matches")
	assert.Equal(t, device.CoreCount, req.Core, "Core matches")
	assert.Equal(t, device.ACLTemplateUUID, req.ACLTemplateUUID, "ACLTemplateUUID matches")
	assert.Equal(t, device.MgmtAclTemplateUuid, req.MgmtAclTemplateUUID, "MgmtAclTemplateUuid matches")
	assert.NotNil(t, req.UserPublicKey, "UserPublicKey is not nil")
	verifyDeviceUserPublicKeyRequest(t, *device.UserPublicKey, *req.UserPublicKey)
	assert.NotNil(t, req.ClusterDetails, "ClusterDetails are not nil")
	verifyClusterDetailsRequest(t, *device.ClusterDetails, *req.ClusterDetails)
}

func verifyDeviceUserPublicKey(t *testing.T, userKey DeviceUserPublicKey, apiUserKey api.DeviceUserPublicKey) {
	assert.Equal(t, apiUserKey.Username, userKey.Username, "Username matches")
	assert.Equal(t, apiUserKey.KeyName, userKey.KeyName, "KeyName matches")
}

func verifyDeviceUserPublicKeyRequest(t *testing.T, userKey DeviceUserPublicKey, apiUserKeyReq api.DeviceUserPublicKeyRequest) {
	assert.Equal(t, apiUserKeyReq.Username, userKey.Username, "Username matches")
	assert.Equal(t, apiUserKeyReq.KeyName, userKey.KeyName, "KeyName matches")
}

func verifyClusterDetailsRequest(t *testing.T, clusterDetails ClusterDetails, apiClusterDetailsReq api.ClusterDetailsRequest) {
	assert.Equal(t, clusterDetails.ClusterName, apiClusterDetailsReq.ClusterName, "ClusterName matches")
	apiClusterNodeDetailReqMap := apiClusterDetailsReq.ClusterNodeDetails
	assert.NotNil(t, apiClusterNodeDetailReqMap, "ClusterNodeDetails are not nil")
	verifyClusterNodeDetailRequest(t, clusterDetails.Node0, apiClusterNodeDetailReqMap["node0"])
	verifyClusterNodeDetailRequest(t, clusterDetails.Node1, apiClusterNodeDetailReqMap["node1"])
}

func verifyClusterNodeDetailRequest(t *testing.T, clusterNodeDetail *ClusterNodeDetail, apiClusterNodeDetailReq api.ClusterNodeDetailRequest) {
	assert.NotNil(t, apiClusterNodeDetailReq, "ClusterNodeDetailRequest is not nil")
	assert.Equal(t, clusterNodeDetail.VendorConfiguration, apiClusterNodeDetailReq.VendorConfiguration, "VendorConfigurations match")
	assert.Equal(t, clusterNodeDetail.LicenseFileId, apiClusterNodeDetailReq.LicenseFileID, "LicenseFileId matches")
	assert.Equal(t, clusterNodeDetail.LicenseToken, apiClusterNodeDetailReq.LicenseToken, "LicenseToken matches")
}
