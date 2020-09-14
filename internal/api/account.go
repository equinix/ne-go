package api

type AccountResponse struct {
	Accounts []Account `json:"accounts,omitempty"`
}

type Account struct {
	Name   string `json:"accountName,omitempty"`
	Number string `json:"accountNumber,omitempty"`
	UCMID  string `json:"accountUcmId,omitempty"`
	Status string `json:"accountStatus,omitempty"`
}
