package client

// CustomScript models a /api/v1/library/custom-scripts response. The self-service
// category and recommended fields are accepted on write but NOT returned by GET (see
// API_REFERENCE.md), so they decode as zero values.
type CustomScript struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name"`
	ExecutionFrequency     string `json:"execution_frequency,omitempty"`
	Script                 string `json:"script,omitempty"`
	RemediationScript      string `json:"remediation_script,omitempty"`
	ShowInSelfService      bool   `json:"show_in_self_service"`
	SelfServiceCategoryID  string `json:"self_service_category_id,omitempty"`
	SelfServiceRecommended bool   `json:"self_service_recommended,omitempty"`
	Active                 bool   `json:"active"`
	Restart                bool   `json:"restart"`
}

// CustomScriptInput is the create/update payload. A nil field is left out so the
// server keeps its value; a non-nil zero value ("" or false) is sent, so fields can
// be cleared or switched off.
type CustomScriptInput struct {
	Name                   string  `json:"name"`
	ExecutionFrequency     string  `json:"execution_frequency"`
	Script                 string  `json:"script"`
	RemediationScript      *string `json:"remediation_script,omitempty"`
	ShowInSelfService      *bool   `json:"show_in_self_service,omitempty"`
	SelfServiceCategoryID  *string `json:"self_service_category_id,omitempty"`
	SelfServiceRecommended *bool   `json:"self_service_recommended,omitempty"`
	Active                 *bool   `json:"active,omitempty"`
	Restart                *bool   `json:"restart,omitempty"`
}

// Tag models /api/v1/tags. There is no GET-by-id endpoint; Read lists and filters.
type Tag struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// SelfServiceCategory models an item of the flat array returned by
// /api/v1/self-service/categories.
type SelfServiceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CustomProfile models /api/v1/library/custom-profiles GET responses. The input
// .mobileconfig is uploaded as multipart `file`; the API returns it as `profile`.
type CustomProfile struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	MDMIdentifier string `json:"mdm_identifier,omitempty"`
	Profile       string `json:"profile,omitempty"`
	Active        bool   `json:"active"`
	RunsOnMac     bool   `json:"runs_on_mac"`
	RunsOnIPhone  bool   `json:"runs_on_iphone"`
	RunsOnIPad    bool   `json:"runs_on_ipad"`
	RunsOnTV      bool   `json:"runs_on_tv"`
	RunsOnVision  bool   `json:"runs_on_vision"`
}
