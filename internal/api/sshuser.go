package api

//SSHUser describes network edge SSH user
type SSHUser struct {
	UUID        *string  `json:"uuid,omitempty"`
	Username    *string  `json:"username,omitempty"`
	Password    *string  `json:"password,omitempty"`
	DeviceUUIDs []string `json:"deviceUUIDs,omitempty"`
}

//SSHUserRequest describes network edge SSH user creation request
type SSHUserRequest struct {
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	DeviceUUID *string `json:"deviceUuid,omitempty"`
}

//SSHUserRequestResponse describes response for SSH user creation request
type SSHUserRequestResponse struct {
	UUID *string `json:"uuid,omitempty"`
}

//SSHUserUpdateRequest describes network edge SSH user update request
type SSHUserUpdateRequest struct {
	Password *string `json:"password,omitempty"`
}

//SSHUsersResponse describes response for a get ssh user list request
type SSHUsersResponse struct {
	TotalCount *int      `json:"totalCount,omitempty"`
	PageSize   *int      `json:"pageSize,omitempty"`
	PageNumber *int      `json:"pageNumber,omitempty"`
	List       []SSHUser `json:"list,omitempty"`
}
