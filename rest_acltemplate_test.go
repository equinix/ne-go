package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/equinix/ne-go/internal/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testACLTemplate = ACLTemplate{
	Name:        String("test"),
	Description: String("Test ACL"),
	MetroCode:   String("SV"),
	InboundRules: []ACLTemplateInboundRule{
		{
			SrcType:     String("SUBNET"),
			SeqNo:       Int(1),
			Subnets:     []string{"10.0.0.0/24"},
			Protocol:    String("TCP"),
			SrcPort:     String("any"),
			DstPort:     String("22"),
			Description: String("Description of the rule"),
		},
		{
			SrcType:  String("DOMAIN"),
			SeqNo:    Int(2),
			Subnets:  []string{"216.221.225.13/32"},
			Protocol: String("TCP"),
			SrcPort:  String("any"),
			DstPort:  String("1024-10000"),
		},
	},
}

func TestCreateACLTemplate(t *testing.T) {
	//given
	newUUID := "299cd6f2-714e-4265-a07c-48944a6ac3bd"
	template := testACLTemplate
	reqBody := api.ACLTemplate{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/aclTemplates", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp := httpmock.NewStringResponse(201, "")
			resp.Header.Add("Location", "/ne/v1/aclTemplates/"+newUUID)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	uuid, err := c.CreateACLTemplate(template)

	//then
	assert.Nil(t, err, "Error is not returned")
	assert.Equal(t, newUUID, *uuid, "UUID matches")
	verifyACLTemplate(t, template, reqBody)
}

func TestGetACLTemplates(t *testing.T) {
	//Given
	var respBody api.ACLTemplatesResponse
	if err := readJSONData("./test-fixtures/ne_acltemplates_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "cannot read test response due to %s", err.Error())
	}
	limit := respBody.Pagination.Limit
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/aclTemplates?limit=%d", baseURL, limit),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//When
	c := NewClient(context.Background(), baseURL, testHc)
	c.PageSize = limit
	templates, err := c.GetACLTemplates()

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, templates, "Client should return a response")
	assert.Equal(t, len(respBody.Data), len(templates), "Number of objects matches")
	for i := range respBody.Data {
		verifyACLTemplate(t, templates[i], respBody.Data[i])
		verifyACLTemplateDeviceDetails(t, templates[i], respBody.Data[i])
	}
}

func TestGetACLTemplate(t *testing.T) {
	//given
	resp := api.ACLTemplate{}
	if err := readJSONData("./test-fixtures/ne_acltemplate_get_resp.json", &resp); err != nil {
		assert.Fail(t, "Cannot read test response")
	}
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	testHc := setupMockedClient("GET", fmt.Sprintf("%s/ne/v1/aclTemplates/%s", baseURL, templateID), 200, resp)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	template, err := c.GetACLTemplate(templateID)

	//then
	assert.NotNil(t, template, "Returned template is not nil")
	assert.Nil(t, err, "Error is not returned")
	verifyACLTemplate(t, *template, resp)
	verifyACLTemplateDeviceDetails(t, *template, resp)
}

func TestReplaceACLTemplate(t *testing.T) {
	//given
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	template := testACLTemplate
	reqBody := api.ACLTemplate{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/ne/v1/aclTemplates/%s", baseURL, templateID),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			return httpmock.NewStringResponse(204, ""), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.ReplaceACLTemplate(templateID, template)

	//then
	assert.Nil(t, err, "Error is not returned")
	verifyACLTemplate(t, template, reqBody)
}

func TestDeleteACLTemplate(t *testing.T) {
	//given
	templateID := "db66bf49-b2d8-4e64-8719-d46406b54039"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/aclTemplates/%s", baseURL, templateID),
		httpmock.NewStringResponder(204, ""))
	defer httpmock.DeactivateAndReset()

	//when
	c := NewClient(context.Background(), baseURL, testHc)
	err := c.DeleteACLTemplate(templateID)

	//then
	assert.Nil(t, err, "Error is not returned")
}

func verifyACLTemplate(t *testing.T, template ACLTemplate, apiTemplate api.ACLTemplate) {
	assert.Equal(t, template.UUID, apiTemplate.UUID, "UUID matches")
	assert.Equal(t, template.Name, apiTemplate.Name, "Name matches")
	assert.Equal(t, template.Description, apiTemplate.Description, "Description matches")
	assert.Equal(t, template.MetroCode, apiTemplate.MetroCode, "MetroCode matches")
	assert.Equal(t, template.DeviceACLStatus, apiTemplate.DeviceACLStatus, "DeviceACLStatus matches")
	assert.Equal(t, len(template.InboundRules), len(apiTemplate.InboundRules), "Number of InboundRules matches")
	for i := range template.InboundRules {
		verifyACLTemplateInboundRule(t, template.InboundRules[i], apiTemplate.InboundRules[i])
	}
}

func verifyACLTemplateInboundRule(t *testing.T, rule ACLTemplateInboundRule, apiRule api.ACLTemplateInboundRule) {
	assert.Equal(t, rule.SeqNo, apiRule.SeqNO, "SeqNo matches")
	assert.Equal(t, rule.SrcType, apiRule.SrcType, "SrcType matches")
	assert.ElementsMatch(t, rule.Subnets, apiRule.Subnets, "Subnets matches")
	assert.Equal(t, rule.Protocol, apiRule.Protocol, "Protocol matches")
	assert.Equal(t, rule.SrcPort, apiRule.SrcPort, "SrcPort matches")
	assert.Equal(t, rule.DstPort, apiRule.DstPort, "DstPort matches")
	assert.Equal(t, rule.Description, apiRule.Description, "Description matches")
}

func verifyACLTemplateDeviceDetails(t *testing.T, template ACLTemplate, apiTemplate api.ACLTemplate) {
	assert.Equal(t, len(template.DeviceDetails), len(apiTemplate.DeviceDetails), "Number of DeviceDetails matches")
	for i := range template.DeviceDetails {
		assert.Equal(t, template.DeviceDetails[i].UUID, apiTemplate.DeviceDetails[i].UUID, "UUID matches")
		assert.Equal(t, template.DeviceDetails[i].Name, apiTemplate.DeviceDetails[i].Name, "Name matches")
		assert.Equal(t, template.DeviceDetails[i].ACLStatus, apiTemplate.DeviceDetails[i].ACLStatus, "ACL Status matches")
	}
}
