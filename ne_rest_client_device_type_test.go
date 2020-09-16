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
	pageSize := respBody.PageSize
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/device/type?size=%d", baseURL, pageSize),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = pageSize
	types, err := c.GetDeviceTypes()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, types, "Client should return a response")
	assert.Equal(t, len(respBody.Content), len(types), "Number of objects matches")
	for i := range respBody.Content {
		verifyDeviceType(t, respBody.Content[i], types[i])
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
		assert.Equalf(t, apiDeviceType.AvailableMetros[i].Code, deviceType.MetroCodes[i], "Code of available metro element %d matches", i)
	}
}
