## 1.1.0 (Unreleased)

FEATURES:

* Equinix Network Edge device links support ([equinix/terraform-provider-equinix#43](https://github.com/equinix/terraform-provider-equinix/issues/43))

## 1.0.1 (March 15, 2021)

BUG FIXES:

* fixed errorCode consts for failed removal of device and SSH public key when
objects were already removed

## 1.0.0 (March 12, 2021)

NOTES:

* first version of Equinix Network Edge Go client

FEATURES:

* Equinix Network Edge `Device` management
  * `CreateDevice` function to provision new device single device
  * `CreateRedundantDevice` function to provision pair of redundant devices
  * `GetDevices` function to fetch details all devices
  * `GetDevice` function fo fetch details of a given device
  * `NewDeviceUpdateRequest` function to create device update request, with option
  to update:
    * device name
    * term length
    * notification addresses
    * additional bandwidth amount
    * ACL template
  * `DeleteDevice` function to deprovision given device
* Uploading Equinix Network Edge device license files
  * `UploadLicenseFile` to upload file from a given `io.Reader`
* Equinix Network Edge `SSHUser` management
  * `CreateSSHUser` function to create new SSH user
  * `GetSSHUsers` function to fetch details all SSH users
  * `GetSSHUser` function to fetch details of a given SSH user
  * `NewSSHUserUpdateRequest` function to create device update request, with option
  to update:
    * password
    * associated devices
  * `DeleteSSHUser` to remove SSH user
* Equinix Network Edge `SSHPublicKey` management
  * `CreateSSHPublicKey` function to create new SSH public key
  * `GetSSHPublicKeys` function to fetch details all SSH public keys
  * `GetSSHPublicKey` function to fetch details of a given SSH public key
  * `DeleteSSHPublicKey` to remove SSH public key
* Equinix Network Edge `ACLTemplate` management
  * `CreateACLTemplate` function to create new ACL template
  * `GetACLTemplates` function to fetch details of all ACL templates
  * `GetACLTemplate` function to fetch details of a given ACL template
  * `ReplaceACLTemplate` function to replace given ACL template
  * `DeleteACLTemplate` function to remove given ACL template
* Equinix Network Edge `BGPConfiguration` management
  * `CreateBGPConfiguration` function to create new BGP configuration
  * `GetBGPConfiguration` function to fetch details of a given BGP configuration
  * `GetBGPConfigurationForConnection` function to fetch details of a BGP configuration
  for a given connection
  * `NewBGPConfigurationUpdateRequest` function to create BGP configuration update
  request, with option to update:
    * local IP address
    * local ASN number
    * remote ASN number
    * remote IP address
    * authentication key
