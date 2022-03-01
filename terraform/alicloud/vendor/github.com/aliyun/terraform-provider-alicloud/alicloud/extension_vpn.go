package alicloud

const (
	Ssl_Cert_Expiring = Status("expiring-soon")
	Ssl_Cert_Normal   = Status("normal")
	Ssl_Cert_Expired  = Status("expired")
)

const (
	IKE_VERSION_1       = string("ikev1")
	IKE_VERSION_2       = string("ikev2")
	IKE_MODE_MAIN       = string("main")
	IKE_MODE_AGGRESSIVE = string("aggressive")
	VPN_ENC_AES         = string("aes")
	VPN_ENC_AES_192     = string("aes192")
	VPN_ENC_AES_256     = string("aes256")
	VPN_ENC_AES_DES     = string("des")
	VPN_ENC_AES_3DES    = string("3des")
	VPN_AUTH_MD5        = string("md5")
	VPN_AUTH_SHA        = string("sha1")
	VPN_AUTH_SHA256     = string("sha256")
	VPN_AUTH_SHA386     = string("sha386")
	VPN_AUTH_SHA512     = string("sha512")
	VPN_PFS_DISABLED    = string("disabled")
	VPN_PFS_G1          = string("group1")
	VPN_PFS_G2          = string("group2")
	VPN_PFS_G5          = string("group5")
	VPN_PFS_G14         = string("group14")
	VPN_PFS_G24         = string("group24")
	VPN_UDP_PROTO       = string("UDP")
	VPN_TCP_PROTO       = string("TCP")
	SSL_VPN_ENC_AES_128 = string("AES-128-CBC")
	SSL_VPN_ENC_AES_192 = string("AES-192-CBC")
	SSL_VPN_ENC_AES_256 = string("AES-256-CBC")
	SSL_VPN_ENC_NONE    = string("none")
)

type IpsecConfig struct {
	IpsecAuthAlg  string
	IpsecEncAlg   string
	IpsecLifetime int
	IpsecPfs      string
}

type IkeConfig struct {
	IkeAuthAlg  string
	IkeEncAlg   string
	IkeLifetime int
	IkeMode     string
	IkePfs      string
	IkeVersion  string
	LocalId     string
	RemoteId    string
	Psk         string
}
