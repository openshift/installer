use std::convert::TryFrom;

use serde::Deserialize;

use super::{
    connection::{DbusDictionary, _from_map},
    NmError,
};

#[derive(Debug, Clone, PartialEq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
pub struct NmLldpNeighbor {
    pub raw: Option<Vec<u8>>,
    pub chassis_id_type: Option<u32>,
    pub chassis_id: Option<String>,
    pub port_id_type: Option<u32>,
    pub port_id: Option<String>,
    pub destination: Option<String>,
    pub port_description: Option<String>,
    pub system_name: Option<String>,
    pub system_description: Option<String>,
    pub system_capabilities: Option<u32>,
    pub management_addresses: Option<Vec<NmLldpNeighborMgmtAddr>>,
    pub ieee_802_1_pvid: Option<u32>,
    pub ieee_802_1_ppvid: Option<u32>,
    pub ieee_802_1_ppvid_flags: Option<u32>,
    pub ieee_802_1_ppvids: Option<Vec<NmLldpNeighbor8021Ppvid>>,
    pub ieee_802_1_vid: Option<u32>,
    pub ieee_802_1_vlan_name: Option<String>,
    pub ieee_802_1_vlans: Option<Vec<NmLldpNeighbor8021Vlan>>,
    pub ieee_802_3_mac_phy_conf: Option<NmLldpNeighbor8023MacPhyConf>,
    pub ieee_802_3_power_via_mdi: Option<NmLldpNeighbor8023PowerViaMdi>,
    pub ieee_802_3_max_frame_size: Option<u32>,
    _other: DbusDictionary,
}

impl TryFrom<DbusDictionary> for NmLldpNeighbor {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            raw: _from_map!(v, "raw", <Vec<u8>>::try_from)?,
            chassis_id_type: _from_map!(v, "chassis-id-type", <u32>::try_from)?,
            chassis_id: _from_map!(v, "chassis-id", <String>::try_from)?,
            port_id_type: _from_map!(v, "port-id-type", <u32>::try_from)?,
            port_id: _from_map!(v, "port-id", <String>::try_from)?,
            destination: _from_map!(v, "destination", <String>::try_from)?,
            port_description: _from_map!(
                v,
                "port-description",
                <String>::try_from
            )?,
            system_name: _from_map!(v, "system-name", <String>::try_from)?,
            system_description: _from_map!(
                v,
                "system-description",
                <String>::try_from
            )?,
            system_capabilities: _from_map!(
                v,
                "system-capabilities",
                <u32>::try_from
            )?,
            management_addresses: _from_map!(
                v,
                "management-addresses",
                parse_mgmt_addrs
            )?,
            ieee_802_1_pvid: _from_map!(v, "ieee-802-1-pvid", <u32>::try_from)?,
            ieee_802_1_ppvid: _from_map!(
                v,
                "ieee-802-1-ppvid",
                <u32>::try_from
            )?,
            ieee_802_1_ppvid_flags: _from_map!(
                v,
                "ieee-802-1-ppvid-flags",
                <u32>::try_from
            )?,
            ieee_802_1_ppvids: _from_map!(
                v,
                "ieee-802-1-ppvids",
                parse_ppvids
            )?,
            ieee_802_1_vid: _from_map!(v, "ieee-802-1-vid", <u32>::try_from)?,
            ieee_802_1_vlan_name: _from_map!(
                v,
                "ieee-802-1-vlan-name",
                <String>::try_from
            )?,
            ieee_802_1_vlans: _from_map!(v, "ieee-802-1-vlans", parse_vlans)?,
            ieee_802_3_mac_phy_conf: _from_map!(
                v,
                "ieee-802-3-mac-phy-conf",
                NmLldpNeighbor8023MacPhyConf::try_from
            )?,
            ieee_802_3_power_via_mdi: _from_map!(
                v,
                "ieee-802-3-power-via-mdi",
                NmLldpNeighbor8023PowerViaMdi::try_from
            )?,
            ieee_802_3_max_frame_size: _from_map!(
                v,
                "ieee-802-3-max-frame-size",
                <u32>::try_from
            )?,
            _other: v,
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
pub struct NmLldpNeighbor8021Ppvid {
    pub ppvid: Option<u32>,
    pub flags: Option<u32>,
}

fn parse_ppvids(
    v: zvariant::OwnedValue,
) -> Result<Vec<NmLldpNeighbor8021Ppvid>, NmError> {
    let mut ret = Vec::new();
    for value in Vec::<DbusDictionary>::try_from(v)? {
        ret.push(NmLldpNeighbor8021Ppvid::try_from(value)?);
    }
    Ok(ret)
}

impl TryFrom<DbusDictionary> for NmLldpNeighbor8021Ppvid {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            ppvid: _from_map!(v, "ppvid", <u32>::try_from)?,
            flags: _from_map!(v, "flags", <u32>::try_from)?,
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
pub struct NmLldpNeighbor8021Vlan {
    pub vid: Option<u32>,
    pub name: Option<String>,
}

fn parse_vlans(
    v: zvariant::OwnedValue,
) -> Result<Vec<NmLldpNeighbor8021Vlan>, NmError> {
    let mut ret = Vec::new();
    for value in Vec::<DbusDictionary>::try_from(v)? {
        ret.push(NmLldpNeighbor8021Vlan::try_from(value)?);
    }
    Ok(ret)
}

impl TryFrom<DbusDictionary> for NmLldpNeighbor8021Vlan {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            vid: _from_map!(v, "vid", <u32>::try_from)?,
            name: _from_map!(v, "name", <String>::try_from)?,
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize)]
pub struct NmLldpNeighbor8023MacPhyConf {
    pub autoneg: Option<u32>,
    pub pmd_autoneg_cap: Option<u32>,
    pub operational_mau_type: Option<u32>,
}

impl TryFrom<zvariant::OwnedValue> for NmLldpNeighbor8023MacPhyConf {
    type Error = NmError;
    fn try_from(v: zvariant::OwnedValue) -> Result<Self, Self::Error> {
        let mut v = DbusDictionary::try_from(v)?;
        Ok(Self {
            autoneg: _from_map!(v, "autoneg", <u32>::try_from)?,
            pmd_autoneg_cap: _from_map!(v, "pmd-autoneg-cap", <u32>::try_from)?,
            operational_mau_type: _from_map!(
                v,
                "operational-mau-type",
                <u32>::try_from
            )?,
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize)]
pub struct NmLldpNeighbor8023PowerViaMdi {
    pub mdi_power_support: Option<u32>,
    pub pse_power_pair: Option<u32>,
    pub power_class: Option<u32>,
}

impl TryFrom<zvariant::OwnedValue> for NmLldpNeighbor8023PowerViaMdi {
    type Error = NmError;
    fn try_from(v: zvariant::OwnedValue) -> Result<Self, Self::Error> {
        let mut v = DbusDictionary::try_from(v)?;
        Ok(Self {
            mdi_power_support: _from_map!(
                v,
                "mdi-power-support",
                <u32>::try_from
            )?,
            pse_power_pair: _from_map!(v, "pse-power-pair", <u32>::try_from)?,
            power_class: _from_map!(v, "power-class", <u32>::try_from)?,
        })
    }
}

#[derive(Debug, Clone, PartialEq, Eq, Default, Deserialize)]
#[serde(try_from = "DbusDictionary")]
pub struct NmLldpNeighborMgmtAddr {
    pub address_subtype: Option<u32>,
    pub address: Option<Vec<u8>>,
    pub interface_number_subtype: Option<u32>,
    pub interface_number: Option<u32>,
}

fn parse_mgmt_addrs(
    v: zvariant::OwnedValue,
) -> Result<Vec<NmLldpNeighborMgmtAddr>, NmError> {
    let mut ret = Vec::new();
    for value in Vec::<DbusDictionary>::try_from(v)? {
        ret.push(NmLldpNeighborMgmtAddr::try_from(value)?);
    }
    Ok(ret)
}

impl TryFrom<DbusDictionary> for NmLldpNeighborMgmtAddr {
    type Error = NmError;
    fn try_from(mut v: DbusDictionary) -> Result<Self, Self::Error> {
        Ok(Self {
            address_subtype: _from_map!(v, "address-subtype", <u32>::try_from)?,
            address: _from_map!(v, "address", <Vec<u8>>::try_from)?,
            interface_number_subtype: _from_map!(
                v,
                "interface-number-subtype",
                <u32>::try_from
            )?,
            interface_number: _from_map!(
                v,
                "interface-number",
                <u32>::try_from
            )?,
        })
    }
}
