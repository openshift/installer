use serde::{ser::SerializeStruct, Deserialize, Serialize, Serializer};

const LLDP_CHASSIS_ID_TYPE: u8 = 1;
const LLDP_PORT_TYPE: u8 = 2;
const LLDP_SYSTEM_NAME_TYPE: u8 = 5;
const LLDP_SYSTEM_DESCRIPTION_TYPE: u8 = 6;
const LLDP_SYSTEM_CAPABILITIES_TYPE: u8 = 7;
const LLDP_MANAGEMENT_ADDRESSES_TYPE: u8 = 8;
const LLDP_ORGANIZATION_SPECIFIC_TYPE: u8 = 127;

const LLDP_CHASSIS_ID_RESERVED: u8 = 0;
const LLDP_CHASSIS_ID_CHASSIS_COMPONENT: u8 = 1;
const LLDP_CHASSIS_ID_IFACE_ALIAS: u8 = 2;
const LLDP_CHASSIS_ID_PORT_COMPONENT: u8 = 3;
const LLDP_CHASSIS_ID_MAC_ADDR: u8 = 4;
const LLDP_CHASSIS_ID_NET_ADDR: u8 = 5;
const LLDP_CHASSIS_ID_IFACE_NAME: u8 = 6;
const LLDP_CHASSIS_ID_LOCAL: u8 = 7;

const LLDP_PORT_ID_RESERVED: u8 = 0;
const LLDP_PORT_ID_IFACE_ALIAS: u8 = 1;
const LLDP_PORT_ID_PORT_COMPONENT: u8 = 2;
const LLDP_PORT_ID_MAC_ADDR: u8 = 3;
const LLDP_PORT_ID_NET_ADDR: u8 = 4;
const LLDP_PORT_ID_IFACE_NAME: u8 = 5;
const LLDP_PORT_ID_AGENT_CIRCUIT_ID: u8 = 6;
const LLDP_PORT_ID_LOCAL: u8 = 7;

const LLDP_ORG_OIU_VLAN: &str = "00:80:c2";
const LLDP_ORG_SUBTYPE_VLAN: u8 = 3;

const LLDP_ORG_OIU_MAC_PHY_CONF: &str = "00:12:0f";
const LLDP_ORG_SUBTYPE_MAC_PHY_CONF: u8 = 1;

const LLDP_ORG_OIU_PPVIDS: &str = "00:80:c2";
const LLDP_ORG_SUBTYPE_PPVIDS: u8 = 2;

const LLDP_ORG_OIU_MAX_FRAME_SIZE: &str = "00:12:0f";
const LLDP_ORG_SUBTYPE_MAX_FRAME_SIZE: u8 = 4;

const LLDP_SYS_CAP_OTHER: u16 = 1;
const LLDP_SYS_CAP_REPEATER: u16 = 2;
const LLDP_SYS_CAP_MAC_BRIDGE: u16 = 3;
const LLDP_SYS_CAP_AP: u16 = 4;
const LLDP_SYS_CAP_ROUTER: u16 = 5;
const LLDP_SYS_CAP_TELEPHONE: u16 = 6;
const LLDP_SYS_CAP_DOCSIS: u16 = 7;
const LLDP_SYS_CAP_STATION_ONLY: u16 = 8;
const LLDP_SYS_CAP_CVLAN: u16 = 9;
const LLDP_SYS_CAP_SVLAN: u16 = 10;
const LLDP_SYS_CAP_TWO_PORT_MAC_RELAY: u16 = 11;

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize, Serialize)]
#[non_exhaustive]
#[serde(deny_unknown_fields)]
pub struct LldpConfig {
    #[serde(deserialize_with = "crate::deserializer::bool_or_string")]
    pub enabled: bool,
    #[serde(skip_deserializing, skip_serializing_if = "Vec::is_empty")]
    pub neighbors: Vec<Vec<LldpNeighborTlv>>,
}

impl LldpConfig {
    pub(crate) fn pre_verify_cleanup(&mut self) {
        self.neighbors = Vec::new();
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[serde(untagged)]
#[non_exhaustive]
pub enum LldpNeighborTlv {
    SystemName(LldpSystemName),
    SystemDescription(LldpSystemDescription),
    SystemCapabilities(LldpSystemCapabilities),
    ChassisId(LldpChassisId),
    PortId(LldpPortId),
    Ieee8021Vlans(LldpVlans),
    Ieee8023MacPhyConf(LldpMacPhyConf),
    Ieee8021Ppvids(LldpPpvids),
    ManagementAddresses(LldpMgmtAddrs),
    Ieee8023MaxFrameSize(LldpMaxFrameSize),
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub struct LldpSystemName(pub String);

impl Serialize for LldpSystemName {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_system_name", 2)?;
        serial_struct.serialize_field("type", &LLDP_SYSTEM_NAME_TYPE)?;
        serial_struct.serialize_field("system-name", &self.0)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub struct LldpSystemDescription(pub String);

impl Serialize for LldpSystemDescription {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_system_description", 2)?;
        serial_struct.serialize_field("type", &LLDP_SYSTEM_DESCRIPTION_TYPE)?;
        serial_struct.serialize_field("system-description", &self.0)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpChassisId {
    pub id: String,
    pub id_type: LldpChassisIdType,
}

impl Serialize for LldpChassisId {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_chassis_id", 4)?;
        serial_struct.serialize_field("_description", &self.id_type)?;
        serial_struct
            .serialize_field("chassis-id-type", &(self.id_type as u8))?;
        serial_struct.serialize_field("type", &LLDP_CHASSIS_ID_TYPE)?;
        serial_struct.serialize_field("chassis-id", &self.id)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum LldpChassisIdType {
    Reserved,
    #[serde(rename = "Chassis component")]
    ChassisComponent,
    #[serde(rename = "Interface alias")]
    InterfaceAlias,
    #[serde(rename = "Port component")]
    PortComponent,
    #[serde(rename = "MAC address")]
    MacAddress,
    #[serde(rename = "Network address")]
    NetworkAddress,
    #[serde(rename = "Interface name")]
    InterfaceName,
    #[serde(rename = "Locally assigned")]
    LocallyAssigned,
}

impl From<LldpChassisIdType> for u8 {
    fn from(v: LldpChassisIdType) -> u8 {
        match v {
            LldpChassisIdType::Reserved => LLDP_CHASSIS_ID_RESERVED,
            LldpChassisIdType::ChassisComponent => {
                LLDP_CHASSIS_ID_CHASSIS_COMPONENT
            }
            LldpChassisIdType::InterfaceAlias => LLDP_CHASSIS_ID_IFACE_ALIAS,
            LldpChassisIdType::PortComponent => LLDP_CHASSIS_ID_PORT_COMPONENT,
            LldpChassisIdType::MacAddress => LLDP_CHASSIS_ID_MAC_ADDR,
            LldpChassisIdType::NetworkAddress => LLDP_CHASSIS_ID_NET_ADDR,
            LldpChassisIdType::InterfaceName => LLDP_CHASSIS_ID_IFACE_NAME,
            LldpChassisIdType::LocallyAssigned => LLDP_CHASSIS_ID_LOCAL,
        }
    }
}

impl From<u8> for LldpChassisIdType {
    fn from(i: u8) -> Self {
        match i {
            LLDP_CHASSIS_ID_CHASSIS_COMPONENT => Self::ChassisComponent,
            LLDP_CHASSIS_ID_IFACE_ALIAS => Self::InterfaceAlias,
            LLDP_CHASSIS_ID_PORT_COMPONENT => Self::PortComponent,
            LLDP_CHASSIS_ID_MAC_ADDR => Self::MacAddress,
            LLDP_CHASSIS_ID_NET_ADDR => Self::NetworkAddress,
            LLDP_CHASSIS_ID_IFACE_NAME => Self::InterfaceName,
            LLDP_CHASSIS_ID_LOCAL => Self::LocallyAssigned,
            _ => Self::Reserved,
        }
    }
}

impl Default for LldpChassisIdType {
    fn default() -> Self {
        Self::Reserved
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpSystemCapabilities(pub u16);

impl Serialize for LldpSystemCapabilities {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_system_capabilities", 2)?;
        serial_struct
            .serialize_field("type", &LLDP_SYSTEM_CAPABILITIES_TYPE)?;
        serial_struct
            .serialize_field("system-capabilities", &parse_sys_caps(self.0))?;
        serial_struct.end()
    }
}

fn parse_sys_caps(caps: u16) -> Vec<LldpSystemCapability> {
    let mut ret = Vec::new();
    if (caps & 1 << (LLDP_SYS_CAP_OTHER - 1)) > 0 {
        ret.push(LldpSystemCapability::Other);
    }
    if (caps & 1 << (LLDP_SYS_CAP_REPEATER - 1)) > 0 {
        ret.push(LldpSystemCapability::Repeater);
    }
    if (caps & 1 << (LLDP_SYS_CAP_MAC_BRIDGE - 1)) > 0 {
        ret.push(LldpSystemCapability::MacBridgeComponent);
    }
    if (caps & 1 << (LLDP_SYS_CAP_AP - 1)) > 0 {
        ret.push(LldpSystemCapability::AccessPoint);
    }
    if (caps & 1 << (LLDP_SYS_CAP_ROUTER - 1)) > 0 {
        ret.push(LldpSystemCapability::Router);
    }
    if (caps & 1 << (LLDP_SYS_CAP_TELEPHONE - 1)) > 0 {
        ret.push(LldpSystemCapability::Telephone);
    }
    if (caps & 1 << (LLDP_SYS_CAP_DOCSIS - 1)) > 0 {
        ret.push(LldpSystemCapability::DocsisCableDevice);
    }
    if (caps & 1 << (LLDP_SYS_CAP_STATION_ONLY - 1)) > 0 {
        ret.push(LldpSystemCapability::StationOnly);
    }
    if (caps & 1 << (LLDP_SYS_CAP_CVLAN - 1)) > 0 {
        ret.push(LldpSystemCapability::CVlanComponent);
    }
    if (caps & 1 << (LLDP_SYS_CAP_SVLAN - 1)) > 0 {
        ret.push(LldpSystemCapability::SVlanComponent);
    }
    if (caps & 1 << (LLDP_SYS_CAP_TWO_PORT_MAC_RELAY - 1)) > 0 {
        ret.push(LldpSystemCapability::TwoPortMacRelayComponent);
    }
    ret
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[non_exhaustive]
pub enum LldpSystemCapability {
    Other,
    Repeater,
    #[serde(rename = "MAC Bridge component")]
    MacBridgeComponent,
    #[serde(rename = "802.11 Access Point (AP)")]
    AccessPoint,
    Router,
    Telephone,
    #[serde(rename = "DOCSIS cable device")]
    DocsisCableDevice,
    #[serde(rename = "Station Only")]
    StationOnly,
    #[serde(rename = "C-VLAN component")]
    CVlanComponent,
    #[serde(rename = "S-VLAN component")]
    SVlanComponent,
    #[serde(rename = "Two-port MAC Relay component")]
    TwoPortMacRelayComponent,
}

impl Default for LldpSystemCapability {
    fn default() -> Self {
        Self::Other
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpPortId {
    pub id: String,
    pub id_type: LldpPortIdType,
}

impl Serialize for LldpPortId {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_port_id", 4)?;
        serial_struct.serialize_field("_description", &self.id_type)?;
        serial_struct.serialize_field("port-id-type", &(self.id_type as u8))?;
        serial_struct.serialize_field("type", &LLDP_PORT_TYPE)?;
        serial_struct.serialize_field("port-id", &self.id)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize)]
pub enum LldpPortIdType {
    Reserved,
    #[serde(rename = "Interface alias")]
    InterfaceAlias,
    #[serde(rename = "Port component")]
    PortComponent,
    #[serde(rename = "MAC address")]
    MacAddress,
    #[serde(rename = "Network address")]
    NetworkAddress,
    #[serde(rename = "Interface name")]
    InterfaceName,
    #[serde(rename = "Agent circuit ID")]
    AgentCircuitId,
    #[serde(rename = "Locally assigned")]
    LocallyAssigned,
}

impl Default for LldpPortIdType {
    fn default() -> Self {
        Self::Reserved
    }
}

impl From<LldpPortIdType> for u8 {
    fn from(v: LldpPortIdType) -> u8 {
        match v {
            LldpPortIdType::Reserved => LLDP_PORT_ID_RESERVED,
            LldpPortIdType::InterfaceAlias => LLDP_PORT_ID_IFACE_ALIAS,
            LldpPortIdType::PortComponent => LLDP_PORT_ID_PORT_COMPONENT,
            LldpPortIdType::MacAddress => LLDP_PORT_ID_MAC_ADDR,
            LldpPortIdType::NetworkAddress => LLDP_PORT_ID_NET_ADDR,
            LldpPortIdType::InterfaceName => LLDP_PORT_ID_IFACE_NAME,
            LldpPortIdType::AgentCircuitId => LLDP_PORT_ID_AGENT_CIRCUIT_ID,
            LldpPortIdType::LocallyAssigned => LLDP_PORT_ID_LOCAL,
        }
    }
}

impl From<u8> for LldpPortIdType {
    fn from(v: u8) -> Self {
        match v {
            LLDP_PORT_ID_IFACE_ALIAS => Self::InterfaceAlias,
            LLDP_PORT_ID_PORT_COMPONENT => Self::PortComponent,
            LLDP_PORT_ID_MAC_ADDR => Self::MacAddress,
            LLDP_PORT_ID_NET_ADDR => Self::NetworkAddress,
            LLDP_PORT_ID_IFACE_NAME => Self::InterfaceName,
            LLDP_PORT_ID_AGENT_CIRCUIT_ID => Self::AgentCircuitId,
            LLDP_PORT_ID_LOCAL => Self::LocallyAssigned,
            _ => Self::Reserved,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpVlans(pub Vec<LldpVlan>);

impl Serialize for LldpVlans {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct = serializer.serialize_struct("lldp_vlans", 4)?;
        serial_struct
            .serialize_field("type", &LLDP_ORGANIZATION_SPECIFIC_TYPE)?;
        serial_struct.serialize_field("ieee-802-1-vlans", &self.0)?;
        serial_struct.serialize_field("oui", LLDP_ORG_OIU_VLAN)?;
        serial_struct.serialize_field("subtype", &LLDP_ORG_SUBTYPE_VLAN)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Default)]
#[non_exhaustive]
pub struct LldpVlan {
    pub name: String,
    pub vid: u32,
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpMacPhyConf {
    pub autoneg: bool,
    pub operational_mau_type: u16,
    pub pmd_autoneg_cap: u16,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[non_exhaustive]
#[serde(rename_all = "kebab-case")]
struct _LldpMacPhyConf {
    autoneg: bool,
    operational_mau_type: u16,
    pmd_autoneg_cap: u16,
}

impl From<&LldpMacPhyConf> for _LldpMacPhyConf {
    fn from(conf: &LldpMacPhyConf) -> Self {
        Self {
            autoneg: conf.autoneg,
            operational_mau_type: conf.operational_mau_type,
            pmd_autoneg_cap: conf.pmd_autoneg_cap,
        }
    }
}

impl Serialize for LldpMacPhyConf {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_mac_phy_conf", 4)?;
        serial_struct
            .serialize_field("type", &LLDP_ORGANIZATION_SPECIFIC_TYPE)?;
        serial_struct.serialize_field(
            "ieee-802-3-mac-phy-conf",
            &_LldpMacPhyConf::from(self),
        )?;
        serial_struct.serialize_field("oui", LLDP_ORG_OIU_MAC_PHY_CONF)?;
        serial_struct
            .serialize_field("subtype", &LLDP_ORG_SUBTYPE_MAC_PHY_CONF)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpPpvids(pub Vec<u32>);
impl Serialize for LldpPpvids {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_ppvids", 4)?;
        serial_struct
            .serialize_field("type", &LLDP_ORGANIZATION_SPECIFIC_TYPE)?;
        serial_struct.serialize_field("ieee-802-1-ppvids", &self.0)?;
        serial_struct.serialize_field("oui", LLDP_ORG_OIU_PPVIDS)?;
        serial_struct.serialize_field("subtype", &LLDP_ORG_SUBTYPE_PPVIDS)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpMgmtAddrs(pub Vec<LldpMgmtAddr>);
impl Serialize for LldpMgmtAddrs {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_mgmt_addrs", 2)?;
        serial_struct
            .serialize_field("type", &LLDP_MANAGEMENT_ADDRESSES_TYPE)?;
        serial_struct.serialize_field("management-addresses", &self.0)?;
        serial_struct.end()
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Default)]
#[non_exhaustive]
#[serde(rename_all = "kebab-case")]
pub struct LldpMgmtAddr {
    pub address: String,
    pub address_subtype: LldpAddressFamily,
    pub interface_number: u32,
    pub interface_number_subtype: u32,
}

#[derive(Debug, Clone, PartialEq, Eq, Serialize)]
#[non_exhaustive]
pub enum LldpAddressFamily {
    Unknown,
    #[serde(rename = "IPv4")]
    Ipv4,
    #[serde(rename = "IPv6")]
    Ipv6,
    #[serde(rename = "MAC")]
    Mac,
}

impl Default for LldpAddressFamily {
    fn default() -> Self {
        Self::Unknown
    }
}

const ADDRESS_FAMILY_IP4: u16 = 1;
const ADDRESS_FAMILY_IP6: u16 = 2;
const ADDRESS_FAMILY_MAC: u16 = 6;

impl From<u16> for LldpAddressFamily {
    fn from(i: u16) -> Self {
        match i {
            ADDRESS_FAMILY_IP4 => Self::Ipv4,
            ADDRESS_FAMILY_IP6 => Self::Ipv6,
            ADDRESS_FAMILY_MAC => Self::Mac,
            _ => Self::Unknown,
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
#[non_exhaustive]
pub struct LldpMaxFrameSize(pub u32);
impl Serialize for LldpMaxFrameSize {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut serial_struct =
            serializer.serialize_struct("lldp_max_frame_size", 4)?;
        serial_struct
            .serialize_field("type", &LLDP_ORGANIZATION_SPECIFIC_TYPE)?;
        serial_struct.serialize_field("ieee-802-3-max-frame-size", &self.0)?;
        serial_struct.serialize_field("oui", LLDP_ORG_OIU_MAX_FRAME_SIZE)?;
        serial_struct
            .serialize_field("subtype", &LLDP_ORG_SUBTYPE_MAX_FRAME_SIZE)?;
        serial_struct.end()
    }
}
