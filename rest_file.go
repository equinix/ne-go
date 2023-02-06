package ne

import (
	"github.com/equinix/ne-go/internal/api"
	"io"
	"net/http"
	"net/url"
)

const (
	//ProcessTypeLicense indicates file type where customer is uploading a license file
	ProcessTypeLicense = "LICENSE"

	//ProcessTypeCloudInit indicates file type where customer is uploading a cloud_init file
	ProcessTypeCloudInit = "CLOUD_INIT"
)

//UploadFile performs multipart upload of a cloud_init/license file from a given reader interface
//along with provided data. Uploaded file identifier is returned on success.
func (c RestClient) UploadFile(metroCode, deviceTypeCode, processType, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (*string, error) {
	path := "/ne/v1/files"
	respBody := api.FileUploadResponse{}
	req := c.R().
		SetFileReader("file", fileName, reader).
		SetFormData(map[string]string{
			"metroCode":            metroCode,
			"deviceTypeCode":       deviceTypeCode,
			"processType":          processType,
			"licenseType":          licenseMode,
			"deviceManagementType": deviceManagementMode,
		}).
		SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.FileUUID, nil
}

//GetFile retrieves file metadata with a given UUID
func (c RestClient) GetFile(uuid string) (*File, error) {
	path := "/ne/v1/files/" + url.PathEscape(uuid)
	respBody := api.File{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	file := mapFileAPIToDomain(respBody)
	return &file, nil
}

func mapFileAPIToDomain(apiFile api.File) File {
	return File{
		UUID:           apiFile.UUID,
		FileName:       apiFile.FileName,
		MetroCode:      apiFile.MetroCode,
		DeviceTypeCode: apiFile.DeviceTypeCode,
		ProcessType:    apiFile.ProcessType,
		Status:         apiFile.Status,
	}
}
