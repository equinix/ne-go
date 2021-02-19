package ne

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUploadLicenseFile(t *testing.T) {
	//given
	licenseFile, err := os.Open("test-fixtures/test_license_file.lic")
	if err != nil {
		assert.Fail(t, "Cannot read test license file")
	}
	defer licenseFile.Close()
	resp := api.LicenseFileUploadResponse{}
	if err := readJSONData("./test-fixtures/ne_licensefile_upload_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	testHc := &http.Client{}
	metroCode := "SV"
	deviceTypeCode := "CSRSDWAN"
	licenseMode := DeviceLicenseModeBYOL
	managementMode := DeviceManagementTypeSelf
	fileName := "CSRSDWAN.cfg"
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/devices/licenseFiles", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := r.ParseMultipartForm(32 << 20); err != nil {
				return httpmock.NewStringResponse(400, err.Error()), nil
			}
			assert.Equal(t, metroCode, r.MultipartForm.Value["metroCode"][0], "Form metroCode matches")
			assert.Equal(t, deviceTypeCode, r.MultipartForm.Value["deviceTypeCode"][0], "Form deviceTypeCode matches")
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
	id, err := c.UploadLicenseFile(metroCode, deviceTypeCode, managementMode, licenseMode, fileName, licenseFile)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, resp.FileID, id, "File identifier matches")
}
