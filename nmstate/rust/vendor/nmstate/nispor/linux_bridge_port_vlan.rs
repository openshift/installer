use crate::{
    BridgePortTunkTag, BridgePortVlanConfig, BridgePortVlanMode,
    BridgePortVlanRange,
};

pub(crate) fn parse_port_vlan_conf(
    np_vlan_entries: &[nispor::BridgeVlanEntry],
) -> Option<BridgePortVlanConfig> {
    let mut ret = BridgePortVlanConfig::new();
    let mut is_native = false;
    let mut trunk_tags = Vec::new();
    let is_access_port = is_access_port(np_vlan_entries);

    for np_vlan_entry in np_vlan_entries {
        let (vlan_min, vlan_max) = get_vlan_tag_range(np_vlan_entry);
        if vlan_min == 1 && vlan_max == 1 {
            continue;
        }
        if is_access_port {
            ret.tag = Some(vlan_min);
        } else if np_vlan_entry.is_pvid && np_vlan_entry.is_egress_untagged {
            ret.tag = Some(vlan_max);
            is_native = true;
        } else if vlan_min == vlan_max {
            trunk_tags.push(BridgePortTunkTag::Id(vlan_min));
        } else {
            trunk_tags.push(BridgePortTunkTag::IdRange(BridgePortVlanRange {
                max: vlan_max,
                min: vlan_min,
            }));
        }
    }
    if trunk_tags.is_empty() {
        ret.mode = Some(BridgePortVlanMode::Access);
    } else {
        ret.mode = Some(BridgePortVlanMode::Trunk);
        ret.enable_native = Some(is_native);
    }
    if ret.mode == Some(BridgePortVlanMode::Access)
        && trunk_tags.is_empty()
        && ret.tag.is_none()
    {
        None
    } else {
        ret.trunk_tags = Some(trunk_tags);

        Some(ret)
    }
}

fn is_access_port(np_vlan_entries: &[nispor::BridgeVlanEntry]) -> bool {
    np_vlan_entries.len() == 1
        && np_vlan_entries[0].is_pvid
        && np_vlan_entries[0].is_egress_untagged
}

fn get_vlan_tag_range(np_vlan_entry: &nispor::BridgeVlanEntry) -> (u16, u16) {
    np_vlan_entry.vid_range.unwrap_or((
        np_vlan_entry.vid.unwrap_or(1),
        np_vlan_entry.vid.unwrap_or(1),
    ))
}
