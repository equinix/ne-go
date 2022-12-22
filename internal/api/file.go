package api

//FileUploadResponse describes response to file upload request
type FileUploadResponse struct {
	FileUUID *string `json:"fileUuid,omitempty"`
}
