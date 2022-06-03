use crate::nm::nm_dbus::{
    NmConnection, NmSettingBridge, NmSettingBridgeVlanRange,
};

use crate::{
    BridgePortTunkTag, BridgePortVlanConfig, BridgePortVlanMode,
    LinuxBridgeInterface, LinuxBridgeOptions, LinuxBridgeStpOptions,
};

pub(crate) fn gen_nm_br_setting(
    br_iface: &LinuxBridgeInterface,
    nm_conn: &mut NmConnection,
) {
    let mut nm_br_set = nm_conn.bridge.as_ref().cloned().unwrap_or_default();

    if let Some(br_conf) = br_iface.bridge.as_ref() {
        if let Some(br_opts) = br_conf.options.as_ref() {
            apply_br_options(&mut nm_br_set, br_opts)
        }

        if br_conf.port.is_some() {
            nm_br_set.vlan_filtering =
                Some(br_iface.vlan_filtering_is_enabled());
        }
    }

    nm_conn.bridge = Some(nm_br_set);
}

fn apply_br_options(
    nm_br_set: &mut NmSettingBridge,
    br_opts: &LinuxBridgeOptions,
) {
    if let Some(v) = br_opts.group_addr.as_ref() {
        nm_br_set.group_address = Some(v.to_string());
    }
    if let Some(v) = br_opts.group_forward_mask.as_ref() {
        nm_br_set.group_forward_mask = Some((*v).into());
    }
    if let Some(v) = br_opts.group_fwd_mask.as_ref() {
        nm_br_set.group_forward_mask = Some((*v).into());
    }
    if let Some(v) = br_opts.hash_max.as_ref() {
        nm_br_set.multicast_hash_max = Some(*v);
    }
    if let Some(v) = br_opts.mac_ageing_time.as_ref() {
        nm_br_set.ageing_time = Some(*v);
    }
    if let Some(v) = br_opts.multicast_last_member_count.as_ref() {
        nm_br_set.multicast_last_member_count = Some(*v);
    }
    if let Some(v) = br_opts.multicast_last_member_interval.as_ref() {
        nm_br_set.multicast_last_member_interval = Some(*v);
    }
    if let Some(v) = br_opts.multicast_membership_interval.as_ref() {
        nm_br_set.multicast_membership_interval = Some(*v);
    }
    if let Some(v) = br_opts.multicast_querier.as_ref() {
        nm_br_set.multicast_querier = Some(*v);
    }
    if let Some(v) = br_opts.multicast_querier_interval.as_ref() {
        nm_br_set.multicast_querier_interval = Some(*v);
    }
    if let Some(v) = br_opts.multicast_query_interval.as_ref() {
        nm_br_set.multicast_query_interval = Some(*v);
    }
    if let Some(v) = br_opts.multicast_query_response_interval.as_ref() {
        nm_br_set.multicast_query_response_interval = Some(*v);
    }
    if let Some(v) = br_opts.multicast_query_use_ifaddr.as_ref() {
        nm_br_set.multicast_query_use_ifaddr = Some(*v);
    }
    if let Some(v) = br_opts.multicast_router.as_ref() {
        nm_br_set.multicast_router = Some(format!("{}", v));
    }
    if let Some(v) = br_opts.multicast_snooping.as_ref() {
        nm_br_set.multicast_snooping = Some(*v);
    }
    if let Some(v) = br_opts.multicast_startup_query_count.as_ref() {
        nm_br_set.multicast_startup_query_count = Some(*v);
    }
    if let Some(v) = br_opts.multicast_startup_query_interval.as_ref() {
        nm_br_set.multicast_startup_query_interval = Some(*v);
    }

    if let Some(stp_opts) = br_opts.stp.as_ref() {
        apply_stp_setting(nm_br_set, stp_opts);
    }
}

fn apply_stp_setting(
    nm_set: &mut NmSettingBridge,
    opts: &LinuxBridgeStpOptions,
) {
    if let Some(v) = opts.enabled {
        nm_set.stp = Some(v);
    }
    if let Some(v) = opts.forward_delay {
        nm_set.forward_delay = Some(v.into());
    }
    if let Some(v) = opts.hello_time {
        nm_set.hello_time = Some(v.into())
    }
    if let Some(v) = opts.max_age {
        nm_set.max_age = Some(v.into())
    };
    if let Some(v) = opts.priority {
        nm_set.priority = Some(v.into())
    };
}

pub(crate) fn gen_nm_br_port_setting(
    br_iface: &LinuxBridgeInterface,
    nm_conn: &mut NmConnection,
) {
    let mut nm_set = nm_conn.bridge_port.as_ref().cloned().unwrap_or_default();
    let br_port_conf = if let Some(i) = nm_conn
        .iface_name()
        .and_then(|iface_name| br_iface.get_port_conf(iface_name))
    {
        i
    } else {
        return;
    };

    if let Some(v) = br_port_conf.stp_hairpin_mode {
        nm_set.hairpin_mode = Some(v);
    }

    if let Some(v) = br_port_conf.stp_path_cost {
        nm_set.path_cost = Some(v);
    }

    if let Some(v) = br_port_conf.stp_priority {
        nm_set.priority = Some(v.into());
    }
    if let Some(v) = br_port_conf.vlan.as_ref() {
        nm_set.vlans = Some(nmstate_port_vlans_to_nm_vlan_range(v));
    }

    nm_conn.bridge_port = Some(nm_set);
}

fn nmstate_port_vlans_to_nm_vlan_range(
    port_vlan_conf: &BridgePortVlanConfig,
) -> Vec<NmSettingBridgeVlanRange> {
    let mut ret = Vec::new();
    match port_vlan_conf.mode {
        Some(BridgePortVlanMode::Trunk) => {
            if let Some(trunk_tags) = &port_vlan_conf.trunk_tags {
                for trunk_tag in trunk_tags.as_slice() {
                    ret.push(trunk_tag_to_nm_vlan_range(trunk_tag));
                }
            }
            if let Some(t) = port_vlan_conf.tag {
                ret.push(access_tag_to_nm_vlan_range(t))
            }
        }
        Some(BridgePortVlanMode::Access) => {
            if let Some(t) = port_vlan_conf.tag {
                ret.push(access_tag_to_nm_vlan_range(t))
            };
        }
        _ => (),
    }

    ret
}

fn trunk_tag_to_nm_vlan_range(
    trunk_tag: &BridgePortTunkTag,
) -> NmSettingBridgeVlanRange {
    let mut ret = NmSettingBridgeVlanRange::default();
    let (vid_min, vid_max) = trunk_tag.get_vlan_tag_range();
    ret.vid_start = vid_min;
    ret.vid_end = vid_max;
    ret.pvid = false;
    ret.untagged = false;
    ret
}

fn access_tag_to_nm_vlan_range(tag: u16) -> NmSettingBridgeVlanRange {
    let mut ret = NmSettingBridgeVlanRange::default();
    ret.vid_start = tag;
    ret.vid_end = tag;
    ret.pvid = true;
    ret.untagged = true;
    ret
}
