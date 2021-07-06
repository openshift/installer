package datatypes

const (
	CONFIG = "config"
	STAT = "stat"

)

// Base
type BaseRes struct {
	Errorcode *int    `json:"errorcode,omitempty"`
	Message   *string `json:"message,omitempty"`
	Severity  *string `json:"severity,omitempty"`
}

// service
type Service struct {
	Name        *string `json:"name,omitempty"`
	Ip          *string `json:"ip,omitempty"`
	Ipaddress   *string `json:"ipaddress,omitempty"`
	ServiceType *string `json:"servicetype,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Maxclient   *string `json:"maxclient,omitempty"`
	Usip        *string `json:"usip,omitempty"`
}

type ServiceReq struct {
	Service *Service `json:"service,omitempty"`
}

type ServiceRes struct {
	BaseRes
	Service []Service `json:"service,omitempty"`
}

// lbvserver
type Lbvserver struct {
	Name        *string `json:"name,omitempty"`
	ServiceType *string `json:"servicetype,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Lbmethod    *string `json:"lbmethod,omitempty"`
	Ipv46       *string `json:"ipv46,omitempty"`
	Persistencetype       *string `json:"persistencetype,omitempty"`
}

type LbvserverReq struct {
	Lbvserver *Lbvserver `json:"lbvserver,omitempty"`
}

type LbvserverRes struct {
	BaseRes
	Lbvserver []Lbvserver `json:"lbvserver,omitempty"`
}

//lbvserver_service_binding
type LbvserverServiceBinding struct {
	Name        *string `json:"name,omitempty"`
	ServiceName *string `json:"serviceName,omitempty"`
}

type LbvserverServiceBindingReq struct {
	LbvserverServiceBinding *LbvserverServiceBinding `json:"lbvserver_service_binding,omitempty"`
}

type LbvserverServiceBindingRes struct {
	BaseRes
	LbvserverServiceBinding []LbvserverServiceBinding `json:"lbvserver_service_binding,omitempty"`
}

// systemfile
type Systemfile struct {
	Filename     *string `json:"filename,omitempty"`
	Filelocation *string `json:"filelocation,omitempty"`
	Filecontent  *string `json:"filecontent,omitempty"`
	Fileencoding *string `json:"fileencoding,omitempty"`
}

type SystemfileReq struct {
	Systemfile *Systemfile `json:"systemfile,omitempty"`
}

type SystemfileRes struct {
	BaseRes
	Systemfile []Systemfile `json:"systemfile,omitempty"`
}

// nsfeature
type Nsfeature struct {
	Feature []string `json:"feature"`
}

type NsfeatureReq struct {
	Nsfeature *Nsfeature `json:"nsfeature,omitempty"`
}

type NsfeatureRes struct {
	BaseRes
	Nsfeature []Nsfeature `json:"nsfeature,omitempty"`
}

// sslcertkey
type Sslcertkey struct {
	Certkey *string `json:"certkey,omitempty"`
	Cert    *string `json:"cert,omitempty"`
	Key     *string `json:"key,omitempty"`
}

type SslcertkeyReq struct {
	Sslcertkey *Sslcertkey `json:"sslcertkey,omitempty"`
}

type SslcertkeyRes struct {
	BaseRes
	Sslcertkey []Sslcertkey `json:"sslcertkey,omitempty"`
}

// sslvserver_sslcertkey_binding
type SslvserverSslcertkeyBinding struct {
	Vservername *string `json:"vservername,omitempty"`
	Certkeyname    *string `json:"certkeyname,omitempty"`
}

type SslvserverSslcertkeyBindingReq struct {
	SslvserverSslcertkeyBinding *SslvserverSslcertkeyBinding `json:"sslvserver_sslcertkey_binding,omitempty"`
}

type SslvserverSslcertkeyBindingRes struct {
	BaseRes
	SslvserverSslcertkeyBinding []SslvserverSslcertkeyBinding `json:"sslvserver_sslcertkey_binding,omitempty"`
}

// systemuser
type Systemuser struct {
	Username *string `json:"username,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type SystemuserReq struct {
	Systemuser *Systemuser `json:"systemuser,omitempty"`
}

type SystemuserRes struct {
	BaseRes
	Systemuser []Systemuser `json:"systemuser,omitempty"`
}

// hanode
type Hanode struct {
	Id *string `json:"id,omitempty"`
	Ipaddress    *string `json:"ipaddress,omitempty"`
	Hastatus    *string `json:"hastatus,omitempty"`
}

type HanodeReq struct {
	Hanode *Hanode `json:"hanode,omitempty"`
}

type HanodeRes struct {
	BaseRes
	Hanode []Hanode `json:"hanode,omitempty"`
}

// nsrpcnode
type Nsrpcnode struct {
	Ipaddress *string `json:"ipaddress,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type NsrpcnodeReq struct {
	Nsrpcnode *Nsrpcnode `json:"nsrpcnode,omitempty"`
}

type NsrpcnodeRes struct {
	BaseRes
	Nsrpcnode []Nsrpcnode `json:"nsrpcnode,omitempty"`
}

// hafiles
type Hafiles struct {
	Mode []string `json:"mode,omitempty"`
}

type HafilesReq struct {
	Hafiles *Hafiles `json:"hafiles,omitempty"`
}

type HafilesRes struct {
	BaseRes
	Hafiles []Hafiles `json:"hafiles,omitempty"`
}

//service_lbmonitor_binding
type ServiceLbmonitorBinding struct {
	Name        *string `json:"name,omitempty"`
	MonitorName *string `json:"monitor_name,omitempty"`
}

type ServiceLbmonitorBindingReq struct {
	ServiceLbmonitorBinding *ServiceLbmonitorBinding `json:"service_lbmonitor_binding,omitempty"`
}

type ServiceLbmonitorBindingRes struct {
	BaseRes
	ServiceLbmonitorBinding []ServiceLbmonitorBinding `json:"service_lbmonitor_binding,omitempty"`
}
