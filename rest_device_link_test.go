package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDeviceLinkGroups(t *testing.T) {
	//Given
	var respBody api.DeviceLinkGroupsGetResponse
	if err := readJSONData("./test-fixtures/ne_device_links_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	limit := respBody.Pagination.Limit
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/links?limit=%d", baseURL, limit),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	linkGroups, err := c.GetDeviceLinkGroups()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, linkGroups, "Client should return a response")
	assert.Equal(t, len(respBody.Data), len(linkGroups), "Number of objects matches")
	for i := range respBody.Data {
		verifyDeviceLinkGroup(t, linkGroups[i], respBody.Data[i])
	}
}

func TestGetDeviceLinkGroup(t *testing.T) {
	//Given
	var respBody api.DeviceLinkGroup
	if err := readJSONData("./test-fixtures/ne_device_link_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	uuid := "testLinkGroup"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/links/%s", baseURL, uuid),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	linkGroup, err := c.GetDeviceLinkGroup(uuid)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, linkGroup, "Client should return a response")
	verifyDeviceLinkGroup(t, *linkGroup, respBody)
}

func TestCreateDeviceLinkGroup(t *testing.T) {
	//Given
	var respBody api.DeviceLinkGroupCreateResponse
	if err := readJSONData("./test-fixtures/ne_device_link_create_resp.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	testLinkGroup := DeviceLinkGroup{
		Name:   String("testLinkGroup"),
		Subnet: String("10.1.2.0/24"),
		Devices: []DeviceLinkGroupDevice{
			{
				DeviceID:    String("c9a5c40c-b90f-4156-8460-6cb5dc98f88d"),
				ASN:         Int(12345),
				InterfaceID: Int(5),
			},
			{
				DeviceID:    String("7312336a-d508-4ac2-8d99-3070c30aed94"),
				ASN:         Int(22335),
				InterfaceID: Int(4),
			},
			{
				DeviceID:    String("01026670-b55e-46db-84ba-004573d7a9d8"),
				ASN:         Int(12451),
				InterfaceID: Int(8),
			},
		},
		Links: []DeviceLinkGroupLink{
			{
				AccountNumber:        String("22314"),
				Throughput:           String("50"),
				ThroughputUnit:       String("Mbps"),
				SourceMetroCode:      String("LD"),
				DestinationMetroCode: String("AM"),
				SourceZoneCode:       String("Zone1"),
				DestinationZoneCode:  String("Zone1"),
			},
			{
				AccountNumber:        String("10314"),
				Throughput:           String("50"),
				ThroughputUnit:       String("Mbps"),
				SourceMetroCode:      String("LD"),
				DestinationMetroCode: String("FR"),
				SourceZoneCode:       String("Zone1"),
				DestinationZoneCode:  String("Zone1"),
			},
		},
	}
	request := api.DeviceLinkGroup{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/links", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(202, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateDeviceLinkGroup(testLinkGroup)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, uuid, "Client should return a response")
	assert.Equal(t, respBody.UUID, uuid, "UUID matches")
	verifyDeviceLinkGroup(t, testLinkGroup, request)
}

func TestUpdateDeviceLinkGroup(t *testing.T) {
	//given
	groupID := "test"
	newGroupName := "newDLGroup"
	newSubnet := "20.1.1.0/24"
	newDevices := []DeviceLinkGroupDevice{
		{
			DeviceID:    String("c9a5c40c-b90f-4156-8460-6cb5dc98f88d"),
			ASN:         Int(12345),
			InterfaceID: Int(5),
		},
		{
			DeviceID:    String("7312336a-d508-4ac2-8d99-3070c30aed94"),
			ASN:         Int(22335),
			InterfaceID: Int(4),
		},
	}
	newLinks := []DeviceLinkGroupLink{
		{
			AccountNumber:        String("22314"),
			Throughput:           String("50"),
			ThroughputUnit:       String("Mbps"),
			SourceMetroCode:      String("LD"),
			DestinationMetroCode: String("AM"),
			SourceZoneCode:       String("Zone1"),
			DestinationZoneCode:  String("Zone1"),
		},
	}
	req := api.DeviceLinkGroup{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PATCH", fmt.Sprintf("%s/ne/v1/links/%s", baseURL, groupID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.NewDeviceLinkGroupUpdateRequest(groupID).
		WithGroupName(newGroupName).
		WithSubnet(newSubnet).
		WithDevices(newDevices).
		WithLinks(newLinks).
		Execute()

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, &newSubnet, req.Subnet, "Subnet matches")
	assert.Equal(t, &newGroupName, req.GroupName, "GroupName matches")
	assert.Equal(t, len(newDevices), len(req.Devices), "Devices number matches")
	for i := range newDevices {
		verifyDeviceLinkGroupDevice(t, newDevices[i], req.Devices[i])
	}
	assert.Equal(t, len(newLinks), len(req.Links), "Links number matches")
	for i := range newLinks {
		verifyDeviceLinkGroupLink(t, newLinks[i], req.Links[i])
	}
}

func TestDeleteDeviceLinkGroup(t *testing.T) {
	//given
	uuid := "testLinkGroup"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/links/%s", baseURL, uuid),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteDeviceLinkGroup(uuid)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifyDeviceLinkGroup(t *testing.T, linkGroup DeviceLinkGroup, apiLinkGroup api.DeviceLinkGroup) {
	assert.Equal(t, apiLinkGroup.UUID, linkGroup.UUID, "UUID matches")
	assert.Equal(t, apiLinkGroup.GroupName, linkGroup.Name, "GroupName matches")
	assert.Equal(t, apiLinkGroup.Subnet, linkGroup.Subnet, "Subnet matches")
	assert.Equal(t, len(apiLinkGroup.Devices), len(linkGroup.Devices), "Length of []Devices matches")
	for i := range apiLinkGroup.Devices {
		verifyDeviceLinkGroupDevice(t, linkGroup.Devices[i], apiLinkGroup.Devices[i])
	}
	assert.Equal(t, len(apiLinkGroup.Links), len(linkGroup.Links), "Length of []Links matches")
	for i := range apiLinkGroup.Links {
		verifyDeviceLinkGroupLink(t, linkGroup.Links[i], apiLinkGroup.Links[i])
	}
}

func verifyDeviceLinkGroupDevice(t *testing.T, linkGroupDevice DeviceLinkGroupDevice, apiLinkGroupDevice api.DeviceLinkGroupDevice) {
	assert.Equal(t, linkGroupDevice.DeviceID, apiLinkGroupDevice.DeviceUUID, "DeviceUUID matches")
	assert.Equal(t, linkGroupDevice.ASN, apiLinkGroupDevice.ASN, "ASN matches")
	assert.Equal(t, linkGroupDevice.InterfaceID, apiLinkGroupDevice.InterfaceID, "InterfaceID matches")
	assert.Equal(t, linkGroupDevice.Status, apiLinkGroupDevice.Status, "Status matches")
	assert.Equal(t, linkGroupDevice.IPAddress, apiLinkGroupDevice.IPAddress, "IPAddress matches")
}

func verifyDeviceLinkGroupLink(t *testing.T, linkGroupLink DeviceLinkGroupLink, apiLinkGroupLink api.DeviceLinkGroupLink) {
	assert.Equal(t, linkGroupLink.AccountNumber, apiLinkGroupLink.AccountNumber, "AccountNumber matches")
	assert.Equal(t, linkGroupLink.Throughput, apiLinkGroupLink.Throughput, "Throughput matches")
	assert.Equal(t, linkGroupLink.ThroughputUnit, apiLinkGroupLink.ThroughputUnit, "ThroughputUnit matches")
	assert.Equal(t, linkGroupLink.SourceMetroCode, apiLinkGroupLink.SourceMetroCode, "SourceMetroCode matches")
	assert.Equal(t, linkGroupLink.DestinationMetroCode, apiLinkGroupLink.DestinationMetroCode, "DestinationMetroCode matches")
	assert.Equal(t, linkGroupLink.SourceZoneCode, apiLinkGroupLink.SourceZoneCode, "SourceZoneCode matches")
	assert.Equal(t, linkGroupLink.DestinationZoneCode, apiLinkGroupLink.DestinationZoneCode, "DestinationZoneCode matches")
}
