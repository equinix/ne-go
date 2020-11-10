package api

//SSHPublicKey describes SSH public key
type SSHPublicKey struct {
	UUID     string `json:"uuid,omitempty"`
	KeyName  string `json:"keyName,omitempty"`
	KeyValue string `json:"keyValue,omitempty"`
}

//SSHPublicKeyCreateResponse describes SSH public key create response
type SSHPublicKeyCreateResponse struct {
	UUID string `json:"uuid,omitempty"`
}
