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
	Code             string                      `json:"deviceTypeCode,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Description      string                      `json:"description,omitempty"`
	Vendor           string                      `json:"vendor,omitempty"`
	Category         string                      `json:"category,omitempty"`
	AvailableMetros  []DeviceTypeAvailableMetro  `json:"availableMetros,omitempty"`
	SoftwarePackages []DeviceTypeSoftwarePackage `json:"softwarePackages,omitempty"`
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
