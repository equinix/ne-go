package ne

import (
	"context"
	"fmt"
	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	//given
	cloudInitFile, err := os.Open("test-fixtures/test_cloud_init_file.txt")
	if err != nil {
		assert.Fail(t, "Cannot read test cloud_init file")
	}
	defer cloudInitFile.Close()
	resp := api.FileUploadResponse{}
	if err := readJSONData("./test-fixtures/ne_file_upload_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	testHc := &http.Client{}
	metroCode := "SV"
	deviceTypeCode := "AVIATRIX_EDGE"
	processType := ProcessTypeCloudInit
	managementMode := DeviceManagementTypeSelf
	licenseMode := DeviceLicenseModeBYOL
	fileName := "AVIATRIX.txt"
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/files", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				return httpmock.NewStringResponse(400, err.Error()), nil
			}
			assert.Equal(t, metroCode, r.MultipartForm.Value["metroCode"][0], "Form metroCode matches")
			assert.Equal(t, deviceTypeCode, r.MultipartForm.Value["deviceTypeCode"][0], "Form deviceTypeCode matches")
			assert.Equal(t, processType, r.MultipartForm.Value["processType"][0], "Form processType matches")
			assert.Equal(t, managementMode, r.MultipartForm.Value["deviceManagementType"][0], "Form deviceManagementType matches")
			assert.Equal(t, licenseMode, r.MultipartForm.Value["licenseType"][0], "Form licenseType matches")
			assert.NotNil(t, r.MultipartForm.File["file"])
			resp, _ := httpmock.NewJsonResponse(201, resp)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	id, err := c.UploadFile(metroCode, deviceTypeCode, processType, managementMode, licenseMode, fileName, cloudInitFile)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, resp.FileUUID, id, "File identifier matches")
}

func TestGetFile(t *testing.T) {
	//given
	resp := api.File{}
	if err := readJSONData("./test-fixtures/ne_file_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	fileID := "26728391-2706-4135-87f2-19822bcb4721"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/files/%s", baseURL, fileID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	file, err := c.GetFile(fileID)

	//then
	assert.NotNil(t, file, "Returned file is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyFile(t, resp, *file)
}

func verifyFile(t *testing.T, apiFile api.File, file File) {
	assert.Equal(t, apiFile.UUID, file.UUID, "UUID matches")
	assert.Equal(t, apiFile.FileName, file.FileName, "FileName matches")
	assert.Equal(t, apiFile.MetroCode, file.MetroCode, "MetroCode matches")
	assert.Equal(t, apiFile.DeviceTypeCode, file.DeviceTypeCode, "DeviceTypeCode matches")
	assert.Equal(t, apiFile.ProcessType, file.ProcessType, "ProcessType matches")
	assert.Equal(t, apiFile.Status, file.Status, "Status matches")
}
