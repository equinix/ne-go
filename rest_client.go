package ne

import (
	"context"
	"net/http"

	"github.com/equinix/rest-go"
)

//RestClient describes REST implementation of Network Edge Client
type RestClient struct {
	*rest.Client
}

//NewClient creates new REST Network Edge client with a given baseURL, context and httpClient
func NewClient(ctx context.Context, baseURL string, httpClient *http.Client) *RestClient {
	rest := rest.NewClient(ctx, baseURL, httpClient)
	rest.SetHeader("User-agent", "equinix/ne-go")
	return &RestClient{rest}
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

const (
	changeTypeCreate = "Add"
	changeTypeUpdate = "Update"
	changeTypeDelete = "Delete"
)
