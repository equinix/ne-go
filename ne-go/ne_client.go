//Package ne implements Network Edge client
package ne

import (
	"fmt"
)

//Client interface describes operations provided by Network Edge client library
type Client interface {
	CreateDevice(device Device) (string, error)
	CreateRedundantDevice(primary Device, secondary Device) (string, string, error)
	GetDevice(uuid string) (*Device, error)
	NewDeviceUpdateRequest(uuid string) DeviceUpdateRequest
	DeleteDevice(uuid string) error

	CreateSSHUser(username string, password string, device string) (string, error)
	GetSSHUser(uuid string) (*SSHUser, error)
	NewSSHUserUpdateRequest(uuid string) SSHUserUpdateRequest
	DeleteSSHUser(uuid string) error
}

//DeviceUpdateRequest describes composite request to update given Network Edge device
type DeviceUpdateRequest interface {
	WithDeviceName(deviceName string) DeviceUpdateRequest
	WithTermLength(termLength int) DeviceUpdateRequest
	WithNotifications(notifications []string) DeviceUpdateRequest
	WithAdditionalBandwidth(additionalBandwidth int) DeviceUpdateRequest
	WithACLs(acls []string) DeviceUpdateRequest
	Execute() error
}

//SSHUserUpdateRequest describes composite request to update given Network Edge SSH user
type SSHUserUpdateRequest interface {
	WithNewPassword(password string) SSHUserUpdateRequest
	WithDeviceChange(old []string, new []string) SSHUserUpdateRequest
	Execute() error
}

//Error describes Network Edge error that occurs during API call processing
type Error struct {
	//ErrorCode is short error identifier
	ErrorCode string
	//ErrorMessage is textual description of an error
	ErrorMessage string
}

//ChangeError describes single error that occured during update of selected target property
type ChangeError struct {
	Type   string
	Target string
	Value  interface{}
	Cause  error
}

func (e ChangeError) Error() string {
	return fmt.Sprintf("change type '%s', target '%s', value '%s', cause: '%s'", e.Type, e.Target, e.Value, e.Cause)
}

//UpdateError describes error that occured during composite update request and consists of multiple atomic change errors
type UpdateError struct {
	Failed []ChangeError
}

//AddChangeError functions add new atomic change error to update error structure
func (e *UpdateError) AddChangeError(changeType string, target string, value interface{}, cause error) {
	e.Failed = append(e.Failed, ChangeError{
		Type:   changeType,
		Target: target,
		Value:  value,
		Cause:  cause})
}

//ChangeErrorsCount returns number of atomic change errors in a given composite update error
func (e UpdateError) ChangeErrorsCount() int {
	return len(e.Failed)
}

func (e UpdateError) Error() string {
	str := fmt.Sprintf("update error: %d changes failed.", len(e.Failed))
	for _, err := range e.Failed {
		str = fmt.Sprintf("%s [%s]", str, err.Error())
	}
	return str
}

//Device describes Network Edge device
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
	VendorConfig        *DeviceVendorConfig
	Version             string
}

//DeviceVendorConfig describes vendor specific configuration attrubues of a Network Edge device
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

//SSHUser describes Network Edge SSH user
type SSHUser struct {
	UUID        string
	Username    string
	Password    string
	MetroCodes  []string
	DeviceUUIDs []string
}
