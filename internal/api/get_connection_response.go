package api

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// GetConnectionResponse get connection by UUID response
type GetConnectionResponse struct {

	// aside encapsulation
	AsideEncapsulation string `json:"asideEncapsulation,omitempty"`

	// authorization key
	AuthorizationKey string `json:"authorizationKey,omitempty"`

	// billing tier
	BillingTier string `json:"billingTier,omitempty"`

	// buyer organization name
	BuyerOrganizationName string `json:"buyerOrganizationName,omitempty"`

	// created by
	CreatedBy string `json:"createdBy,omitempty"`

	// created by email
	CreatedByEmail string `json:"createdByEmail,omitempty"`

	// created by full name
	CreatedByFullName string `json:"createdByFullName,omitempty"`

	// created date
	CreatedDate string `json:"createdDate,omitempty"`

	// last updated by
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`

	// last updated by email
	LastUpdatedByEmail string `json:"lastUpdatedByEmail,omitempty"`

	// last updated by full name
	LastUpdatedByFullName string `json:"lastUpdatedByFullName,omitempty"`

	// last updated date
	LastUpdatedDate string `json:"lastUpdatedDate,omitempty"`

	// metro code
	MetroCode string `json:"metroCode,omitempty"`

	// metro description
	MetroDescription string `json:"metroDescription,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// named tag
	NamedTag string `json:"namedTag,omitempty"`

	// notifications
	Notifications []string `json:"notifications"`

	// private
	Private bool `json:"private,omitempty"`

	// provider status
	ProviderStatus string `json:"providerStatus,omitempty"`

	// purchase order number
	PurchaseOrderNumber string `json:"purchaseOrderNumber,omitempty"`

	// redundancy group
	RedundancyGroup string `json:"redundancyGroup,omitempty"`

	// redundancy type
	RedundancyType string `json:"redundancyType,omitempty"`

	// redundant UUID
	RedundantUUID string `json:"redundantUUID,omitempty"`

	// remote
	Remote bool `json:"remote,omitempty"`

	// self
	Self bool `json:"self,omitempty"`

	// seller region
	SellerRegion string `json:"sellerRegion,omitempty"`

	// seller metro code
	SellerMetroCode string `json:"sellerMetroCode,omitempty"`

	// seller metro description
	SellerMetroDescription string `json:"sellerMetroDescription,omitempty"`

	// seller organization name
	SellerOrganizationName string `json:"sellerOrganizationName,omitempty"`

	// seller service name
	SellerServiceName string `json:"sellerServiceName,omitempty"`

	// seller service UUID
	SellerServiceUUID string `json:"sellerServiceUUID,omitempty"`

	// speed
	Speed int64 `json:"speed,omitempty"`

	// speed unit
	SpeedUnit string `json:"speedUnit,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// uuid
	UUID string `json:"uuid,omitempty"`

	// virtual device UUID
	VirtualDeviceUUID string `json:"virtualDeviceUUID,omitempty"`

	// vlan s tag
	VlanSTag int64 `json:"vlanSTag,omitempty"`

	// z side port name
	ZSidePortName string `json:"zSidePortName,omitempty"`

	// z side port UUID
	ZSidePortUUID string `json:"zSidePortUUID,omitempty"`

	// z side vlan c tag
	ZSideVlanCTag int64 `json:"zSideVlanCTag,omitempty"`

	// z side vlan s tag
	ZSideVlanSTag int64 `json:"zSideVlanSTag,omitempty"`

	// action details
	ActionDetails []ActionDetail `json:"actionDetails,omitempty"`
}

// ActionDetail required actions details 
type ActionDetail struct {

	// Action Type
	ActionType string `json:"actionType,omitempty"`

	// Operation Id
	OperationID string `json:"operationId,omitempty"`

	// Action Message
	ActionMessage string `json:"actionMessage,omitempty"`

	// Action Required
	ActionRequiredData []ActionRequiredData `json:"actionRequiredData,omitempty"`
}


// ActionRequiredData required actions data
type ActionRequiredData struct {

	// Key
	Key string `json:"key,omitempty"`

	// Label
	Label string `json:"label,omitempty"`

	// Value
	Value string `json:"value,omitempty"`
	
	// Editable
	Editable bool `json:"editable,omitempty"`

	// Validation Pattern
	ValidationPattern string `json:"validationPattern,omitempty"`
}


// Validate validates this g e t connection by u Uid response
func (m *GetConnectionResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *GetConnectionResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *GetConnectionResponse) UnmarshalBinary(b []byte) error {
	var res GetConnectionResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
