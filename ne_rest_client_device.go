package ne

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/equinix/ne-go/internal/api"
	"github.com/go-resty/resty/v2"
)

const (
	deviceManagementTypeSelf      = "SELF-CONFIGURED"
	deviceManagementTypeEquinix   = "EQUINIX-CONFIGURED"
	deviceLicenseModeSubscription = "Sub"
	deviceLicenseModeBYOL         = "BYOL"
)

type restDeviceUpdateRequest struct {
	uuid                string
	deviceFields        map[string]interface{}
	deviceName          string
	termLength          int
	notifications       []string
	additionalBandwidth int
	acls                []string
	c                   RestClient
}

//CreateDevice creates given Network Edge device and returns its UUID upon successful creation
func (c RestClient) CreateDevice(device Device) (string, error) {
	url := fmt.Sprintf("%s/ne/v1/device", c.baseURL)
	reqBody := createDeviceRequest(device)
	respBody := api.DeviceRequestResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)

	if err := c.execute(req, resty.MethodPost, url); err != nil {
		return "", err
	}
	return respBody.UUID, nil
}

//CreateRedundantDevice creates HA device setup from given primary and secondary devices and
//returns their UUIDS upon successful creation
func (c RestClient) CreateRedundantDevice(primary Device, secondary Device) (string, string, error) {
	url := fmt.Sprintf("%s/ne/v1/device", c.baseURL)
	reqBody := createRedundantDeviceRequest(primary, secondary)
	respBody := api.DeviceRequestResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)

	if err := c.execute(req, resty.MethodPost, url); err != nil {
		return "", "", err
	}
	return respBody.UUID, respBody.SecondaryUUID, nil
}

//GetDevice fetches details of a device with a given UUID
func (c RestClient) GetDevice(uuid string) (*Device, error) {
	url := fmt.Sprintf("%s/ne/v1/device/%s", c.baseURL, url.PathEscape(uuid))
	result := api.Device{}
	request := c.R().SetResult(&result)
	if err := c.execute(request, resty.MethodGet, url); err != nil {
		return nil, err
	}
	device, err := mapDeviceAPIToDomain(result)
	if err != nil {
		return nil, fmt.Errorf("error when reading device data: %s", err)
	}
	return device, nil
}

//NewDeviceUpdateRequest creates new composite update request for a device with a given UUID
func (c RestClient) NewDeviceUpdateRequest(uuid string) DeviceUpdateRequest {
	return &restDeviceUpdateRequest{
		uuid:         uuid,
		deviceFields: make(map[string]interface{}),
		c:            c}
}

//DeleteDevice deletes device with a given UUID
func (c RestClient) DeleteDevice(uuid string) error {
	url := fmt.Sprintf("%s/ne/v1/device/%s", c.baseURL, url.PathEscape(uuid))
	req := c.R().SetQueryParam("deleteRedundantDevice", "true")
	if err := c.execute(req, resty.MethodDelete, url); err != nil {
		return err
	}
	return nil
}

//WithDeviceName sets new device name in a composite device update request
func (req *restDeviceUpdateRequest) WithDeviceName(deviceName string) DeviceUpdateRequest {
	req.deviceFields["deviceName"] = deviceName
	return req
}

//WithTermLength sets new term length in a composite device update request
func (req *restDeviceUpdateRequest) WithTermLength(termLength int) DeviceUpdateRequest {
	req.deviceFields["termLength"] = termLength
	return req
}

//WithNotifications sets new notifications in a composite device update request
func (req *restDeviceUpdateRequest) WithNotifications(notifications []string) DeviceUpdateRequest {
	req.deviceFields["notifications"] = notifications
	return req
}

//WithAdditionalBandwidth sets new additional bandwidth in a composite device update request
func (req *restDeviceUpdateRequest) WithAdditionalBandwidth(additionalBandwidth int) DeviceUpdateRequest {
	req.additionalBandwidth = additionalBandwidth
	return req
}

//WithAdditionalBandwidth sets new ACLs in a composite device update request
func (req *restDeviceUpdateRequest) WithACLs(acls []string) DeviceUpdateRequest {
	req.acls = acls
	return req
}

//Execute attempts to update device according new data set in composite update request.
//This is not atomic operation and if any update will fail, other changes won't be reverted.
//UpdateError will be returned if any of requested data failed to update
func (req *restDeviceUpdateRequest) Execute() error {
	updateErr := UpdateError{}
	if err := req.c.replaceDeviceFields(req.uuid, req.deviceFields); err != nil {
		updateErr.AddChangeError(changeTypeUpdate, "deviceFields", req.deviceFields, err)
	}
	if len(req.acls) > 0 {
		if err := req.c.replaceDeviceACLs(req.uuid, req.acls); err != nil {
			updateErr.AddChangeError(changeTypeUpdate, "acl", req.acls, err)
		}
	}
	if req.additionalBandwidth > 0 {
		if err := req.c.replaceDeviceAdditionalBandwidth(req.uuid, req.additionalBandwidth); err != nil {
			updateErr.AddChangeError(changeTypeUpdate, "additionalBandwidth", req.additionalBandwidth, err)
		}
	}
	if updateErr.ChangeErrorsCount() > 0 {
		return updateErr
	}
	return nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func mapDeviceAPIToDomain(apiDevice api.Device) (*Device, error) {
	device := Device{}
	device.UUID = apiDevice.UUID
	device.Name = apiDevice.Name
	device.TypeCode = apiDevice.DeviceTypeCode
	device.Status = apiDevice.Status
	device.LicenseStatus = apiDevice.LicenseStatus
	device.MetroCode = apiDevice.MetroCode
	device.IBX = apiDevice.IBX
	device.Region = apiDevice.Region
	if val, err := strconv.Atoi(apiDevice.Throughput); err == nil {
		device.Throughput = val
	} else {
		return nil, fmt.Errorf("can't parse throughput: %v", err)
	}
	device.ThroughputUnit = apiDevice.ThroughputUnit
	device.HostName = apiDevice.HostName
	device.PackageCode = apiDevice.PackageCode
	device.Version = apiDevice.Version
	if apiDevice.LicenseType == deviceLicenseModeBYOL {
		device.IsBYOL = true
	}
	device.ACLs = apiDevice.ACL
	device.SSHIPAddress = apiDevice.SSHIPAddress
	device.SSHIPFqdn = apiDevice.SSHIPFqdn
	device.AccountNumber = apiDevice.AccountNumber
	device.Notifications = apiDevice.Notifications
	device.PurchaseOrderNumber = apiDevice.PurchaseOrderNumber
	device.RedundancyType = apiDevice.RedundancyType
	device.RedundantUUID = apiDevice.RedundantUUID
	device.TermLength = apiDevice.TermLength
	device.AdditionalBandwidth = apiDevice.AdditionalBandwidth
	device.OrderReference = apiDevice.OrderReference
	device.InterfaceCount = apiDevice.InterfaceCount
	if apiDevice.Core != nil {
		device.CoreCount = apiDevice.Core.Core
	}
	if apiDevice.DeviceManagementType == deviceManagementTypeSelf {
		device.IsSelfManaged = true
	}
	device.Interfaces = mapDeviceInterfacesAPIToDomain(apiDevice.Interfaces)
	device.VendorConfiguration = apiDevice.VendorConfig
	return &device, nil
}

func mapDeviceInterfacesAPIToDomain(apiInterfaces []api.DeviceInterface) []DeviceInterface {
	transformed := make([]DeviceInterface, len(apiInterfaces))
	for i := range apiInterfaces {
		transformed[i] = DeviceInterface{
			ID:                apiInterfaces[i].ID,
			Name:              apiInterfaces[i].Name,
			Status:            apiInterfaces[i].Status,
			OperationalStatus: apiInterfaces[i].OperationalStatus,
			MACAddress:        apiInterfaces[i].MACAddress,
			IPAddress:         apiInterfaces[i].IPAddress,
			AssignedType:      apiInterfaces[i].AssignedType,
			Type:              apiInterfaces[i].Type,
		}
	}
	return transformed
}

func createDeviceRequest(device Device) api.DeviceRequest {
	req := api.DeviceRequest{}
	if device.Throughput > 0 {
		req.Throughput = strconv.Itoa(device.Throughput)
	}
	req.ThroughputUnit = device.ThroughputUnit
	req.MetroCode = device.MetroCode
	req.DeviceTypeCode = device.TypeCode
	req.TermLength = strconv.Itoa(device.TermLength)
	req.LicenseMode = deviceLicenseModeSubscription
	if device.IsBYOL {
		req.LicenseMode = deviceLicenseModeBYOL
	}
	req.LicenseToken = device.LicenseToken
	req.PackageCode = device.PackageCode
	req.VirtualDeviceName = device.Name
	req.Notifications = device.Notifications
	req.HostNamePrefix = device.HostName
	req.OrderReference = device.OrderReference
	req.PurchaseOrderNumber = device.PurchaseOrderNumber
	req.AccountNumber = device.AccountNumber
	req.Version = device.Version
	req.InterfaceCount = device.InterfaceCount
	req.DeviceManagementType = deviceManagementTypeEquinix
	if device.IsSelfManaged {
		req.DeviceManagementType = deviceManagementTypeSelf
	}
	req.Core = device.CoreCount
	if device.AdditionalBandwidth > 0 {
		req.AdditionalBandwidth = strconv.Itoa(device.AdditionalBandwidth)
	}
	req.FqdnACL = mapDeviceACLsToFQDNACLs(device.ACLs)
	req.VendorConfig = device.VendorConfiguration
	return req
}

func mapDeviceACLsToFQDNACLs(acls []string) []api.DeviceFqdnACL {
	transformed := make([]api.DeviceFqdnACL, len(acls))
	for i := range acls {
		transformed[i] = api.DeviceFqdnACL{
			CIDRs: []string{acls[i]},
			Type:  "SUBNET",
		}
	}
	return transformed
}

func createRedundantDeviceRequest(primary Device, secondary Device) api.DeviceRequest {
	req := createDeviceRequest(primary)
	secReq := api.SecondaryDeviceRequest{}
	secReq.MetroCode = secondary.MetroCode
	secReq.VirtualDeviceName = secondary.Name
	secReq.Notifications = secondary.Notifications
	secReq.HostNamePrefix = secondary.HostName
	secReq.AccountNumber = secondary.AccountNumber
	if secondary.AdditionalBandwidth > 0 {
		secReq.AdditionalBandwidth = strconv.Itoa(secondary.AdditionalBandwidth)
	}
	secReq.FqdnACL = mapDeviceACLsToFQDNACLs(secondary.ACLs)
	secReq.VendorConfig = secondary.VendorConfiguration
	req.Secondary = &secReq
	return req
}

func (c RestClient) replaceDeviceACLs(uuid string, acls []string) error {
	url := fmt.Sprintf("%s/ne/v1/device/%s/fqdn-acl", c.baseURL, url.PathEscape(uuid))
	reqBody := mapDeviceACLsToFQDNACLs(acls)
	req := c.R().SetBody(reqBody)
	if err := c.execute(req, resty.MethodPut, url); err != nil {
		return err
	}
	return nil
}

func (c RestClient) replaceDeviceAdditionalBandwidth(uuid string, bandwidth int) error {
	url := fmt.Sprintf("%s/ne/v1/device/additionalbandwidth/%s", c.baseURL, url.PathEscape(uuid))
	reqBody := api.DeviceAdditionalBandwidthUpdateRequest{AdditionalBandwidth: bandwidth}
	req := c.R().SetBody(reqBody)
	if err := c.execute(req, resty.MethodPut, url); err != nil {
		return err
	}
	return nil
}

func (c RestClient) replaceDeviceFields(uuid string, fields map[string]interface{}) error {
	reqBody := api.DeviceUpdateRequest{}
	okToSend := false
	if v, ok := fields["deviceName"]; ok {
		reqBody.VirtualDeviceName = v.(string)
		okToSend = true
	}
	if v, ok := fields["termLength"]; ok {
		reqBody.TermLength = v.(int)
		okToSend = true
	}
	if v, ok := fields["notifications"]; ok {
		reqBody.Notifications = v.([]string)
		okToSend = true
	}
	if okToSend {
		url := fmt.Sprintf("%s/ne/v1/device/%s", c.baseURL, url.PathEscape(uuid))
		req := c.R().SetBody(&reqBody)
		if err := c.execute(req, resty.MethodPatch, url); err != nil {
			return err
		}
	}
	return nil
}
