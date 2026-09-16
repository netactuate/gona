package gona

// DeleteServerRequest contains optional fields for deleting a cloud server.
type DeleteServerRequest struct {
	CancelBilling *bool   `json:"cancel_billing,omitempty"`
	ForcePassword *string `json:"force_password,omitempty"`
	Password      *string `json:"password,omitempty"`
}

// ServerActionRequest contains optional fields for cloud server power actions.
type ServerActionRequest struct {
	Force *bool `json:"force,omitempty"`
}
