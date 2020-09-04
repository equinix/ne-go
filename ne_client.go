//Package ne implements Network Edge client
package ne

import (
	"fmt"
)

const (
	//DeviceStateInitializing Equinix is allocating resources and creating new device
	DeviceStateInitializing = "INITIALIZING"
	//DeviceStateProvisioning Network Edge device is booting
	DeviceStateProvisioning = "PROVISIONING"
	//DeviceStateWaitingPrimary secondary Network Edge device is waiting for provisioning of its redundant (primary) device
	DeviceStateWaitingPrimary = "WAITING_FOR_PRIMARY"
	//DeviceStateWaitingSecondary primary Network Edge device is waiting for provisioning of its redundant (secondary) device
	DeviceStateWaitingSecondary = "WAITING_FOR_SECONDARY"
	//DeviceStateFailed Network Edge device creation and provisioning have failed
	DeviceStateFailed = "FAILED"
	//DeviceStateProvisioned Network Edge device was successfully provisioned and is fully operational
	DeviceStateProvisioned = "PROVISIONED"
	//DeviceStateDeprovisioning Network Edge device is in process of deprovisioning
	DeviceStateDeprovisioning = "DEPROVISIONING"
	//DeviceStateDeprovisioned Network Edge device was successfully deprovisioned
	DeviceStateDeprovisioned = "DEPROVISIONED"

	//DeviceLicenseStateApplying license is in registration process
	DeviceLicenseStateApplying = "APPLYING_LICENSE"
	//DeviceLicenseStateRegistered license was successfully registered
	DeviceLicenseStateRegistered = "REGISTERED"
	//DeviceLicenseStateFailed license registration has failed
	DeviceLicenseStateFailed = "REGISTRATION_FAILED"

	//ErrorCodeDeviceRemoved is used on attempt to remove device that is deprovisioning or already deprovisioned
	ErrorCodeDeviceRemoved = "IC-NE-VD-030"
)

//Client interface describes operations provided by Network Edge client library
type Client interface {
	CreateDevice(device Device) (string, error)
	CreateRedundantDevice(primary Device, secondary Device) (string, string, error)
	GetDevice(uuid string) (*Device, error)
	GetDevices(statuses []string) ([]Device, error)
	NewDeviceUpdateRequest(uuid string) DeviceUpdateRequest
	DeleteDevice(uuid string) error
	GetDeviceACLs(uuid string) (*DeviceACLs, error)

	CreateSSHUser(username string, password string, device string) (string, error)
	GetSSHUsers() ([]SSHUser, error)
	GetSSHUser(uuid string) (*SSHUser, error)
	NewSSHUserUpdateRequest(uuid string) SSHUserUpdateRequest
	DeleteSSHUser(uuid string) error

	GetL2Connection(uuid string) (*L2Connection, error)
	CreateL2Connection(conn L2Connection) (*L2Connection, error)
	CreateL2RedundantConnection(priConn L2Connection, secConn L2Connection) (*L2Connection, error)
	DeleteL2Connection(uuid string) error
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

//ChangeError describes single error that occurred during update of selected target property
type ChangeError struct {
	Type   string
	Target string
	Value  interface{}
	Cause  error
}

func (e ChangeError) Error() string {
	return fmt.Sprintf("change type '%s', target '%s', value '%s', cause: '%s'", e.Type, e.Target, e.Value, e.Cause)
}

//UpdateError describes error that occurred during composite update request and consists of multiple atomic change errors
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
	UUID                string
	Name                string
	TypeCode            string
	Status              string
	LicenseStatus       string
	MetroCode           string
	IBX                 string
	Region              string
	Throughput          int
	ThroughputUnit      string
	HostName            string
	PackageCode         string
	Version             string
	IsBYOL              bool
	LicenseToken        string
	ACLs                []string
	SSHIPAddress        string
	SSHIPFqdn           string
	AccountNumber       string
	Notifications       []string
	PurchaseOrderNumber string
	RedundancyType      string
	RedundantUUID       string
	TermLength          int
	AdditionalBandwidth int
	OrderReference      string
	InterfaceCount      int
	CoreCount           int
	IsSelfManaged       bool
	Interfaces          []DeviceInterface
	VendorConfiguration map[string]string
}

//DeviceInterface describes Network Edge device interface
type DeviceInterface struct {
	ID                int
	Name              string
	Status            string
	OperationalStatus string
	MACAddress        string
	IPAddress         string
	AssignedType      string
	Type              string
}

//DeviceACLs describes device ACLs - list of CIDRs along with their current provisioning status
type DeviceACLs struct {
	ACLs   []string
	Status string
}

//SSHUser describes Network Edge SSH user
type SSHUser struct {
	UUID        string
	Username    string
	Password    string
	DeviceUUIDs []string
}

//L2Connection describes layer 2 connection from a Network Edge device
type L2Connection struct {
	UUID                     string
	AuthorizationKey         string
	Name                     string
	NamedTag                 string
	Notifications            []string
	ProfileUUID              string
	PurchaseOrderNumber      string
	RedundantUUID            string
	SellerRegion             string
	SellerMetroCode          string
	SellerHostedConnectionID string
	Speed                    int
	SpeedUnit                string
	Status                   string
	VirtualDeviceUUID        string
	VlanSTag                 int
	ZSidePortUUID            string
	ZSideVlanSTag            int
	ZSideVlanCTag            int
}
