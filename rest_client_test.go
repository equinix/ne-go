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

func TestParseResourceIdFromLocationHeader(t *testing.T) {
	//given
	resourceID := "3c11e8d9-80da-4a04-ae22-a35313d64717"
	header := "/v1/publicKeys/" + resourceID
	//when
	id, err := parseResourceIDFromLocationHeader(header)
	//then
	assert.Nil(t, err, "Error is not returned")
	assert.NotNil(t, id, "Returned resource ID is not nil")
	assert.Equal(t, resourceID, *id, "Resource IDs match")
}

func TestParseResourceIdFromLocationHeader_negative(t *testing.T) {
	//given
	headerOne := "dummy"
	headerTwo := "/dummy"
	//when
	_, errOne := parseResourceIDFromLocationHeader(headerOne)
	_, errTwo := parseResourceIDFromLocationHeader(headerTwo)
	//then
	assert.NotNil(t, errOne, "Error is returned")
	assert.NotNil(t, errTwo, "Error is returned")
}

type mockedHeaderProvider struct {
	headers map[string][]string
}

func (m *mockedHeaderProvider) Header() http.Header {
	return m.headers
}

func TestGetLocationHeaderValue(t *testing.T) {
	//given
	locationValue := "/path/to/resource"
	mock := &mockedHeaderProvider{
		map[string][]string{
			"Location": {locationValue},
		},
	}
	//when
	value, err := getLocationHeaderValue(mock)
	//then
	assert.Nil(t, err, "Error is nil")
	assert.NotNil(t, value, "Value is not nil")
	assert.Equal(t, locationValue, *value, "Values match")
}

func TestGetLocationHeaderValue_negative(t *testing.T) {
	//given
	//when
	_, errOne := getLocationHeaderValue(&mockedHeaderProvider{})
	_, errTwo := getLocationHeaderValue(&mockedHeaderProvider{
		map[string][]string{
			"Location": {"valueOne", "valueTwo"},
		}})
	//then
	assert.NotNil(t, errOne, "Error is returned")
	assert.NotNil(t, errTwo, "Error is returned")
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
