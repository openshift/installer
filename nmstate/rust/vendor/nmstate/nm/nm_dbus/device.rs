use std::convert::TryFrom;

use log::warn;

use super::{
    connection::DbusDictionary,
    dbus::{NM_DBUS_INTERFACE_DEV, NM_DBUS_INTERFACE_ROOT},
    lldp::NmLldpNeighbor,
    ErrorKind, NmError,
};

const NM_DEVICE_TYPE_UNKNOWN: u32 = 0;
const NM_DEVICE_TYPE_ETHERNET: u32 = 1;
const NM_DEVICE_TYPE_WIFI: u32 = 2;
// const NM_DEVICE_TYPE_UNUSED1: u32 = 3;
// const NM_DEVICE_TYPE_UNUSED2: u32 = 4;
const NM_DEVICE_TYPE_BT: u32 = 5;
const NM_DEVICE_TYPE_OLPC_MESH: u32 = 6;
const NM_DEVICE_TYPE_WIMAX: u32 = 7;
const NM_DEVICE_TYPE_MODEM: u32 = 8;
const NM_DEVICE_TYPE_INFINIBAND: u32 = 9;
const NM_DEVICE_TYPE_BOND: u32 = 10;
const NM_DEVICE_TYPE_VLAN: u32 = 11;
const NM_DEVICE_TYPE_ADSL: u32 = 12;
const NM_DEVICE_TYPE_BRIDGE: u32 = 13;
const NM_DEVICE_TYPE_GENERIC: u32 = 14;
const NM_DEVICE_TYPE_TEAM: u32 = 15;
const NM_DEVICE_TYPE_TUN: u32 = 16;
const NM_DEVICE_TYPE_IP_TUNNEL: u32 = 17;
const NM_DEVICE_TYPE_MACVLAN: u32 = 18;
const NM_DEVICE_TYPE_VXLAN: u32 = 19;
const NM_DEVICE_TYPE_VETH: u32 = 20;
const NM_DEVICE_TYPE_MACSEC: u32 = 21;
const NM_DEVICE_TYPE_DUMMY: u32 = 22;
const NM_DEVICE_TYPE_PPP: u32 = 23;
const NM_DEVICE_TYPE_OVS_INTERFACE: u32 = 24;
const NM_DEVICE_TYPE_OVS_PORT: u32 = 25;
const NM_DEVICE_TYPE_OVS_BRIDGE: u32 = 26;
const NM_DEVICE_TYPE_WPAN: u32 = 27;
const NM_DEVICE_TYPE_6LOWPAN: u32 = 28;
const NM_DEVICE_TYPE_WIREGUARD: u32 = 29;
const NM_DEVICE_TYPE_WIFI_P2P: u32 = 30;
const NM_DEVICE_TYPE_VRF: u32 = 31;

const NM_DEVICE_STATE_UNKNOWN: u32 = 0;
const NM_DEVICE_STATE_UNMANAGED: u32 = 10;
const NM_DEVICE_STATE_UNAVAILABLE: u32 = 20;
const NM_DEVICE_STATE_DISCONNECTED: u32 = 30;
const NM_DEVICE_STATE_PREPARE: u32 = 40;
const NM_DEVICE_STATE_CONFIG: u32 = 50;
const NM_DEVICE_STATE_NEED_AUTH: u32 = 60;
const NM_DEVICE_STATE_IP_CONFIG: u32 = 70;
const NM_DEVICE_STATE_IP_CHECK: u32 = 80;
const NM_DEVICE_STATE_SECONDARIES: u32 = 90;
const NM_DEVICE_STATE_ACTIVATED: u32 = 100;
const NM_DEVICE_STATE_DEACTIVATING: u32 = 110;
const NM_DEVICE_STATE_FAILED: u32 = 120;

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
pub enum NmDeviceState {
    Unknown,
    Unmanaged,
    Unavailable,
    Disconnected,
    Prepare,
    Config,
    NeedAuth,
    IpConfig,
    IpCheck,
    Secondaries,
    Activated,
    Deactivating,
    Failed,
}

impl Default for NmDeviceState {
    fn default() -> Self {
        Self::Unknown
    }
}

impl From<u32> for NmDeviceState {
    fn from(i: u32) -> Self {
        match i {
            NM_DEVICE_STATE_UNKNOWN => Self::Unknown,
            NM_DEVICE_STATE_UNMANAGED => Self::Unmanaged,
            NM_DEVICE_STATE_UNAVAILABLE => Self::Unavailable,
            NM_DEVICE_STATE_DISCONNECTED => Self::Disconnected,
            NM_DEVICE_STATE_PREPARE => Self::Prepare,
            NM_DEVICE_STATE_CONFIG => Self::Config,
            NM_DEVICE_STATE_NEED_AUTH => Self::NeedAuth,
            NM_DEVICE_STATE_IP_CONFIG => Self::IpConfig,
            NM_DEVICE_STATE_IP_CHECK => Self::IpCheck,
            NM_DEVICE_STATE_SECONDARIES => Self::Secondaries,
            NM_DEVICE_STATE_ACTIVATED => Self::Activated,
            NM_DEVICE_STATE_DEACTIVATING => Self::Deactivating,
            NM_DEVICE_STATE_FAILED => Self::Failed,
            _ => {
                warn!("Unknown Device state reason {}", i);
                Self::Unknown
            }
        }
    }
}

const NM_DEVICE_STATE_REASON_NONE: u32 = 0;
const NM_DEVICE_STATE_REASON_UNKNOWN: u32 = 1;
const NM_DEVICE_STATE_REASON_NOW_MANAGED: u32 = 2;
const NM_DEVICE_STATE_REASON_NOW_UNMANAGED: u32 = 3;
const NM_DEVICE_STATE_REASON_CONFIG_FAILED: u32 = 4;
const NM_DEVICE_STATE_REASON_IP_CONFIG_UNAVAILABLE: u32 = 5;
const NM_DEVICE_STATE_REASON_IP_CONFIG_EXPIRED: u32 = 6;
const NM_DEVICE_STATE_REASON_NO_SECRETS: u32 = 7;
const NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT: u32 = 8;
const NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED: u32 = 9;
const NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED: u32 = 10;
const NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT: u32 = 11;
const NM_DEVICE_STATE_REASON_PPP_START_FAILED: u32 = 12;
const NM_DEVICE_STATE_REASON_PPP_DISCONNECT: u32 = 13;
const NM_DEVICE_STATE_REASON_PPP_FAILED: u32 = 14;
const NM_DEVICE_STATE_REASON_DHCP_START_FAILED: u32 = 15;
const NM_DEVICE_STATE_REASON_DHCP_ERROR: u32 = 16;
const NM_DEVICE_STATE_REASON_DHCP_FAILED: u32 = 17;
const NM_DEVICE_STATE_REASON_SHARED_START_FAILED: u32 = 18;
const NM_DEVICE_STATE_REASON_SHARED_FAILED: u32 = 19;
const NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED: u32 = 20;
const NM_DEVICE_STATE_REASON_AUTOIP_ERROR: u32 = 21;
const NM_DEVICE_STATE_REASON_AUTOIP_FAILED: u32 = 22;
const NM_DEVICE_STATE_REASON_MODEM_BUSY: u32 = 23;
const NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE: u32 = 24;
const NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER: u32 = 25;
const NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT: u32 = 26;
const NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED: u32 = 27;
const NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED: u32 = 28;
const NM_DEVICE_STATE_REASON_GSM_APN_FAILED: u32 = 29;
const NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING: u32 = 30;
const NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED: u32 = 31;
const NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT: u32 = 32;
const NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED: u32 = 33;
const NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED: u32 = 34;
const NM_DEVICE_STATE_REASON_FIRMWARE_MISSING: u32 = 35;
const NM_DEVICE_STATE_REASON_REMOVED: u32 = 36;
const NM_DEVICE_STATE_REASON_SLEEPING: u32 = 37;
const NM_DEVICE_STATE_REASON_CONNECTION_REMOVED: u32 = 38;
const NM_DEVICE_STATE_REASON_USER_REQUESTED: u32 = 39;
const NM_DEVICE_STATE_REASON_CARRIER: u32 = 40;
const NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED: u32 = 41;
const NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE: u32 = 42;
const NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND: u32 = 43;
const NM_DEVICE_STATE_REASON_BT_FAILED: u32 = 44;
const NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED: u32 = 45;
const NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED: u32 = 46;
const NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED: u32 = 47;
const NM_DEVICE_STATE_REASON_GSM_SIM_WRONG: u32 = 48;
const NM_DEVICE_STATE_REASON_INFINIBAND_MODE: u32 = 49;
const NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED: u32 = 50;
const NM_DEVICE_STATE_REASON_BR2684_FAILED: u32 = 51;
const NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE: u32 = 52;
const NM_DEVICE_STATE_REASON_SSID_NOT_FOUND: u32 = 53;
const NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED: u32 = 54;
const NM_DEVICE_STATE_REASON_DCB_FCOE_FAILED: u32 = 55;
const NM_DEVICE_STATE_REASON_TEAMD_CONTROL_FAILED: u32 = 56;
const NM_DEVICE_STATE_REASON_MODEM_FAILED: u32 = 57;
const NM_DEVICE_STATE_REASON_MODEM_AVAILABLE: u32 = 58;
const NM_DEVICE_STATE_REASON_SIM_PIN_INCORRECT: u32 = 59;
const NM_DEVICE_STATE_REASON_NEW_ACTIVATION: u32 = 60;
const NM_DEVICE_STATE_REASON_PARENT_CHANGED: u32 = 61;
const NM_DEVICE_STATE_REASON_PARENT_MANAGED_CHANGED: u32 = 62;
const NM_DEVICE_STATE_REASON_OVSDB_FAILED: u32 = 63;
const NM_DEVICE_STATE_REASON_IP_ADDRESS_DUPLICATE: u32 = 64;
const NM_DEVICE_STATE_REASON_IP_METHOD_UNSUPPORTED: u32 = 65;
const NM_DEVICE_STATE_REASON_SRIOV_CONFIGURATION_FAILED: u32 = 66;
const NM_DEVICE_STATE_REASON_PEER_NOT_FOUND: u32 = 67;

#[derive(Debug, Clone, PartialEq, Eq)]
#[non_exhaustive]
pub enum NmDeviceStateReason {
    Null,
    Unknown,
    NowManaged,
    NowUnmanaged,
    ConfigFailed,
    IpConfigUnavailable,
    IpConfigExpired,
    NoSecrets,
    SupplicantDisconnect,
    SupplicantConfigFailed,
    SupplicantFailed,
    SupplicantTimeout,
    PppStartFailed,
    PppDisconnect,
    PppFailed,
    DhcpStartFailed,
    DhcpError,
    DhcpFailed,
    SharedStartFailed,
    SharedFailed,
    AutoipStartFailed,
    AutoipError,
    AutoipFailed,
    ModemBusy,
    ModemNoDialTone,
    ModemNoCarrier,
    ModemDialTimeout,
    ModemDialFailed,
    ModemInitFailed,
    GsmApnFailed,
    GsmRegistrationNotSearching,
    GsmRegistrationDenied,
    GsmRegistrationTimeout,
    GsmRegistrationFailed,
    GsmPinCheckFailed,
    FirmwareMissing,
    Removed,
    Sleeping,
    ConnectionRemoved,
    UserRequested,
    Carrier,
    ConnectionAssumed,
    SupplicantAvailable,
    ModemNotFound,
    BtFailed,
    GsmSimNotInserted,
    GsmSimPinRequired,
    GsmSimPukRequired,
    GsmSimWrong,
    InfinibandMode,
    DependencyFailed,
    Br2684Failed,
    ModemManagerUnavailable,
    SsidNotFound,
    SecondaryConnectionFailed,
    DcbFcoeFailed,
    TeamdControlFailed,
    ModemFailed,
    ModemAvailable,
    SimPinIncorrect,
    NewActivation,
    ParentChanged,
    ParentManagedChanged,
    OvsdbFailed,
    IpAddressDuplicate,
    IpMethodUnsupported,
    SriovConfigurationFailed,
    PeerNotFound,
}

impl Default for NmDeviceStateReason {
    fn default() -> Self {
        Self::Unknown
    }
}

impl From<u32> for NmDeviceStateReason {
    fn from(i: u32) -> Self {
        match i {
            NM_DEVICE_STATE_REASON_NONE => Self::Null,
            NM_DEVICE_STATE_REASON_UNKNOWN => Self::Unknown,
            NM_DEVICE_STATE_REASON_NOW_MANAGED => Self::NowManaged,
            NM_DEVICE_STATE_REASON_NOW_UNMANAGED => Self::NowUnmanaged,
            NM_DEVICE_STATE_REASON_CONFIG_FAILED => Self::ConfigFailed,
            NM_DEVICE_STATE_REASON_IP_CONFIG_UNAVAILABLE => {
                Self::IpConfigUnavailable
            }
            NM_DEVICE_STATE_REASON_IP_CONFIG_EXPIRED => Self::IpConfigExpired,
            NM_DEVICE_STATE_REASON_NO_SECRETS => Self::NoSecrets,
            NM_DEVICE_STATE_REASON_SUPPLICANT_DISCONNECT => {
                Self::SupplicantDisconnect
            }
            NM_DEVICE_STATE_REASON_SUPPLICANT_CONFIG_FAILED => {
                Self::SupplicantConfigFailed
            }
            NM_DEVICE_STATE_REASON_SUPPLICANT_FAILED => Self::SupplicantFailed,
            NM_DEVICE_STATE_REASON_SUPPLICANT_TIMEOUT => {
                Self::SupplicantTimeout
            }
            NM_DEVICE_STATE_REASON_PPP_START_FAILED => Self::PppStartFailed,
            NM_DEVICE_STATE_REASON_PPP_DISCONNECT => Self::PppDisconnect,
            NM_DEVICE_STATE_REASON_PPP_FAILED => Self::PppFailed,
            NM_DEVICE_STATE_REASON_DHCP_START_FAILED => Self::DhcpStartFailed,
            NM_DEVICE_STATE_REASON_DHCP_ERROR => Self::DhcpError,
            NM_DEVICE_STATE_REASON_DHCP_FAILED => Self::DhcpFailed,
            NM_DEVICE_STATE_REASON_SHARED_START_FAILED => {
                Self::SharedStartFailed
            }
            NM_DEVICE_STATE_REASON_SHARED_FAILED => Self::SharedFailed,
            NM_DEVICE_STATE_REASON_AUTOIP_START_FAILED => {
                Self::AutoipStartFailed
            }
            NM_DEVICE_STATE_REASON_AUTOIP_ERROR => Self::AutoipError,
            NM_DEVICE_STATE_REASON_AUTOIP_FAILED => Self::AutoipFailed,
            NM_DEVICE_STATE_REASON_MODEM_BUSY => Self::ModemBusy,
            NM_DEVICE_STATE_REASON_MODEM_NO_DIAL_TONE => Self::ModemNoDialTone,
            NM_DEVICE_STATE_REASON_MODEM_NO_CARRIER => Self::ModemNoCarrier,
            NM_DEVICE_STATE_REASON_MODEM_DIAL_TIMEOUT => Self::ModemDialTimeout,
            NM_DEVICE_STATE_REASON_MODEM_DIAL_FAILED => Self::ModemDialFailed,
            NM_DEVICE_STATE_REASON_MODEM_INIT_FAILED => Self::ModemInitFailed,
            NM_DEVICE_STATE_REASON_GSM_APN_FAILED => Self::GsmApnFailed,
            NM_DEVICE_STATE_REASON_GSM_REGISTRATION_NOT_SEARCHING => {
                Self::GsmRegistrationNotSearching
            }
            NM_DEVICE_STATE_REASON_GSM_REGISTRATION_DENIED => {
                Self::GsmRegistrationDenied
            }
            NM_DEVICE_STATE_REASON_GSM_REGISTRATION_TIMEOUT => {
                Self::GsmRegistrationTimeout
            }
            NM_DEVICE_STATE_REASON_GSM_REGISTRATION_FAILED => {
                Self::GsmRegistrationFailed
            }
            NM_DEVICE_STATE_REASON_GSM_PIN_CHECK_FAILED => {
                Self::GsmPinCheckFailed
            }
            NM_DEVICE_STATE_REASON_FIRMWARE_MISSING => Self::FirmwareMissing,
            NM_DEVICE_STATE_REASON_REMOVED => Self::Removed,
            NM_DEVICE_STATE_REASON_SLEEPING => Self::Sleeping,
            NM_DEVICE_STATE_REASON_CONNECTION_REMOVED => {
                Self::ConnectionRemoved
            }
            NM_DEVICE_STATE_REASON_USER_REQUESTED => Self::UserRequested,
            NM_DEVICE_STATE_REASON_CARRIER => Self::Carrier,
            NM_DEVICE_STATE_REASON_CONNECTION_ASSUMED => {
                Self::ConnectionAssumed
            }
            NM_DEVICE_STATE_REASON_SUPPLICANT_AVAILABLE => {
                Self::SupplicantAvailable
            }
            NM_DEVICE_STATE_REASON_MODEM_NOT_FOUND => Self::ModemNotFound,
            NM_DEVICE_STATE_REASON_BT_FAILED => Self::BtFailed,
            NM_DEVICE_STATE_REASON_GSM_SIM_NOT_INSERTED => {
                Self::GsmSimNotInserted
            }
            NM_DEVICE_STATE_REASON_GSM_SIM_PIN_REQUIRED => {
                Self::GsmSimPinRequired
            }
            NM_DEVICE_STATE_REASON_GSM_SIM_PUK_REQUIRED => {
                Self::GsmSimPukRequired
            }
            NM_DEVICE_STATE_REASON_GSM_SIM_WRONG => Self::GsmSimWrong,
            NM_DEVICE_STATE_REASON_INFINIBAND_MODE => Self::InfinibandMode,
            NM_DEVICE_STATE_REASON_DEPENDENCY_FAILED => Self::DependencyFailed,
            NM_DEVICE_STATE_REASON_BR2684_FAILED => Self::Br2684Failed,
            NM_DEVICE_STATE_REASON_MODEM_MANAGER_UNAVAILABLE => {
                Self::ModemManagerUnavailable
            }
            NM_DEVICE_STATE_REASON_SSID_NOT_FOUND => Self::SsidNotFound,
            NM_DEVICE_STATE_REASON_SECONDARY_CONNECTION_FAILED => {
                Self::SecondaryConnectionFailed
            }
            NM_DEVICE_STATE_REASON_DCB_FCOE_FAILED => Self::DcbFcoeFailed,
            NM_DEVICE_STATE_REASON_TEAMD_CONTROL_FAILED => {
                Self::TeamdControlFailed
            }
            NM_DEVICE_STATE_REASON_MODEM_FAILED => Self::ModemFailed,
            NM_DEVICE_STATE_REASON_MODEM_AVAILABLE => Self::ModemAvailable,
            NM_DEVICE_STATE_REASON_SIM_PIN_INCORRECT => Self::SimPinIncorrect,
            NM_DEVICE_STATE_REASON_NEW_ACTIVATION => Self::NewActivation,
            NM_DEVICE_STATE_REASON_PARENT_CHANGED => Self::ParentChanged,
            NM_DEVICE_STATE_REASON_PARENT_MANAGED_CHANGED => {
                Self::ParentManagedChanged
            }
            NM_DEVICE_STATE_REASON_OVSDB_FAILED => Self::OvsdbFailed,
            NM_DEVICE_STATE_REASON_IP_ADDRESS_DUPLICATE => {
                Self::IpAddressDuplicate
            }
            NM_DEVICE_STATE_REASON_IP_METHOD_UNSUPPORTED => {
                Self::IpMethodUnsupported
            }
            NM_DEVICE_STATE_REASON_SRIOV_CONFIGURATION_FAILED => {
                Self::SriovConfigurationFailed
            }
            NM_DEVICE_STATE_REASON_PEER_NOT_FOUND => Self::PeerNotFound,
            _ => {
                warn!("Unknown Device state reason {}", i);
                Self::Unknown
            }
        }
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default)]
pub struct NmDevice {
    pub name: String,
    pub iface_type: String,
    pub state: NmDeviceState,
    pub state_reason: NmDeviceStateReason,
    pub is_mac_vtap: bool,
    pub obj_path: String,
    pub real: bool,
}

fn nm_dev_name_get(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<String, NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.get_property::<String>("Interface") {
        Ok(n) => Ok(n),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!(
                "Failed to retrieve interface name of device {}: {}",
                obj_path, e
            ),
        )),
    }
}

fn nm_dev_iface_type_get(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<String, NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.get_property::<u32>("DeviceType") {
        Ok(i) => Ok(match i {
            // Using the NM_SETTING_*_NAME string
            NM_DEVICE_TYPE_UNKNOWN => "unknown".to_string(),
            NM_DEVICE_TYPE_ETHERNET => "802-3-ethernet".to_string(),
            NM_DEVICE_TYPE_WIFI => "802-11-wireless".to_string(),
            NM_DEVICE_TYPE_BT => "bluetooth".to_string(),
            NM_DEVICE_TYPE_OLPC_MESH => "802-11-olpc-mesh".to_string(),
            NM_DEVICE_TYPE_WIMAX => "wimax".to_string(),
            NM_DEVICE_TYPE_MODEM => "modem".to_string(),
            NM_DEVICE_TYPE_INFINIBAND => "infiniband".to_string(),
            NM_DEVICE_TYPE_BOND => "bond".to_string(),
            NM_DEVICE_TYPE_VLAN => "vlan".to_string(),
            NM_DEVICE_TYPE_ADSL => "adsl".to_string(),
            NM_DEVICE_TYPE_BRIDGE => "bridge".to_string(),
            NM_DEVICE_TYPE_GENERIC => "generic".to_string(),
            NM_DEVICE_TYPE_TEAM => "team".to_string(),
            NM_DEVICE_TYPE_TUN => "tun".to_string(),
            NM_DEVICE_TYPE_IP_TUNNEL => "ip-tunnel".to_string(),
            NM_DEVICE_TYPE_MACVLAN => "macvlan".to_string(),
            NM_DEVICE_TYPE_VXLAN => "vxlan".to_string(),
            NM_DEVICE_TYPE_VETH => "veth".to_string(),
            NM_DEVICE_TYPE_MACSEC => "macsec".to_string(),
            NM_DEVICE_TYPE_DUMMY => "dummy".to_string(),
            NM_DEVICE_TYPE_PPP => "ppp".to_string(),
            NM_DEVICE_TYPE_OVS_INTERFACE => "ovs-interface".to_string(),
            NM_DEVICE_TYPE_OVS_PORT => "ovs-port".to_string(),
            NM_DEVICE_TYPE_OVS_BRIDGE => "ovs-bridge".to_string(),
            NM_DEVICE_TYPE_WPAN => "wpan".to_string(),
            NM_DEVICE_TYPE_6LOWPAN => "6lowpan".to_string(),
            NM_DEVICE_TYPE_WIREGUARD => "wireguard".to_string(),
            NM_DEVICE_TYPE_WIFI_P2P => "wifi-p2p".to_string(),
            NM_DEVICE_TYPE_VRF => "vrf".to_string(),
            _ => format!("unknown({})", i),
        }),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!(
                "Failed to retrieve device type of device {}: {}",
                obj_path, e
            ),
        )),
    }
}

fn nm_dev_state_reason_get(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<(NmDeviceState, NmDeviceStateReason), NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.get_property::<(u32, u32)>("StateReason") {
        Ok((state, state_reason)) => Ok((state.into(), state_reason.into())),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!(
                "Failed to retrieve state reason of device {}: {}",
                obj_path, e
            ),
        )),
    }
}

fn nm_dev_is_mac_vtap_get(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<bool, NmError> {
    let dbus_iface = format!("{}.Macvlan", NM_DBUS_INTERFACE_DEV);
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        &dbus_iface,
    )?;
    match proxy.get_property::<bool>("Tab") {
        Ok(v) => Ok(v),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!(
                "Failed to retrieve Macvlan.Tab(tap) of device {}: {}",
                obj_path, e
            ),
        )),
    }
}

fn nm_dev_real_get(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<bool, NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.get_property::<bool>("Real") {
        Ok(r) => Ok(r),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!("Failed to retrieve real of device {}: {}", obj_path, e),
        )),
    }
}

pub(crate) fn nm_dev_from_obj_path(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<NmDevice, NmError> {
    let real = nm_dev_real_get(dbus_conn, obj_path)?;
    let (state, state_reason) = nm_dev_state_reason_get(dbus_conn, obj_path)?;
    let mut dev = NmDevice {
        name: nm_dev_name_get(dbus_conn, obj_path)?,
        iface_type: nm_dev_iface_type_get(dbus_conn, obj_path)?,
        state,
        state_reason,
        obj_path: obj_path.to_string(),
        is_mac_vtap: false,
        real,
    };
    if dev.iface_type == "macvlan" {
        dev.is_mac_vtap = nm_dev_is_mac_vtap_get(dbus_conn, obj_path)?;
    }
    Ok(dev)
}

pub(crate) fn nm_dev_delete(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<(), NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.call::<(), ()>("Delete", &()) {
        Ok(()) => Ok(()),
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!("Failed to delete device {}: {}", obj_path, e),
        )),
    }
}

pub(crate) fn nm_dev_get_llpd(
    dbus_conn: &zbus::Connection,
    obj_path: &str,
) -> Result<Vec<NmLldpNeighbor>, NmError> {
    let proxy = zbus::Proxy::new(
        dbus_conn,
        NM_DBUS_INTERFACE_ROOT,
        obj_path,
        NM_DBUS_INTERFACE_DEV,
    )?;
    match proxy.get_property::<Vec<DbusDictionary>>("LldpNeighbors") {
        Ok(v) => {
            let mut ret = Vec::new();
            for value in v {
                ret.push(NmLldpNeighbor::try_from(value)?);
            }
            Ok(ret)
        }
        Err(e) => Err(NmError::new(
            ErrorKind::Bug,
            format!(
                "Failed to retrieve LLDP neighbors of device {}: {}",
                obj_path, e
            ),
        )),
    }
}
