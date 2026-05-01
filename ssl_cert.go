package unifi

import "fmt"

// SSLCertificate represents the active SSL certificate from /api/s/{site}/stat/active.
type SSLCertificate struct {
	ID          string                `fake:"{uuid}"      json:"_id"`
	CertType    string                `json:"cert_type"`
	Chain       []SSLCertificateChain `json:"chain"`
	Fingerprint string                `json:"fingerprint"`
	IsActive    FlexBool              `json:"is_active"`
	IsValid     FlexBool              `json:"is_valid"`
	Issuer      string                `json:"issuer"`
	Serial      string                `json:"serial"`
	Status      string                `json:"status"`
	Subject     string                `json:"subject"`
	ValidFrom   FlexInt               `json:"valid_from"`
	ValidTo     FlexInt               `json:"valid_to"`

	SiteName string `json:"-"`
}

// SSLCertificateChain represents a certificate in the chain (root, intermediate, etc.).
type SSLCertificateChain struct {
	CertType    string  `json:"cert_type"`
	Fingerprint string  `json:"fingerprint"`
	Issuer      string  `json:"issuer"`
	Serial      string  `json:"serial"`
	Subject     string  `json:"subject"`
	ValidFrom   FlexInt `json:"valid_from"`
	ValidTo     FlexInt `json:"valid_to"`
}

// GetSSLCertificate returns the active SSL certificate info for a single site.
// Uses the legacy API endpoint: GET /api/s/{site}/stat/active.
// When the controller returns no certificate data, a zero-value SSLCertificate is returned
// with only SiteName set. Callers can detect this case by checking cert.ID == "".
func (u *Unifi) GetSSLCertificate(site *Site) (*SSLCertificate, error) {
	if u == nil {
		return nil, ErrNilUnifi
	}

	if site == nil || site.Name == "" {
		return nil, ErrNoSiteProvided
	}

	u.DebugLog("Polling Controller for SSL certificate, site %s", site.Name)

	path := fmt.Sprintf(APISSLCertPath, site.Name)

	var response struct {
		Data []SSLCertificate `json:"data"`
	}

	if err := u.GetData(path, &response); err != nil {
		return nil, fmt.Errorf("fetching SSL certificate for site %s: %w", site.Name, err)
	}

	if len(response.Data) == 0 {
		u.DebugLog("No SSL certificate found for site %s", site.Name)

		return &SSLCertificate{SiteName: site.SiteName}, nil
	}

	if len(response.Data) > 1 {
		u.DebugLog("SSL certificate response for site %s contained %d entries; using first", site.Name, len(response.Data))
	}

	response.Data[0].SiteName = site.SiteName

	return &response.Data[0], nil
}
