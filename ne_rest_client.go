package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/equinix/ne-go/internal/api"
	"github.com/go-resty/resty/v2"
)

//RestClient describes REST implementation of Network Edge Client
type RestClient struct {
	//PageSize determines default page size for GET requests on resource collections
	PageSize int
	baseURL  string
	ctx      context.Context
	*resty.Client
}

//RestError describes Network Edge error specific to REST implementation
type RestError struct {
	HTTPCode int
	Message  string
	Errors   []Error
}

func (e RestError) Error() string {
	return fmt.Sprintf("network edge rest error: httpCode: %v, message: %v", e.HTTPCode, e.Message)
}

//NewClient creates new REST Network Edge client with a given baseURL, context and httpClient
func NewClient(ctx context.Context, baseURL string, httpClient *http.Client) *RestClient {
	resty := resty.NewWithClient(httpClient)
	resty.SetHeader("User-agent", "equinix/ne-go")
	resty.SetHeader("Accept", "application/json")
	return &RestClient{
		100,
		baseURL,
		ctx,
		resty}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

const (
	changeTypeCreate = "Add"
	changeTypeUpdate = "Update"
	changeTypeDelete = "Delete"
)

func (c RestClient) execute(req *resty.Request, method string, url string) error {
	resp, err := req.SetContext(c.ctx).Execute(method, url)
	if err != nil {
		restErr := RestError{Message: fmt.Sprintf("operation failed: %s", err)}
		if resp != nil {
			restErr.HTTPCode = resp.StatusCode()
		}
		return restErr
	}
	if resp.IsError() {
		err := transformErrorBody(resp.Body())
		err.HTTPCode = resp.StatusCode()
		return err
	}
	return nil
}

func transformErrorBody(body []byte) RestError {
	apiError := api.ErrorResponse{}
	if err := json.Unmarshal(body, &apiError); err == nil {
		return mapErrorAPIToDomain(apiError)
	}
	apiErrors := api.ErrorResponses{}
	if err := json.Unmarshal(body, &apiErrors); err == nil {
		return mapErrorsAPIToDomain(apiErrors)
	}
	return RestError{
		Message: string(body)}
}

func mapErrorAPIToDomain(apiError api.ErrorResponse) RestError {
	return RestError{
		Message: apiError.ErrorMessage,
		Errors: []Error{{
			apiError.ErrorCode,
			fmt.Sprintf("[Error: Property: %v, %v]", apiError.Property, apiError.ErrorMessage),
		}},
	}
}

func mapErrorsAPIToDomain(apiErrors api.ErrorResponses) RestError {
	errors := make([]Error, len(apiErrors))
	msg := ""
	for i, v := range apiErrors {
		errors[i] = Error{v.ErrorCode, v.ErrorMessage}
		msg = msg + fmt.Sprintf(" [Error %v: Property: %v, %v]", i+1, v.Property, v.ErrorMessage)
	}
	return RestError{
		Message: "Multiple errors occurred: " + msg,
		Errors:  errors,
	}
}
