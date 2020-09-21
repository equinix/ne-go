package ne

import (
	"fmt"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
)

//GetDeviceTypes retrieves list of devices types along with their details
func (c RestClient) GetDeviceTypes() ([]DeviceType, error) {
	url := fmt.Sprintf("%s/ne/v1/device/type", c.baseURL)
	content, err := c.GetPaginated(url, &api.DeviceTypeResponse{},
		DefaultPagingConfig())
	if err != nil {
		return nil, err
	}
	transformed := make([]DeviceType, len(content))
	for i := range content {
		transformed[i] = mapDeviceTypeAPIToDomain(content[i].(api.DeviceType))
	}
	return transformed, nil
}

//GetDeviceSoftwareVersions retrieves list of available software versions for a given device type
func (c RestClient) GetDeviceSoftwareVersions(deviceTypeCode string) ([]DeviceSoftwareVersion, error) {
	reqURL := fmt.Sprintf("%s/ne/v1/device/type", c.baseURL)
	content, err := c.GetPaginated(reqURL, &api.DeviceTypeResponse{},
		DefaultPagingConfig().
			SetAdditionalParams(map[string]string{"deviceTypeCode": url.QueryEscape(deviceTypeCode)}))
	if err != nil {
		return nil, err
	}
	if len(content) < 1 {
		return nil, fmt.Errorf("device type query returned no results for given type code: %s", deviceTypeCode)
	}
	if len(content) > 1 {
		return nil, fmt.Errorf("device type query returned more than one result for a given type code: %s", deviceTypeCode)
	}
	return mapDeviceTypeAPIToDeviceSoftwareVersions(content[0].(api.DeviceType)), nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func mapDeviceTypeAPIToDomain(apiDevice api.DeviceType) DeviceType {
	return DeviceType{
		Name:        apiDevice.Name,
		Code:        apiDevice.Code,
		Description: apiDevice.Description,
		Vendor:      apiDevice.Vendor,
		Category:    apiDevice.Category,
		MetroCodes:  mapDeviceTypeAvailableMetrosAPIToDomain(apiDevice.AvailableMetros),
	}
}

func mapDeviceTypeAvailableMetrosAPIToDomain(apiMetros []api.DeviceTypeAvailableMetro) []string {
	transformed := make([]string, len(apiMetros))
	for i := range apiMetros {
		transformed[i] = apiMetros[i].Code
	}
	return transformed
}

func mapDeviceTypeAPIToDeviceSoftwareVersions(apiType api.DeviceType) []DeviceSoftwareVersion {
	versionMap := make(map[string]*DeviceSoftwareVersion)
	for _, apiPkg := range apiType.SoftwarePackages {
		for _, apiVer := range apiPkg.VersionDetails {
			ver, ok := versionMap[apiVer.Version]
			if !ok {
				ver = mapDeviceSoftwareVersionAPIToDomain(apiVer)
				versionMap[apiVer.Version] = ver
			}
			ver.PackageCodes = append(ver.PackageCodes, apiPkg.Code)
		}
	}
	transformed := make([]DeviceSoftwareVersion, 0, len(versionMap))
	for ver := range versionMap {
		transformed = append(transformed, *versionMap[ver])
	}
	return transformed
}

func mapDeviceSoftwareVersionAPIToDomain(apiVer api.DeviceTypeVersionDetails) *DeviceSoftwareVersion {
	return &DeviceSoftwareVersion{
		Version:          apiVer.Version,
		ImageName:        apiVer.ImageName,
		Date:             apiVer.Date,
		Status:           apiVer.Status,
		IsStable:         apiVer.IsStable,
		ReleaseNotesLink: apiVer.ReleaseNotesLink,
		PackageCodes:     []string{},
	}
}
