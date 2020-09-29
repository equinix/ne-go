package api

//DeviceTypeResponse describes response for Network Edge device types query
type DeviceTypeResponse struct {
	TotalCount int          `json:"totalCount,omitempty"`
	PageNumber int          `json:"pageNumber,omitempty"`
	PageSize   int          `json:"pageSize,omitempty"`
	Content    []DeviceType `json:"content,omitempty"`
}

//DeviceType describes Network Edge device type
type DeviceType struct {
	Code                  string                      `json:"deviceTypeCode,omitempty"`
	Name                  string                      `json:"name,omitempty"`
	Description           string                      `json:"description,omitempty"`
	Vendor                string                      `json:"vendor,omitempty"`
	Category              string                      `json:"category,omitempty"`
	AvailableMetros       []DeviceTypeAvailableMetro  `json:"availableMetros,omitempty"`
	SoftwarePackages      []DeviceTypeSoftwarePackage `json:"softwarePackages,omitempty"`
	DeviceManagementTypes DeviceManagementTypes       `json:"deviceManagementTypes,omitempty"`
}

//DeviceTypeAvailableMetro describes metro in which network edge device is available
type DeviceTypeAvailableMetro struct {
	Code        string `json:"metroCode,omitempty"`
	Description string `json:"metroDescription,omitempty"`
	Region      string `json:"region,omitempty"`
}

//DeviceTypeSoftwarePackage describes device software package details
type DeviceTypeSoftwarePackage struct {
	Code           string                     `json:"packageCode,omitempty"`
	Name           string                     `json:"name,omitempty"`
	VersionDetails []DeviceTypeVersionDetails `json:"versionDetails,omitempty"`
}

//DeviceTypeVersionDetails describes device software version details
type DeviceTypeVersionDetails struct {
	Version                   string   `json:"version,omitempty"`
	ImageName                 string   `json:"imageName,omitempty"`
	Date                      string   `json:"versionDate,omitempty"`
	Status                    string   `json:"status,omitempty"`
	IsStable                  bool     `json:"stableVersion,omitempty"`
	AllowedUpgradableVersions []string `json:"allowedUpgradableVersions,omitempty"`
	IsUpgradeAllowed          bool     `json:"upgradeAllowed,omitempty"`
	ReleaseNotesLink          string   `json:"releaseNotesLink,omitempty"`
}

type DeviceManagementTypes struct {
	EquinixConfigured DeviceManagementType `json:"EQUINIX-CONFIGURED,omitempty"`
	SelfConfigured    DeviceManagementType `json:"SELF-CONFIGURED,omitempty"`
}

type DeviceManagementType struct {
	Type           string               `json:"type,omitempty"`
	LicenseOptions DeviceLicenseOptions `json:"licenseOptions,omitempty"`
	IsSupported    bool                 `json:"supported,omitempty"`
}

type DeviceLicenseOptions struct {
	Sub  DeviceLicenseOption `json:"SUB,omitempty"`
	BYOL DeviceLicenseOption `json:"BYOL,omitempty"`
}

type DeviceLicenseOption struct {
	Type        string       `json:"type,omitempty"`
	Name        string       `json:"name,omitempty"`
	Cores       []DeviceCore `json:"cores,omitempty"`
	IsSupported bool         `json:"supported,omitempty"`
}

type DeviceCore struct {
	Core         int                 `json:"core,omitempty"`
	Memory       int                 `json:"memory,omitempty"`
	Unit         string              `json:"unit,omitempty"`
	Flavor       string              `json:"flavor,omitempty"`
	PackageCodes []DevicePackageCode `json:"packageCodes,omitempty"`
	IsSupported  bool                `json:"supported,omitempty"`
}

type DevicePackageCode struct {
	PackageCode string `json:"packageCode,omitempty"`
	IsSupported bool   `json:"supported,omitempty"`
}
