package ne

import (
	"context"
	"encoding/json"
	"fmt"
	"ne-go/internal/api"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8888"
)

var testPrimaryConnection = L2Connection{
	Name:                "name",
	ProfileUUID:         "profileUUID",
	Speed:               50,
	SpeedUnit:           "MB",
	Notifications:       []string{"someone@eu.equinix.com"},
	PurchaseOrderNumber: "orderNumber",
	VirtualDeviceUUID:   "a8d656a8-bbee-49df-9ea0-99b4d654eb6f",
	NamedTag:            "Private",
	SellerRegion:        "eu-central-1",
	SellerMetroCode:     "FR",
	AuthorizationKey:    "authorizationKey"}

func TestGetL2Connection(t *testing.T) {
	//Given
	respBody := api.GetConnectionResponse{}
	if err := readJSONData("./test-fixtures/ne_connection_get_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannot read test response due to %s", err.Error())
	}
	connID := "connId"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/ne/v1/l2/connections/%s", baseURL, connID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)

	//When
	neClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := neClient.GetL2Connection(connID)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyL2Connection(t, *conn, respBody)
}

func TestCreateL2Connection(t *testing.T) {
	//Given
	respBody := api.PostConnectionResponse{}
	if err := readJSONData("./test-fixtures/ne_connection_create_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	reqBody := api.PostConnectionRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/l2/connections", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()
	newConnection := testPrimaryConnection

	//When
	neClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := neClient.CreateL2Connection(newConnection)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyL2ConnectionRequest(t, *conn, reqBody)
	assert.Equal(t, conn.UUID, respBody.PrimaryConnectionID, "UUID matches")
}

func TestCreateRedundantL2Connection(t *testing.T) {
	//Given
	respBody := api.PostConnectionResponse{}
	if err := readJSONData("./test-fixtures/ne_connection_create_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannont read test response due to %s", err.Error())
	}
	reqBody := api.PostConnectionRequest{}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/ne/v1/l2/connections", baseURL),
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		},
	)
	defer httpmock.DeactivateAndReset()
	newPriConn := testPrimaryConnection
	newSecConn := L2Connection{
		Name:              "secName",
		VirtualDeviceUUID: "secondaryVirtualDeviceUUID",
		ZSidePortUUID:     "secondaryZSidePortUUID",
		ZSideVlanSTag:     717,
		ZSideVlanCTag:     718}

	//When
	neClient := NewClient(context.Background(), baseURL, testHc)
	conn, err := neClient.CreateL2RedundantConnection(newPriConn, newSecConn)

	//Then
	assert.Nil(t, err, "Client should not return an error")
	assert.NotNil(t, conn, "Client should return a response")
	verifyRedundantL2ConnectionRequest(t, newPriConn, newSecConn, reqBody)
	assert.Equal(t, conn.UUID, respBody.PrimaryConnectionID, "UUID matches")
	assert.Equal(t, conn.RedundantUUID, respBody.SecondaryConnectionID, "RedundantUUID matches")
}

func TestDeleteL2Connection(t *testing.T) {
	//Given
	respBody := api.DeleteConnectionResponse{}
	if err := readJSONData("./test-fixtures/ne_connection_delete_resp.json", &respBody); err != nil {
		assert.Failf(t, "Cannot read test response due to %s", err.Error())
	}
	connID := "connId"
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/ne/v1/l2/connections/%s", baseURL, connID),
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(200, respBody)
			return resp, nil
		})
	defer httpmock.DeactivateAndReset()

	//When
	ecxClient := NewClient(context.Background(), baseURL, testHc)
	err := ecxClient.DeleteL2Connection(connID)

	//Then
	assert.Nil(t, err, "Client should not return an error")
}

func verifyL2Connection(t *testing.T, conn L2Connection, resp api.GetConnectionResponse) {
	assert.Equal(t, resp.UUID, conn.UUID, "UUID matches")
	assert.Equal(t, resp.Name, conn.Name, "Name matches")
	assert.Equal(t, resp.SellerServiceUUID, conn.ProfileUUID, "Name matches")
	assert.Equal(t, int(resp.Speed), conn.Speed, "Speed matches")
	assert.Equal(t, resp.SpeedUnit, conn.SpeedUnit, "SpeedUnit matches")
	assert.Equal(t, resp.Status, conn.Status, "Status matches")
	assert.ElementsMatch(t, resp.Notifications, conn.Notifications, "Notifications match")
	assert.Equal(t, resp.PurchaseOrderNumber, conn.PurchaseOrderNumber, "PurchaseOrderNumber match")
	assert.Equal(t, resp.VirtualDeviceUUID, conn.VirtualDeviceUUID, "PrimaryVirtualDeviceUUID matches")
	assert.Equal(t, int(resp.VlanSTag), conn.VlanSTag, "PrimaryVlanSTag matches")
	assert.Equal(t, resp.NamedTag, conn.NamedTag, "NamedTag matches")
	assert.Equal(t, resp.ZSidePortUUID, conn.ZSidePortUUID, "PrimaryZSidePortUUID matches")
	assert.Equal(t, int(resp.ZSideVlanSTag), conn.ZSideVlanSTag, "PrimaryZSideVlanSTag matches")
	assert.Equal(t, int(resp.ZSideVlanCTag), conn.ZSideVlanCTag, "PrimaryZSideVlanCTag matches")
	assert.Equal(t, resp.SellerMetroCode, conn.SellerMetroCode, "SellerMetroCode matches")
	assert.Equal(t, resp.AuthorizationKey, conn.AuthorizationKey, "AuthorizationKey matches")
	assert.Equal(t, resp.RedundantUUID, conn.RedundantUUID, "RedundantUUID key matches")
}

func verifyL2ConnectionRequest(t *testing.T, conn L2Connection, req api.PostConnectionRequest) {
	assert.Equal(t, conn.Name, req.PrimaryName, "Name matches")
	assert.Equal(t, conn.ProfileUUID, req.ProfileUUID, "ProfileUUID matches")
	assert.Equal(t, int64(conn.Speed), req.Speed, "Speed matches")
	assert.Equal(t, conn.SpeedUnit, req.SpeedUnit, "SpeedUnit matches")
	assert.ElementsMatch(t, conn.Notifications, req.Notifications, "Notifications match")
	assert.Equal(t, conn.PurchaseOrderNumber, req.PurchaseOrderNumber, "PurchaseOrderNumber matches")
	assert.Equal(t, conn.VirtualDeviceUUID, req.VirtualDeviceUUID, "PrimaryVirtualDeviceUUID matches")
	assert.Equal(t, conn.NamedTag, req.NamedTag, "NamedTag matches")
	assert.Equal(t, conn.ZSidePortUUID, req.PrimaryZSidePortUUID, "PrimaryZSidePortUUID matches")
	assert.Equal(t, int64(conn.ZSideVlanSTag), req.PrimaryZSideVlanSTag, "PrimaryZSideVlanSTag matches")
	assert.Equal(t, int64(conn.ZSideVlanCTag), req.PrimaryZSideVlanCTag, "PrimaryZSideVlanCTag matches")
	assert.Equal(t, conn.SellerRegion, req.SellerRegion, "SellerRegion matches")
	assert.Equal(t, conn.SellerMetroCode, req.SellerMetroCode, "SellerMetroCode matches")
	assert.Equal(t, conn.AuthorizationKey, req.AuthorizationKey, "Authorization key matches")
}

func verifyRedundantL2ConnectionRequest(t *testing.T, primary L2Connection, secondary L2Connection, req api.PostConnectionRequest) {
	verifyL2ConnectionRequest(t, primary, req)
	assert.Equal(t, secondary.Name, req.SecondaryName, "SecondaryName matches")
	assert.Equal(t, int64(secondary.Speed), req.SecondarySpeed, "SecondarySpeed matches")
	assert.Equal(t, secondary.SpeedUnit, req.SecondarySpeedUnit, "SecondarySpeedUnit matches")
	assert.Equal(t, secondary.VirtualDeviceUUID, req.SecondaryVirtualDeviceUUID, "SecondaryVirtualDeviceUUID matches")
	assert.Equal(t, secondary.ZSidePortUUID, req.SecondaryZSidePortUUID, "SecondaryZSidePortUUID matches")
	assert.Equal(t, int64(secondary.ZSideVlanSTag), req.SecondaryZSideVlanSTag, "SecondaryZSideVlanSTag matches")
	assert.Equal(t, int64(secondary.ZSideVlanCTag), req.SecondaryZSideVlanCTag, "SecondaryZSideVlanCTag matches")
}
