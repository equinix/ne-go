package ne

import (
	"fmt"

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
