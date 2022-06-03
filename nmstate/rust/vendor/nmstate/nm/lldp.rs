use std::convert::TryFrom;
use std::fmt::Write;

use crate::nm::nm_dbus::{
    NmConnection, NmLldpNeighbor, NmLldpNeighbor8021Vlan,
};

use crate::{
    LldpAddressFamily, LldpChassisId, LldpConfig, LldpMacPhyConf,
    LldpMaxFrameSize, LldpMgmtAddr, LldpMgmtAddrs, LldpNeighborTlv, LldpPortId,
    LldpPpvids, LldpSystemCapabilities, LldpSystemDescription, LldpSystemName,
    LldpVlan, LldpVlans,
};

pub(crate) fn is_lldp_enabled(nm_conn: &NmConnection) -> bool {
    nm_conn.connection.as_ref().and_then(|s| s.lldp) == Some(true)
}

pub(crate) fn get_lldp(nm_infos: Vec<NmLldpNeighbor>) -> LldpConfig {
    let mut neighbors = Vec::new();
    for nm_info in nm_infos {
        let tlvs = nm_neighbor_to_nmstate(&nm_info);
        if !tlvs.is_empty() {
            neighbors.push(tlvs);
        }
    }
    LldpConfig {
        enabled: true,
        neighbors,
    }
}

fn nm_neighbor_to_nmstate(nm_info: &NmLldpNeighbor) -> Vec<LldpNeighborTlv> {
    let mut ret = Vec::new();

    if let Some(c) = get_sys_name(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_sys_description(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_sys_caps(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_chassis_id(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_port_id(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_vlans(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_mac_phy_conf(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_ppvids(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_mgmt_addrs(nm_info) {
        ret.push(c)
    }
    if let Some(c) = get_max_frame_size(nm_info) {
        ret.push(c)
    }

    ret
}

impl From<&[NmLldpNeighbor8021Vlan]> for LldpVlans {
    fn from(nm_vlans: &[NmLldpNeighbor8021Vlan]) -> Self {
        let mut vlans = Vec::new();
        for nm_vlan in nm_vlans {
            if let (Some(vid), Some(name)) = (nm_vlan.vid, &nm_vlan.name) {
                vlans.push(LldpVlan {
                    vid,
                    name: name.to_string(),
                });
            }
        }

        LldpVlans(vlans)
    }
}

fn u8_addr_to_mac_string(data: &[u8]) -> String {
    let mut addr = String::new();
    for (i, &val) in data.iter().enumerate() {
        let _ = write!(addr, "{:02X}", val);
        if i != data.len() - 1 {
            addr.push(':');
        }
    }
    addr
}

fn get_chassis_id(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let (Some(id_type), Some(id)) =
        (nm_info.chassis_id_type, nm_info.chassis_id.as_deref())
    {
        if let Ok(id_type) = u8::try_from(id_type) {
            let chassis_id = LldpChassisId {
                id: id.to_string(),
                id_type: id_type.into(),
            };
            return Some(LldpNeighborTlv::ChassisId(chassis_id));
        } else {
            log::warn!(
                "Got unsupported chassis_id_type {}, expecting u8",
                id_type
            );
        }
    }
    None
}

fn get_port_id(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let (Some(id_type), Some(id)) =
        (nm_info.port_id_type, nm_info.port_id.as_deref())
    {
        if let Ok(id_type) = u8::try_from(id_type) {
            let port_id = LldpPortId {
                id: id.to_string(),
                id_type: id_type.into(),
            };
            return Some(LldpNeighborTlv::PortId(port_id));
        } else {
            log::warn!(
                "Got unsupported port_id_type {}, expecting u8",
                id_type
            );
        }
    }
    None
}

fn get_sys_caps(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(s) = nm_info.system_capabilities {
        if let Ok(caps) = u16::try_from(s) {
            return Some(LldpNeighborTlv::SystemCapabilities(
                LldpSystemCapabilities(caps),
            ));
        }
    }
    None
}

fn get_sys_name(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(s) = nm_info.system_name.as_deref() {
        return Some(LldpNeighborTlv::SystemName(LldpSystemName(
            s.to_string(),
        )));
    }
    None
}

fn get_sys_description(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(s) = nm_info.system_description.as_deref() {
        return Some(LldpNeighborTlv::SystemDescription(
            LldpSystemDescription(s.to_string()),
        ));
    }
    None
}

fn get_vlans(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(nm_vlans) = nm_info.ieee_802_1_vlans.as_deref() {
        return Some(LldpNeighborTlv::Ieee8021Vlans(nm_vlans.into()));
    }
    None
}

fn get_mac_phy_conf(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(nm_conf) = nm_info.ieee_802_3_mac_phy_conf.as_ref() {
        if let (Some(a), Some(p), Some(o)) = (
            nm_conf.autoneg,
            nm_conf.pmd_autoneg_cap,
            nm_conf.operational_mau_type,
        ) {
            if let (Ok(o), Ok(p)) = (u16::try_from(o), u16::try_from(p)) {
                let conf = LldpMacPhyConf {
                    autoneg: a > 0,
                    operational_mau_type: o,
                    pmd_autoneg_cap: p,
                };
                return Some(LldpNeighborTlv::Ieee8023MacPhyConf(conf));
            }
        }
    }
    None
}

fn get_ppvids(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(nm_ppvids) = nm_info.ieee_802_1_ppvids.as_ref() {
        let mut ppvids = Vec::new();
        for nm_ppvid in nm_ppvids {
            if let Some(p) = nm_ppvid.ppvid {
                ppvids.push(p);
            }
        }
        return Some(LldpNeighborTlv::Ieee8021Ppvids(LldpPpvids(ppvids)));
    }
    None
}

fn get_mgmt_addrs(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(nm_mgmt_addrs) = nm_info.management_addresses.as_deref() {
        let mut addrs = Vec::new();
        for nm_mgmt_addr in nm_mgmt_addrs {
            if let (Some(at), Some(a), Some(it), Some(i)) = (
                nm_mgmt_addr.address_subtype,
                nm_mgmt_addr.address.as_ref(),
                nm_mgmt_addr.interface_number_subtype,
                nm_mgmt_addr.interface_number,
            ) {
                if let Ok(at) = u16::try_from(at) {
                    let mut mgmt_addr = LldpMgmtAddr::default();
                    let at = LldpAddressFamily::from(at);
                    let addr = match at {
                        LldpAddressFamily::Ipv4 => {
                            if a.len() != 4 {
                                continue;
                            }
                            std::net::Ipv4Addr::new(a[0], a[1], a[2], a[3])
                                .to_string()
                        }
                        LldpAddressFamily::Ipv6 => {
                            if a.len() != 16 {
                                continue;
                            }
                            let mut buff = [0u8; 16];
                            buff.copy_from_slice(&a[..16]);
                            std::net::Ipv6Addr::from(buff).to_string()
                        }
                        LldpAddressFamily::Mac => u8_addr_to_mac_string(a),
                        _ => {
                            continue;
                        }
                    };

                    mgmt_addr.address_subtype = at;
                    mgmt_addr.address = addr;
                    mgmt_addr.interface_number_subtype = it;
                    mgmt_addr.interface_number = i;
                    addrs.push(mgmt_addr);
                }
            }
        }
        return Some(LldpNeighborTlv::ManagementAddresses(LldpMgmtAddrs(
            addrs,
        )));
    }
    None
}

fn get_max_frame_size(nm_info: &NmLldpNeighbor) -> Option<LldpNeighborTlv> {
    if let Some(s) = nm_info.ieee_802_3_max_frame_size {
        return Some(LldpNeighborTlv::Ieee8023MaxFrameSize(LldpMaxFrameSize(
            s,
        )));
    }
    None
}
