package api

//File describes Network Edge uploaded file
type File struct {
	UUID           *string `json:"uuid,omitempty"`
	FileName       *string `json:"fileName,omitempty"`
	MetroCode      *string `json:"metroCode,omitempty"`
	DeviceTypeCode *string `json:"deviceTypeCode,omitempty"`
	ProcessType    *string `json:"processType,omitempty"`
	Status         *string `json:"status,omitempty"`
}

//FileUploadResponse describes response to file upload request
type FileUploadResponse struct {
	FileUUID *string `json:"fileUuid,omitempty"`
}
