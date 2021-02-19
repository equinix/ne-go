package ne

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDeviceTypes(t *testing.T) {
	//Given
	respBody := api.DeviceTypeResponse{}
	if err := readJSONData("./test-fixtures/ne_device_types_get.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	limit := respBody.Pagination.Limit
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/deviceTypes?limit=%d", baseURL, limit),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	types, err := c.GetDeviceTypes()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, types, "Client should return a response")
	assert.Equal(t, len(respBody.Data), len(types), "Number of objects matches")
	for i := range respBody.Data {
		verifyDeviceType(t, respBody.Data[i], types[i])
	}
}

func TestGetDeviceSoftwareVersions(t *testing.T) {
	//Given
	respBody := api.DeviceTypeResponse{}
	if err := readJSONData("./test-fixtures/ne_devices_types_csr1000v_get.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	limit := respBody.Pagination.Limit
	deviceTypeCode := "CSR1000V"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/deviceTypes?deviceTypeCode=%s&limit=%d", baseURL, deviceTypeCode, limit),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	versions, err := c.GetDeviceSoftwareVersions(deviceTypeCode)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, versions, "Client should return a response")
	apiType := respBody.Data[0]
	apiVerMap := make(map[string]api.DeviceTypeVersionDetails)
	for _, pkg := range apiType.SoftwarePackages {
		for _, ver := range pkg.VersionDetails {
			if _, ok := apiVerMap[StringValue(ver.Version)]; !ok {
				apiVerMap[StringValue(ver.Version)] = ver
			}
		}
	}
	assert.Equal(t, len(apiVerMap), len(versions), "Number of versions matches")
	for _, ver := range versions {
		apiVer := apiVerMap[StringValue(ver.Version)]
		verifyDeviceSoftwareVersion(t, apiVer, ver)
	}
}

func TestGetDevicePlatforms(t *testing.T) {
	//Given
	respBody := api.DeviceTypeResponse{}
	if err := readJSONData("./test-fixtures/ne_devices_types_csr1000v_get.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	limit := respBody.Pagination.Limit
	deviceTypeCode := "CSR1000V"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/deviceTypes?deviceTypeCode=%s&limit=%d", baseURL, deviceTypeCode, limit),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	platforms, err := c.GetDevicePlatforms(deviceTypeCode)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, platforms, "Client should return a response")
	assert.Equal(t, 3, len(platforms), "Number of platforms matches")
	for _, version := range platforms {
		assert.ElementsMatch(t, version.PackageCodes, []string{"APPX", "AX", "IPBASE", "SEC"}, "PackageCodes match")
		assert.ElementsMatch(t, version.ManagementTypes, []string{"EQUINIX-CONFIGURED", "SELF-CONFIGURED"}, "ManagementTypes match")
		assert.ElementsMatch(t, version.LicenseOptions, []string{"BYOL", "Sub"}, "LicenseOptions match")
	}
}

func verifyDeviceType(t *testing.T, apiDeviceType api.DeviceType, deviceType DeviceType) {
	assert.Equal(t, apiDeviceType.Name, deviceType.Name, "Name matches")
	assert.Equal(t, apiDeviceType.Description, deviceType.Description, "Description matches")
	assert.Equal(t, apiDeviceType.Code, deviceType.Code, "Code matches")
	assert.Equal(t, apiDeviceType.Vendor, deviceType.Vendor, "Vendor matches")
	assert.Equal(t, apiDeviceType.Category, deviceType.Category, "Category matches")
	assert.Equal(t, len(apiDeviceType.AvailableMetros), len(deviceType.MetroCodes), "Number of available metros matches")
	for i := range apiDeviceType.AvailableMetros {
		assert.Equalf(t, *apiDeviceType.AvailableMetros[i].Code, deviceType.MetroCodes[i], "Code of available metro element %d matches", i)
	}
}

func verifyDeviceSoftwareVersion(t *testing.T, apiVer api.DeviceTypeVersionDetails, ver DeviceSoftwareVersion) {
	assert.Equal(t, apiVer.Version, ver.Version, "Version matches")
	assert.Equal(t, apiVer.ImageName, ver.ImageName, "ImageName matches")
	assert.Equal(t, apiVer.Date, ver.Date, "Date matches")
	assert.Equal(t, apiVer.Status, ver.Status, "Status matches")
	assert.Equal(t, apiVer.IsStable, ver.IsStable, "IsStable matches")
	assert.Equal(t, apiVer.ReleaseNotesLink, ver.ReleaseNotesLink, "ReleaseNotesLink matches")
}
