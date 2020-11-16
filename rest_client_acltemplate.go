package ne

import (
	"net/url"

	"github.com/equinix/ne-go/internal/api"
	"github.com/equinix/rest-go"
	"github.com/go-resty/resty/v2"
)

//CreateACLTemplate creates new ACL template with a given model
//On successful creation, template's UUID is returned
func (c RestClient) CreateACLTemplate(template ACLTemplate) (string, error) {
	path := "/ne/v1/device/acl-template"
	reqBody := mapACLTemplateDomainToAPI(template)
	respBody := api.ACLTemplateCreateResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.Execute(req, resty.MethodPost, path); err != nil {
		return "", err
	}
	return respBody.UUID, nil
}

//GetACLTemplates retrieves list of all ACL templates along with their details
func (c RestClient) GetACLTemplates() ([]ACLTemplate, error) {
	path := "/ne/v1/device/acl-template"
	content, err := c.GetPaginated(path, &api.ACLTemplatesResponse{},
		rest.DefaultPagingConfig().SetPageParamName("pageNumber"))
	if err != nil {
		return nil, err
	}
	transformed := make([]ACLTemplate, len(content))
	for i := range content {
		transformed[i] = mapACLTemplateAPIToDomain(content[i].(api.ACLTemplate))
	}
	return transformed, nil
}

//GetACLTemplate retrieves ACL template with a given UUID
func (c RestClient) GetACLTemplate(uuid string) (*ACLTemplate, error) {
	path := "/ne/v1/device/acl-template/" + url.PathEscape(uuid)
	respBody := api.ACLTemplate{}
	req := c.R().SetResult(&respBody)
	if err := c.Execute(req, resty.MethodGet, path); err != nil {
		return nil, err
	}
	template := mapACLTemplateAPIToDomain(respBody)
	return &template, nil
}

//ReplaceACLTemplate replaces ACL template under given UUID with
//a new one with a given model
func (c RestClient) ReplaceACLTemplate(uuid string, template ACLTemplate) error {
	path := "/ne/v1/device/acl-template/" + url.PathEscape(uuid)
	updateTemplate := ACLTemplate{
		Name:         template.Name,
		Description:  template.Description,
		MetroCode:    template.MetroCode,
		InboundRules: template.InboundRules,
	}
	reqBody := mapACLTemplateDomainToAPI(updateTemplate)
	req := c.R().SetBody(&reqBody)
	if err := c.Execute(req, resty.MethodPut, path); err != nil {
		return err
	}
	return nil
}

//DeleteACLTemplate removes ACL template with a given UUID
func (c RestClient) DeleteACLTemplate(uuid string) error {
	path := "/ne/v1/device/acl-template/" + url.PathEscape(uuid)
	if err := c.Execute(c.R(), resty.MethodDelete, path); err != nil {
		return err
	}
	return nil
}

func mapACLTemplateDomainToAPI(template ACLTemplate) api.ACLTemplate {
	return api.ACLTemplate{
		UUID:              template.UUID,
		Name:              template.Name,
		Description:       template.Description,
		MetroCode:         template.MetroCode,
		VirtualDeviceUUID: template.DeviceUUID,
		DeviceACLStatus:   template.DeviceACLStatus,
		InboundRules:      mapACLTemplateInboundRulesDomainToAPI(template.InboundRules),
	}
}

func mapACLTemplateInboundRulesDomainToAPI(rules []ACLTemplateInboundRule) []api.ACLTemplateInboundRule {
	transformed := make([]api.ACLTemplateInboundRule, len(rules))
	for i := range rules {
		transformed[i] = mapACLTemplateInboundRuleDomainToAPI(rules[i])
	}
	return transformed
}

func mapACLTemplateInboundRuleDomainToAPI(rule ACLTemplateInboundRule) api.ACLTemplateInboundRule {
	return api.ACLTemplateInboundRule{
		SrcType:  rule.SrcType,
		Protocol: rule.Protocol,
		SrcPort:  rule.SrcPort,
		DstPort:  rule.DstPort,
		FQDN:     rule.FQDN,
		Subnets:  rule.Subnets,
		SeqNO:    rule.SeqNo,
	}
}

func mapACLTemplateAPIToDomain(apiTemplate api.ACLTemplate) ACLTemplate {
	return ACLTemplate{
		UUID:            apiTemplate.UUID,
		Name:            apiTemplate.Name,
		Description:     apiTemplate.Description,
		MetroCode:       apiTemplate.MetroCode,
		DeviceUUID:      apiTemplate.VirtualDeviceUUID,
		DeviceACLStatus: apiTemplate.DeviceACLStatus,
		InboundRules:    mapACLTemplateInboundRulesAPIToDomain(apiTemplate.InboundRules),
	}
}

func mapACLTemplateInboundRulesAPIToDomain(apiRules []api.ACLTemplateInboundRule) []ACLTemplateInboundRule {
	transformed := make([]ACLTemplateInboundRule, len(apiRules))
	for i := range apiRules {
		transformed[i] = mapACLTemplateInboundRuleAPIToDomain(apiRules[i])
	}
	return transformed
}

func mapACLTemplateInboundRuleAPIToDomain(apiRule api.ACLTemplateInboundRule) ACLTemplateInboundRule {
	return ACLTemplateInboundRule{
		SrcType:  apiRule.SrcType,
		Protocol: apiRule.Protocol,
		SrcPort:  apiRule.SrcPort,
		DstPort:  apiRule.DstPort,
		FQDN:     apiRule.FQDN,
		Subnets:  apiRule.Subnets,
		SeqNo:    apiRule.SeqNO,
	}
}
