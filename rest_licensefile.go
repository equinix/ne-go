package ne

import (
	"io"

	"github.com/equinix/ne-go/internal/api"
	"github.com/go-resty/resty/v2"
)

//UploadLicenseFile performs multipart upload of a license file from a given reader interface
//along with provided data. Uploaded file identifier is returned on success.
func (c RestClient) UploadLicenseFile(metroCode, deviceTypeCode, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (string, error) {
	path := "/ne/v1/device/license/file"
	respBody := api.LicenseFileUploadResponse{}
	req := c.R().
		SetFileReader("file", fileName, reader).
		SetFormData(map[string]string{
			"metroCode":            metroCode,
			"deviceTypeCode":       deviceTypeCode,
			"licenseType":          licenseMode,
			"deviceManagementType": deviceManagementMode,
		}).
		SetResult(&respBody)
	if err := c.Execute(req, resty.MethodPost, path); err != nil {
		return "", err
	}
	return respBody.FileID, nil
}
