package api

//ErrorResponses multiple errors response
type ErrorResponses []ErrorResponse

//ErrorResponse error response
type ErrorResponse struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	MoreInfo     string `json:"moreInfo,omitempty"`
	Property     string `json:"property,omitempty"`
}
