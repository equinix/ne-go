package ne

import (
	"fmt"
)

type Client interface {
	CreateDevice(device Device) (string, error)
	//create redundant device ? or device create request ?
	GetDevice(uuid string) (*Device, error)
	NewDeviceUpdateRequest(uuid string) DeviceUpdateRequest
	DeleteDevice(uuid string) error

	CreateSSHUser(username string, password string, device string) (string, error)
	GetSSHUser(uuid string) (*SSHUser, error)
	NewSSHUserUpdateRequest(uuid string) SSHUserUpdateRequest
	//delete composite operation ?
}

type DeviceUpdateRequest interface {
	WithDeviceName(deviceName string) DeviceUpdateRequest
	WithTermLength(termLength int) DeviceUpdateRequest
	WithNotifications(notifications []string) DeviceUpdateRequest
	WithAdditionalBandwidth(additionalBandwidth int) DeviceUpdateRequest
	WithACLs(acls []string) DeviceUpdateRequest
	Execute() error
}

type SSHUserUpdateRequest interface {
	WithNewPassword(password string) SSHUserUpdateRequest
	WithNewDevices(uuids []string) SSHUserUpdateRequest
	WithRemovedDevices(uuids []string) SSHUserUpdateRequest
	Execute() error
}

//Error describes Network Edge error that occurs during API call processing
type Error struct {
	//ErrorCode is short error identifier
	ErrorCode string
	//ErrorMessage is textual description of an error
	ErrorMessage string
}

const (
	ChangeTypeCreate = "Add"
	ChangeTypeUpdate = "Update"
	ChangeTypeDelete = "Delete"
)

type ChangeError struct {
	Type   string
	Target string
	Value  interface{}
	Cause  error
}

func (e ChangeError) Error() string {
	return fmt.Sprintf("change type '%s', target '%s', value '%s', cause: '%s'", e.Type, e.Target, e.Value, e.Cause)
}

type UpdateError struct {
	failed []ChangeError
}

func (e UpdateError) Error() string {
	return fmt.Sprintf("update error: %d changes failed", len(e.failed))
}

type Device struct {
	AccountName         string
	AccountNumber       string
	ACL                 []string
	AdditionalBandwidth int
	Controller1         string
	Controller2         string
	DeviceSerialNo      string
	DeviceTypeCategory  string
	DeviceTypeCode      string
	DeviceTypeName      string
	DeviceTypeVendor    string
	Expiry              string
	HostName            string
	//Interfaces          []DeviceInterface
	LicenseFileID       string
	LicenseKey          string
	LicenseName         string
	LicenseSecret       string
	LicenseStatus       string
	LicenseType         string
	LicenseToken        string
	LocalID             string
	ManagementGatewayIP string
	ManagementIP        string
	MetroCode           string
	MetroName           string
	Name                string
	Notifications       []string
	PackageCode         string
	PackageName         string
	PrimaryDNSName      string
	PublicGatewayIP     string
	PublicIP            string
	PurchaseOrderNumber string
	RedundancyType      string
	RedundantUUID       string
	Region              string
	RemoteID            string
	SecondaryDNSName    string
	SerialNumber        string
	SiteID              string
	SSHIPAddress        string
	SSHIPFqdn           string
	Status              string
	SystemIPAddress     string
	TermLength          int
	Throughput          int
	ThroughputUnit      string
	UUID                string
	VendorConfig        DeviceVendorConfig
}

type DeviceInterface struct {
	Asn               int
	AssignedType      string
	ID                int
	IPAddress         string
	IPV4Mask          string
	IPV4Subnet        string
	MacAddress        string
	Name              string
	OperationalStatus string
	Status            string
	Type              string
}

type DeviceVendorConfig struct {
	SiteID          string
	SystemIPAddress string
	/* Controller1     string
	Controller2     string
	AdminPassword   string
	LocalID         string
	RemoteID        string
	SerialNumber    string
	LicenseKey string
	LicenseSecret string
	*/
}

type SSHUser struct {
	UUID        string
	Username    string
	Password    string
	MetroCodes  []string
	DeviceUUIDs []string
}
