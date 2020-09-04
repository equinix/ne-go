package api

//Device describes network edge device
type Device struct {
	UUID                 string                 `json:"uuid,omitempty"`
	Name                 string                 `json:"name,omitempty"`
	DeviceTypeCode       string                 `json:"deviceTypeCode,omitempty"`
	Status               string                 `json:"status,omitempty"`
	LicenseStatus        string                 `json:"licenseStatus,omitempty"`
	MetroCode            string                 `json:"metroCode,omitempty"`
	IBX                  string                 `json:"ibx,omitempty"`
	Region               string                 `json:"region,omitempty"`
	Throughput           int                    `json:"throughput,omitempty,string"`
	ThroughputUnit       string                 `json:"throughputUnit,omitempty"`
	HostName             string                 `json:"hostName,omitempty"`
	PackageCode          string                 `json:"packageCode,omitempty"`
	Version              string                 `json:"version,omitempty"`
	LicenseType          string                 `json:"licenseType,omitempty"`
	ACL                  []string               `json:"acl,omitempty"`
	SSHIPAddress         string                 `json:"sshIpAddress,omitempty"`
	SSHIPFqdn            string                 `json:"sshIpFqdn,omitempty"`
	AccountNumber        string                 `json:"accountNumber,omitempty"`
	Notifications        []string               `json:"notifications,omitempty"`
	PurchaseOrderNumber  string                 `json:"purchaseOrderNumber,omitempty"`
	RedundancyType       string                 `json:"redundancyType,omitempty"`
	RedundantUUID        string                 `json:"redundantUUID,omitempty"`
	TermLength           int                    `json:"termLength,omitempty"`
	AdditionalBandwidth  int                    `json:"additionalBandwidth,omitempty"`
	OrderReference       string                 `json:"orderReference,omitempty"`
	InterfaceCount       int                    `json:"interfaceCount,omitempty"`
	Core                 *DeviceCoreInformation `json:"core,omitempty"`
	DeviceManagementType string                 `json:"deviceManagementType,omitempty"`
	Interfaces           []DeviceInterface      `json:"interfaces,omitempty"`
	VendorConfig         map[string]string      `json:"vendorConfig,omitempty"`
}

//DeviceRequest describes network edge device creation request
type DeviceRequest struct {
	Throughput           int                     `json:"throughput,omitempty,string"`
	ThroughputUnit       string                  `json:"throughputUnit,omitempty"`
	MetroCode            string                  `json:"metroCode,omitempty"`
	DeviceTypeCode       string                  `json:"deviceTypeCode,omitempty"`
	TermLength           string                  `json:"termLength,omitempty"`
	LicenseMode          string                  `json:"licenseMode,omitempty"`
	LicenseToken         string                  `json:"licenseToken,omitempty"`
	PackageCode          string                  `json:"packageCode,omitempty"`
	VirtualDeviceName    string                  `json:"virtualDeviceName,omitempty"`
	Notifications        []string                `json:"notifications,omitempty"`
	HostNamePrefix       string                  `json:"hostNamePrefix,omitempty"`
	OrderReference       string                  `json:"orderReference,omitempty"`
	PurchaseOrderNumber  string                  `json:"purchaseOrderNumber,omitempty"`
	AccountNumber        string                  `json:"accountNumber,omitempty"`
	Version              string                  `json:"version,omitempty"`
	InterfaceCount       int                     `json:"interfaceCount,omitempty"`
	DeviceManagementType string                  `json:"deviceManagementType,omitempty"`
	Core                 int                     `json:"core,omitempty"`
	AdditionalBandwidth  int                     `json:"additionalBandwidth,omitempty,string"`
	FqdnACL              []DeviceFqdnACL         `json:"fqdnAcl,omitempty"`
	VendorConfig         map[string]string       `json:"vendorConfig,omitempty"`
	Secondary            *SecondaryDeviceRequest `json:"secondary,omitempty"`
}

//SecondaryDeviceRequest describes secondary device part of device creation request
type SecondaryDeviceRequest struct {
	MetroCode           string            `json:"metroCode,omitempty"`
	VirtualDeviceName   string            `json:"virtualDeviceName,omitempty"`
	Notifications       []string          `json:"notifications,omitempty"`
	HostNamePrefix      string            `json:"hostNamePrefix,omitempty"`
	AccountNumber       string            `json:"accountNumber,omitempty"`
	AdditionalBandwidth int               `json:"additionalBandwidth,omitempty,string"`
	FqdnACL             []DeviceFqdnACL   `json:"fqdnAcl,omitempty"`
	VendorConfig        map[string]string `json:"vendorConfig,omitempty"`
}

//DeviceFqdnACL describes device ACL FQDN format
type DeviceFqdnACL struct {
	CIDRs []string `json:"cidrs,omitempty"`
	Type  string   `json:"type,omitempty"`
}

//DeviceFqdnACLResponse describes respone for device FQDN ACLs fetch request
type DeviceFqdnACLResponse struct {
	FqdnACLs []DeviceFqdnACL `json:"fqdnAcl,omitempty"`
	Status   string          `json:"status,omitempty"`
}

//DeviceInterface describes device network interface
type DeviceInterface struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Status            string `json:"status,omitempty"`
	OperationalStatus string `json:"operationalStatus,omitempty"`
	MACAddress        string `json:"macAddress,omitempty"`
	IPAddress         string `json:"ipAddress,omitempty"`
	AssignedType      string `json:"assignedType,omitempty"`
	Type              string `json:"type,omitempty"`
}

//DeviceCoreInformation describes device core and memory information
type DeviceCoreInformation struct {
	Core   int    `json:"core,omitempty"`
	Memory int    `json:"memory,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

//DeviceRequestResponse describes response for device creation request
type DeviceRequestResponse struct {
	UUID          string `json:"uuid,omitempty"`
	SecondaryUUID string `json:"secondaryUUID,omitempty"`
}

//DeviceUpdateRequest describes network device update request
type DeviceUpdateRequest struct {
	Notifications     []string `json:"notifications"`
	TermLength        int      `json:"termLength,omitempty"`
	VirtualDeviceName string   `json:"virtualDeviceName,omitempty"`
}

//DeviceAdditionalBandwidthUpdateRequest describes network device additional bandwidth update request
type DeviceAdditionalBandwidthUpdateRequest struct {
	AdditionalBandwidth int `json:"additionalBandwidth"`
}

//DevicesResponse describes response for a get device list request
type DevicesResponse struct {
	TotalCount int      `json:"totalCount,omitempty"`
	PageSize   int      `json:"pageSize,omitempty"`
	Content    []Device `json:"content,omitempty"`
}
