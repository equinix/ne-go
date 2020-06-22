package ne

import (
	"fmt"
	"ne-go/v1/internal/api"
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
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

func (c RestClient) CreateDevice(device Device) (string, error) {
	url := fmt.Sprintf("%s/ne/v1/device", c.baseURL)
	reqBody := createDeviceRequest(device)
	respBody := api.VirtualDeviceCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)

	if err := c.execute(req, resty.MethodPost, url); err != nil {
		return "", err
	}
	return respBody.UUID, nil
}

func (c RestClient) GetDevice(uuid string) (*Device, error) {
	url := fmt.Sprintf("%s/ne/v1/device/%s", c.baseURL, url.PathEscape(uuid))
	result := api.VirtualDeviceDetailsResponse{}
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

func (c RestClient) NewDeviceUpdateRequest(uuid string) DeviceUpdateRequest {
	return &restDeviceUpdateRequest{
		uuid:         uuid,
		deviceFields: make(map[string]interface{}),
		c:            c}
}

func (c RestClient) DeleteDevice(uuid string) error {
	url := fmt.Sprintf("%s/ne/v1/device/%s", c.baseURL, url.PathEscape(uuid))
	if err := c.execute(c.R(), resty.MethodDelete, url); err != nil {
		return err
	}
	return nil
}

func (req *restDeviceUpdateRequest) WithDeviceName(deviceName string) DeviceUpdateRequest {
	req.deviceFields["deviceName"] = deviceName
	return req
}

func (req *restDeviceUpdateRequest) WithTermLength(termLength int) DeviceUpdateRequest {
	req.deviceFields["termLength"] = termLength
	return req
}

func (req *restDeviceUpdateRequest) WithNotifications(notifications []string) DeviceUpdateRequest {
	req.deviceFields["notifications"] = notifications
	return req
}

func (req *restDeviceUpdateRequest) WithAdditionalBandwidth(additionalBandwidth int) DeviceUpdateRequest {
	req.additionalBandwidth = additionalBandwidth
	return req
}

func (req *restDeviceUpdateRequest) WithACLs(acls []string) DeviceUpdateRequest {
	req.acls = acls
	return req
}

func (req *restDeviceUpdateRequest) Execute() error {
	updateErr := UpdateError{}
	if err := req.c.replaceDeviceFields(req.uuid, req.deviceFields); err != nil {
		updateErr.failed = append(updateErr.failed, ChangeError{
			Type:   ChangeTypeUpdate,
			Target: "deviceFields",
			Value:  req.deviceFields,
			Cause:  err})
	}

	if len(req.acls) > 0 {
		if err := req.c.replaceDeviceACLs(req.uuid, req.acls); err != nil {
			updateErr.failed = append(updateErr.failed, ChangeError{
				Type:   ChangeTypeUpdate,
				Target: "acl",
				Value:  req.acls,
				Cause:  err})
		}
	}

	if req.additionalBandwidth > 0 {
		if err := req.c.replaceDeviceAdditionalBandwidth(req.uuid, req.additionalBandwidth); err != nil {
			updateErr.failed = append(updateErr.failed, ChangeError{
				Type:   ChangeTypeUpdate,
				Target: "additionalBandwidth",
				Value:  req.additionalBandwidth,
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

func mapDeviceAPIToDomain(apiDevice api.VirtualDeviceDetailsResponse) (*Device, error) {
	dev := Device{}
	dev.AccountName = apiDevice.UUID
	dev.AccountNumber = apiDevice.AccountNumber
	dev.ACL = apiDevice.ACL
	dev.AdditionalBandwidth = int(apiDevice.AdditionalBandwidth)
	dev.Controller1 = apiDevice.Controller1
	dev.Controller2 = apiDevice.Controller2
	dev.DeviceSerialNo = apiDevice.DeviceSerialNo
	dev.DeviceTypeCategory = apiDevice.DeviceTypeCategory
	dev.DeviceTypeCode = apiDevice.DeviceTypeCode
	dev.DeviceTypeName = apiDevice.DeviceTypeName
	dev.DeviceTypeVendor = apiDevice.DeviceTypeVendor
	dev.Expiry = apiDevice.Expiry
	dev.HostName = apiDevice.HostName
	dev.LicenseFileID = apiDevice.LicenseFileID
	dev.LicenseKey = apiDevice.LicenseKey
	dev.LicenseName = apiDevice.LicenseName
	dev.LicenseSecret = apiDevice.LicenseSecret
	dev.LicenseStatus = apiDevice.LicenseStatus
	dev.LicenseType = apiDevice.LicenseType
	dev.LocalID = apiDevice.LocalID
	dev.ManagementGatewayIP = apiDevice.ManagementGatewayIP
	dev.ManagementIP = apiDevice.ManagementIP
	dev.MetroCode = apiDevice.MetroCode
	dev.MetroName = apiDevice.MetroName
	dev.Name = apiDevice.Name
	dev.Notifications = apiDevice.Notifications
	dev.PackageCode = apiDevice.PackageCode
	dev.PackageName = apiDevice.PackageName
	dev.PrimaryDNSName = apiDevice.PrimaryDNSName
	dev.PublicGatewayIP = apiDevice.PublicGatewayIP
	dev.PublicIP = apiDevice.PublicIP
	dev.PurchaseOrderNumber = apiDevice.PublicIP
	dev.RedundancyType = apiDevice.RedundancyType
	dev.RedundantUUID = apiDevice.RedundantUUID
	dev.Region = apiDevice.Region
	dev.RemoteID = apiDevice.RemoteID
	dev.SecondaryDNSName = apiDevice.SecondaryDNSName
	dev.SerialNumber = apiDevice.SerialNumber
	dev.SiteID = apiDevice.SiteID
	dev.SSHIPAddress = apiDevice.SSHIPAddress
	dev.SSHIPFqdn = apiDevice.SSHIPFqdn
	dev.Status = apiDevice.Status
	dev.SystemIPAddress = apiDevice.SystemIPAddress
	dev.TermLength = int(apiDevice.TermLength)
	if val, err := strconv.Atoi(apiDevice.Throughput); err == nil {
		dev.Throughput = val
	} else {
		return nil, fmt.Errorf("can't parse throughput: %v", err)
	}
	dev.ThroughputUnit = apiDevice.ThroughputUnit
	dev.UUID = apiDevice.UUID
	if apiDevice.VendorConfig != nil {
		dev.VendorConfig = mapDeviceVendorConfigAPIToDomain(*apiDevice.VendorConfig)
	}
	return &dev, nil
}

func createDeviceRequest(device Device) api.VirtualDeviceRequest {
	req := api.VirtualDeviceRequest{}
	req.AccountNumber = device.AccountNumber
	//seems to be not needed if account is already there
	//req.AccountReferenceID = device.AccountReferenceID
	req.ACL = device.ACL
	req.AdditionalBandwidth = int32(device.AdditionalBandwidth)
	req.DeviceTypeCode = &device.DeviceTypeCode
	req.HostNamePrefix = &device.HostName
	req.LicenseFileID = device.LicenseFileID
	req.LicenseKey = device.LicenseKey
	req.LicenseMode = &device.LicenseType
	req.LicenseSecret = device.LicenseSecret
	req.LicenseToken = device.LicenseToken
	req.MetroCode = &device.MetroCode
	req.Notifications = device.Notifications
	req.PackageCode = device.PackageCode
	//req.Secondary = device.Secondary
	req.SiteID = device.SiteID
	//req.SSHUsers = device.SSHUsers
	req.SystemIPAddress = device.SystemIPAddress
	req.Throughput = int32(device.Throughput)
	req.ThroughputUnit = device.ThroughputUnit
	req.VirtualDeviceName = &device.Name
	return req
}

func mapDeviceVendorConfigAPIToDomain(api api.VendorConfig) DeviceVendorConfig {
	return DeviceVendorConfig{
		SiteID:          api.SiteID,
		SystemIPAddress: api.SystemIPAddress,
	}
}

func (c RestClient) replaceDeviceACLs(uuid string, acls []string) error {
	url := fmt.Sprintf("%s/ne/v1/device/%s/acl", c.baseURL, url.PathEscape(uuid))
	req := c.R().SetBody(acls)
	if err := c.execute(req, resty.MethodPut, url); err != nil {
		return err
	}
	return nil
}

func (c RestClient) replaceDeviceAdditionalBandwidth(uuid string, bandwidth int) error {
	url := fmt.Sprintf("%s/ne/v1/device/additionalbandwidth/%s", c.baseURL, url.PathEscape(uuid))
	bandwidthConv := int32(bandwidth)
	reqBody := api.AdditionalBandwidthUpdateRequest{AdditionalBandwidth: &bandwidthConv}
	req := c.R().SetBody(reqBody)
	if err := c.execute(req, resty.MethodPut, url); err != nil {
		return err
	}
	return nil
}

func (c RestClient) replaceDeviceFields(uuid string, fields map[string]interface{}) error {
	reqBody := api.VirtualDeviceInternalPatchRequestDto{}
	okToSend := false
	if v, ok := fields["deviceName"]; ok {
		reqBody.VirtualDeviceName = v.(string)
		okToSend = true
	}
	if v, ok := fields["termLength"]; ok {
		reqBody.TermLength = int64(v.(int))
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
