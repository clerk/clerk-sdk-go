package clerk

import "encoding/json"

// CustomAttribute is used for both reading and writing a custom attribute, and
// the attribute's path has two names on the wire: the legacy scim_path and its
// replacement directory_path. Responses carry both, requests accept either, and
// a request that sends both with different values is rejected.
//
// One struct with two mutable fields for one value cannot tell which field the
// caller edited, so the SDK keeps a single field authoritative and marshals a
// single path key:
//
//   - On read, the path lands in SCIMPath. DirectoryPath stays nil, whichever
//     name the response used.
//   - On write, DirectoryPath wins when it is set and is sent as
//     directory_path; otherwise SCIMPath is sent as scim_path. The two names
//     are never both present, so a read-modify-write can never conflict.
//   - Clearing the path means leaving both nil.
type CustomAttribute struct {
	Name    *string `json:"name"`
	Key     *string `json:"key"`
	SSOPath *string `json:"sso_path"`
	// SCIMPath is the legacy name for DirectoryPath. It is the field the path
	// is decoded into, so existing code that reads it keeps working.
	//
	// Deprecated: use DirectoryPath on write. SCIMPath remains the field
	// populated on read until the legacy name is retired.
	SCIMPath *string `json:"scim_path,omitempty"`
	// DirectoryPath is the new name for SCIMPath. It is write-only: set it to
	// send the path as directory_path. Reads leave it nil and populate
	// SCIMPath instead.
	DirectoryPath *string `json:"directory_path,omitempty"`
	MultiValued   *bool   `json:"multi_valued,omitempty"`
}

// customAttributeJSON is the wire shape of a CustomAttribute. It carries both
// path names so that either can be read or written, while CustomAttribute
// exposes a single decoded path.
type customAttributeJSON struct {
	Name          *string `json:"name"`
	Key           *string `json:"key"`
	SSOPath       *string `json:"sso_path"`
	SCIMPath      *string `json:"scim_path,omitempty"`
	DirectoryPath *string `json:"directory_path,omitempty"`
	MultiValued   *bool   `json:"multi_valued,omitempty"`
}

// MarshalJSON emits exactly one path key: directory_path when DirectoryPath is
// set, scim_path otherwise. Sending both is what the API rejects when they
// disagree, so the SDK never sends both.
func (attr CustomAttribute) MarshalJSON() ([]byte, error) {
	out := customAttributeJSON{
		Name:        attr.Name,
		Key:         attr.Key,
		SSOPath:     attr.SSOPath,
		MultiValued: attr.MultiValued,
	}
	if attr.DirectoryPath != nil {
		out.DirectoryPath = attr.DirectoryPath
	} else {
		out.SCIMPath = attr.SCIMPath
	}
	return json.Marshal(out)
}

// UnmarshalJSON decodes the path from either name into SCIMPath and leaves
// DirectoryPath nil, so that a decoded attribute sent back unchanged carries
// one path key with the value it was read with.
func (attr *CustomAttribute) UnmarshalJSON(data []byte) error {
	var in customAttributeJSON
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	attr.Name = in.Name
	attr.Key = in.Key
	attr.SSOPath = in.SSOPath
	attr.MultiValued = in.MultiValued
	attr.SCIMPath = in.SCIMPath
	if in.DirectoryPath != nil {
		attr.SCIMPath = in.DirectoryPath
	}
	attr.DirectoryPath = nil

	return nil
}

type SAMLConnection struct {
	APIResource
	ID     string `json:"id"`
	Object string `json:"object"`
	Name   string `json:"name"`
	// Deprecated: Use `domains` instead.
	Domain                           string                         `json:"domain"`
	Domains                          []string                       `json:"domains"`
	IdpEntityID                      *string                        `json:"idp_entity_id"`
	OrganizationID                   *string                        `json:"organization_id"`
	EnterpriseConnectionID           *string                        `json:"enterprise_connection_id"`
	IdpSsoURL                        *string                        `json:"idp_sso_url"`
	IdpCertificate                   *string                        `json:"idp_certificate"`
	IdpCertificateIssuedAt           *int64                         `json:"idp_certificate_issued_at"`
	IdpCertificateExpiresAt          *int64                         `json:"idp_certificate_expires_at"`
	IdpMetadataURL                   *string                        `json:"idp_metadata_url"`
	IdpMetadata                      *string                        `json:"idp_metadata"`
	AcsURL                           string                         `json:"acs_url"`
	SPEntityID                       string                         `json:"sp_entity_id"`
	SPMetadataURL                    string                         `json:"sp_metadata_url"`
	AttributeMapping                 SAMLConnectionAttributeMapping `json:"attribute_mapping"`
	Active                           bool                           `json:"active"`
	Provider                         string                         `json:"provider"`
	UserCount                        int64                          `json:"user_count"`
	SyncUserAttributes               bool                           `json:"sync_user_attributes"`
	AllowSubdomains                  bool                           `json:"allow_subdomains"`
	AllowIdpInitiated                bool                           `json:"allow_idp_initiated"`
	DisableAdditionalIdentifications bool                           `json:"disable_additional_identifications"`
	ForceAuthn                       bool                           `json:"force_authn"`
	DisableJITProvisioning           *bool                          `json:"disable_jit_provisioning,omitempty"`
	AllowOrganizationAccountLinking  bool                           `json:"allow_organization_account_linking"`
	CustomAttributes                 *[]CustomAttribute             `json:"custom_attributes,omitempty"`
	Authenticatable                  bool                           `json:"authenticatable"`
	LoginHint                        *LoginHint                     `json:"login_hint,omitempty"`
	CreatedAt                        int64                          `json:"created_at"`
	UpdatedAt                        int64                          `json:"updated_at"`
}

type SAMLConnectionAttributeMapping struct {
	UserID       string `json:"user_id"`
	EmailAddress string `json:"email_address"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// LoginHint describes the login hint a SAML connection sends to the identity
// provider when initiating a sign-in.
type LoginHint struct {
	Mode   string  `json:"mode"`
	Source *string `json:"source,omitempty"`
}

type SAMLConnectionList struct {
	APIResource
	SAMLConnections []*SAMLConnection `json:"data"`
	TotalCount      int64             `json:"total_count"`
}
