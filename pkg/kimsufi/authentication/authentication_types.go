package authentication

type CurrentCredentialResponse struct {
	ApplicationID int                     `json:"applicationId"`
	Creation      string                  `json:"creation"`
	CredentialID  int                     `json:"credentialId"`
	Expiration    string                  `json:"expiration"`
	LastUse       string                  `json:"lastUse"`
	OVHSupport    bool                    `json:"ovhSupport"`
	Rules         []CurrentCredentialRule `json:"rules"`
	Status        string                  `json:"status"`
}

type CurrentCredentialRule struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
