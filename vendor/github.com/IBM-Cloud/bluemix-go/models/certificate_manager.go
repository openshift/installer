package models

// CertificateInfo struct for cert-import & cert-reimport success response.
type CertificateInfo struct {
	ID           string                   `json:"_id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	Domains      []string                 `json:"domains"`
	RotateKeys   bool                     `json:"rotate_keys"`
	Status       string                   `json:"status"`
	Issuer       string                   `json:"issuer"`
	BeginsOn     int64                    `json:"begins_on"`
	ExpiresOn    int64                    `json:"expires_on"`
	Algorithm    string                   `json:"algorithm"`
	KeyAlgorithm string                   `json:"key_algorithm"`
	Imported     bool                     `json:"imported"`
	HasPrevious  bool                     `json:"has_previous"`
	IssuanceInfo *CertificateIssuanceInfo `json:"issuance_info"`
	SerialNumber string                   `json:"serial_number,omitempty"`
	OrderPolicy  OrderPolicy              `json:"order_policy,omitempty"`
}

//CertificateIssuanceInfo struct
type CertificateIssuanceInfo struct {
	Status         string `json:"status"`
	Code           string `json:"code"`
	AdditionalInfo string `json:"additional_info"`
	Auto           bool   `json:"auto"`
	OrderedOn      int64  `json:"ordered_on"`
}

// CertificateImportData struct for holding user-provided certificates and keys for cert-import.
type CertificateImportData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        Data   `json:"data"`
}

//Data of Imported Certificate
type Data struct {
	Content                 string `json:"content"`
	Privatekey              string `json:"priv_key,omitempty"`
	IntermediateCertificate string `json:"intermediate,omitempty"`
}

// CertificateDelete struct for cert-delete success response.
type CertificateDelete struct {
	Message string
}

// CertificateMetadataUpdate struct for cert-metadata-update's request body.
type CertificateMetadataUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CertificateReimportData struct for holding user-provided certificates and keys for cert-reimport.
type CertificateReimportData struct {
	Content                 string `json:"content"`
	Privatekey              string `json:"priv_key,omitempty"`
	IntermediateCertificate string `json:"intermediate,omitempty"`
}

//CertificateGetData ...
type CertificateGetData struct {
	ID           string                  `json:"_id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Domains      []string                `json:"domains"`
	Status       string                  `json:"status"`
	Issuer       string                  `json:"issuer"`
	BeginsOn     int64                   `json:"begins_on"`
	ExpiresOn    int64                   `json:"expires_on"`
	Algorithm    string                  `json:"algorithm"`
	KeyAlgorithm string                  `json:"key_algorithm"`
	Imported     bool                    `json:"imported"`
	HasPrevious  bool                    `json:"has_previous"`
	IssuanceInfo CertificateIssuanceInfo `json:"issuance_info"`
	Data         *Data                   `json:"data"`
	DataKeyID    string                  `json:"data_key_id"`
}

// CertificateOrderData struct for holding user-provided order data for cert-order.
type CertificateOrderData struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description,omitempty"`
	Domains                []string `json:"domains"`
	DomainValidationMethod string   `json:"domain_validation_method"`
	DNSProviderInstanceCrn string   `json:"dns_provider_instance_crn,omitempty"`
	Issuer                 string   `json:"issuer,omitempty"`
	Algorithm              string   `json:"algorithm,omitempty"`
	KeyAlgorithm           string   `json:"key_algorithm,omitempty"`
	AutoRenewEnabled       bool     `json:"auto_renew_enabled,omitempty"`
}

// CertificateRenewData struct for holding user-provided renew data for cert-renew.
type CertificateRenewData struct {
	RotateKeys bool `json:"rotate_keys"`
}

//CertificatesInfo List of certificates
type CertificatesInfo struct {
	CertificateList []CertificateInfo `json:"certificates"`
}

//OrderPolicy ...
type OrderPolicy struct {
	Name             string `json:"name,omitempty"`
	AutoRenewEnabled bool   `json:"auto_renew_enabled,omitempty"`
}
