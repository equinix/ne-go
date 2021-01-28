package ne

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8888"
)

func TestClientImplementation(t *testing.T) {
	//given
	cli := NewClient(context.Background(), baseURL, &http.Client{})
	//then
	assert.Implements(t, (*Client)(nil), cli, "Rest client implements Client interface")
}

func readJSONData(filePath string, target interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func setupMockedClient(method string, url string, respCode int, resp interface{}) *http.Client {
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder(method, url,
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(respCode, resp)
			return resp, nil
		},
	)
	return testHc
}
