package ne

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/equinix/ne-go/internal/api"
	"github.com/go-resty/resty/v2"
)

//GetL2Connection operation retrieves layer 2 connection with a given UUID
func (c RestClient) GetL2Connection(uuid string) (*L2Connection, error) {
	url := fmt.Sprintf("%s/ne/v1/l2/connections/%s", c.baseURL, url.PathEscape(uuid))
	respBody := api.GetConnectionResponse{}
	req := c.R().SetResult(&respBody)
	if err := c.execute(req, resty.MethodGet, url); err != nil {
		return nil, err
	}
	return mapGETToL2Connection(respBody), nil
}

//CreateL2Connection operation creates non-redundant layer 2 connection with a given connection structure.
//Upon successful creation, connection structure, enriched with assigned UUID, will be returned
func (c RestClient) CreateL2Connection(l2connection L2Connection) (*L2Connection, error) {
	url := fmt.Sprintf("%s/ne/v1/l2/connections", c.baseURL)
	reqBody := createL2ConnectionRequest(l2connection)
	respBody := api.PostConnectionResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.execute(req, resty.MethodPost, url); err != nil {
		return nil, err
	}
	l2connection.UUID = respBody.PrimaryConnectionID
	return &l2connection, nil
}

//CreateL2RedundantConnection operation creates redundant layer2 connection with given connection structures.
//Primary connection structure is used as a baseline for underlaying API call, whereas secondary connection structure provices
//supplementary information only.
//Upon successful creation, primary connection structure, enriched with assigned UUID and redundant connection UUID, will be returned
func (c RestClient) CreateL2RedundantConnection(primary L2Connection, secondary L2Connection) (*L2Connection, error) {
	url := fmt.Sprintf("%s/ne/v1/l2/connections", c.baseURL)
	reqBody := createL2RedundantConnectionRequest(primary, secondary)
	respBody := api.PostConnectionResponse{}
	req := c.R().SetBody(&reqBody).SetResult(&respBody)
	if err := c.execute(req, resty.MethodPost, url); err != nil {
		return nil, err
	}
	primary.UUID = respBody.PrimaryConnectionID
	primary.RedundantUUID = respBody.SecondaryConnectionID
	return &primary, nil
}

//DeleteL2Connection deletes layer 2 connection with a given UUID
func (c RestClient) DeleteL2Connection(uuid string) error {
	url := fmt.Sprintf("%s/ne/v1/l2/connections/%s", c.baseURL, url.PathEscape(uuid))
	respBody := api.DeleteConnectionResponse{}
	req := c.R().SetResult(&respBody)
	if err := c.execute(req, resty.MethodDelete, url); err != nil {
		return err
	}
	return nil
}

func mapGETToL2Connection(getResponse api.GetConnectionResponse) *L2Connection {
	return &L2Connection{
		UUID:                     getResponse.UUID,
		AuthorizationKey:         getResponse.AuthorizationKey,
		Name:                     getResponse.Name,
		NamedTag:                 getResponse.NamedTag,
		Notifications:            getResponse.Notifications,
		ProfileUUID:              getResponse.SellerServiceUUID,
		Speed:                    int(getResponse.Speed),
		SpeedUnit:                getResponse.SpeedUnit,
		Status:                   getResponse.Status,
		PurchaseOrderNumber:      getResponse.PurchaseOrderNumber,
		VirtualDeviceUUID:        getResponse.VirtualDeviceUUID,
		VlanSTag:                 int(getResponse.VlanSTag),
		ZSidePortUUID:            getResponse.ZSidePortUUID,
		ZSideVlanSTag:            int(getResponse.ZSideVlanSTag),
		ZSideVlanCTag:            int(getResponse.ZSideVlanCTag),
		SellerRegion:             getResponse.SellerRegion,
		SellerMetroCode:          getResponse.SellerMetroCode,
		SellerHostedConnectionID: getHostedConnectionID(getResponse),
		RedundantUUID:            getResponse.RedundantUUID}
}

func createL2ConnectionRequest(l2connection L2Connection) api.PostConnectionRequest {
	return api.PostConnectionRequest{
		AuthorizationKey:     l2connection.AuthorizationKey,
		NamedTag:             l2connection.NamedTag,
		Notifications:        l2connection.Notifications,
		PrimaryName:          l2connection.Name,
		PrimaryZSidePortUUID: l2connection.ZSidePortUUID,
		PrimaryZSideVlanSTag: int64(l2connection.ZSideVlanSTag),
		PrimaryZSideVlanCTag: int64(l2connection.ZSideVlanCTag),
		ProfileUUID:          l2connection.ProfileUUID,
		PurchaseOrderNumber:  l2connection.PurchaseOrderNumber,
		SellerRegion:         l2connection.SellerRegion,
		SellerMetroCode:      l2connection.SellerMetroCode,
		Speed:                int64(l2connection.Speed),
		SpeedUnit:            l2connection.SpeedUnit,
		VirtualDeviceUUID:    l2connection.VirtualDeviceUUID}
}

func createL2RedundantConnectionRequest(primary L2Connection, secondary L2Connection) api.PostConnectionRequest {
	connReq := createL2ConnectionRequest(primary)
	connReq.SecondaryAuthorizationKey = secondary.AuthorizationKey
	connReq.SecondaryName = secondary.Name
	connReq.SecondaryNotifications = secondary.Notifications
	connReq.SecondaryProfileUUID = secondary.ProfileUUID
	connReq.SecondarySellerMetroCode = secondary.SellerMetroCode
	connReq.SecondarySellerRegion = secondary.SellerRegion
	connReq.SecondarySpeed = int64(secondary.Speed)
	connReq.SecondarySpeedUnit = secondary.SpeedUnit
	connReq.SecondaryVirtualDeviceUUID = secondary.VirtualDeviceUUID
	connReq.SecondaryZSidePortUUID = secondary.ZSidePortUUID
	connReq.SecondaryZSideVlanSTag = int64(secondary.ZSideVlanSTag)
	connReq.SecondaryZSideVlanCTag = int64(secondary.ZSideVlanCTag)
	return connReq
}

// function to obtain the SP side connection ID when it is included in the response
func getHostedConnectionID(getResponse api.GetConnectionResponse) string {
	actionDetails := getResponse.ActionDetails
	for _, v := range actionDetails {
		for _, d := range v.ActionRequiredData {
			if strings.Contains(
				strings.ToLower(d.Key),
				strings.ToLower("ConnectionId")) {
				return d.Value
			}
		}
	}
	return ""
}
