package api

type DeviceLinkGroup struct {
	UUID           *string                    `json:"uuid,omitempty"`
	GroupName      *string                    `json:"groupName,omitempty"`
	Subnet         *string                    `json:"subnet,omitempty"`
	Status         *string                    `json:"status,omitempty"`
	Devices        []DeviceLinkGroupDevice    `json:"linkDevices,omitempty"`
	Links          []DeviceLinkGroupLink      `json:"links"`
	ProjectID      *string                    `json:"projectId,omitempty"`
	MetroLinks     []DeviceLinkGroupMetroLink `json:"metroLinks"`
	RedundancyType *string                    `json:"redundancyType,omitempty"`
}

type DeviceLinkGroupUpdateRequest struct {
	GroupName      *string                    `json:"groupName,omitempty"`
	Subnet         *string                    `json:"subnet,omitempty"`
	Devices        []DeviceLinkGroupDevice    `json:"linkDevices,omitempty"`
	Links          []DeviceLinkGroupLink      `json:"links,omitempty"`
	MetroLinks     []DeviceLinkGroupMetroLink `json:"metroLinks,omitempty"`
	RedundancyType *string                    `json:"redundancyType,omitempty"`
}

type DeviceLinkGroupDevice struct {
	DeviceUUID  *string `json:"deviceUuid,omitempty"`
	ASN         *int    `json:"asn,omitempty"`
	InterfaceID *int    `json:"interfaceId,omitempty"`
	Status      *string `json:"status,omitempty"`
	IPAddress   *string `json:"ipAssigned,omitempty"`
}

type DeviceLinkGroupLink struct {
	AccountNumber        *string `json:"accountNumber,omitempty"`
	Throughput           *string `json:"throughput,omitempty"`
	ThroughputUnit       *string `json:"throughputUnit,omitempty"`
	SourceMetroCode      *string `json:"sourceMetroCode,omitempty"`
	DestinationMetroCode *string `json:"destinationMetroCode,omitempty"`
	SourceZoneCode       *string `json:"sourceZoneCode,omitempty"`
	DestinationZoneCode  *string `json:"destinationZoneCode,omitempty"`
}
type DeviceLinkGroupMetroLink struct {
	AccountNumber  *string `json:"accountNumber,omitempty"`
	MetroCode      *string `json:"metroCode,omitempty"`
	Throughput     *string `json:"throughput,omitempty"`
	ThroughputUnit *string `json:"throughputUnit,omitempty"`
}

type DeviceLinkGroupCreateResponse struct {
	UUID *string `json:"uuid,omitempty"`
}

type DeviceLinkGroupsGetResponse struct {
	Pagination Pagination        `json:"pagination,omitempty"`
	Data       []DeviceLinkGroup `json:"data,omitempty"`
}
