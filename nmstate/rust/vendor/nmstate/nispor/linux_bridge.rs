use std::convert::TryFrom;

use log::warn;

use crate::{
    nispor::linux_bridge_port_vlan::parse_port_vlan_conf, BaseInterface,
    LinuxBridgeConfig, LinuxBridgeInterface, LinuxBridgeMulticastRouterType,
    LinuxBridgeOptions, LinuxBridgePortConfig, LinuxBridgeStpOptions,
};

pub(crate) fn np_bridge_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> LinuxBridgeInterface {
    let mut br_iface = LinuxBridgeInterface::new();
    let mut br_conf = LinuxBridgeConfig::new();
    br_iface.base = base_iface;
    br_conf.options = Some(np_bridge_options_to_nmstate(np_iface));
    if let Some(np_bridge) = &np_iface.bridge {
        br_conf.port = Some(
            np_bridge
                .ports
                .as_slice()
                .iter()
                .map(|iface_name| {
                    let mut port_conf = LinuxBridgePortConfig::new();
                    port_conf.name = iface_name.to_string();
                    port_conf
                })
                .collect(),
        );
    }
    br_iface.bridge = Some(br_conf);
    br_iface
}

pub(crate) fn append_bridge_port_config(
    br_iface: &mut LinuxBridgeInterface,
    np_iface: &nispor::Iface,
    port_np_ifaces: Vec<&nispor::Iface>,
) {
    let mut port_confs: Vec<LinuxBridgePortConfig> = Vec::new();
    for port_np_iface in port_np_ifaces {
        let mut port_conf = LinuxBridgePortConfig::new();
        port_conf.name = port_np_iface.name.to_string();
        if let Some(np_port_info) = &port_np_iface.bridge_port {
            port_conf.stp_hairpin_mode = Some(np_port_info.hairpin_mode);
            port_conf.stp_path_cost = Some(np_port_info.stp_path_cost);
            port_conf.stp_priority = Some(np_port_info.stp_priority);
            if np_iface
                .bridge
                .as_ref()
                .and_then(|br_info| br_info.vlan_filtering)
                == Some(true)
            {
                port_conf.vlan = np_port_info
                    .vlans
                    .as_ref()
                    .and_then(|v| parse_port_vlan_conf(v.as_slice()));
            }
        }
        port_confs.push(port_conf);
    }

    if let Some(mut br_conf) = br_iface.bridge.as_mut() {
        br_conf.port = Some(port_confs);
    }
}

fn np_bridge_options_to_nmstate(
    np_iface: &nispor::Iface,
) -> LinuxBridgeOptions {
    let mut options = LinuxBridgeOptions::default();
    if let Some(ref np_bridge) = &np_iface.bridge {
        options.stp = Some(get_stp_options(np_bridge));
        options.gc_timer = np_bridge.gc_timer;
        options.group_addr = np_bridge
            .group_addr
            .as_ref()
            .map(|addr| addr.to_uppercase());
        options.group_forward_mask = np_bridge.group_fwd_mask;
        options.group_fwd_mask = np_bridge.group_fwd_mask;
        options.hash_max = np_bridge.multicast_hash_max;
        options.hello_timer = np_bridge.hello_timer;
        options.mac_ageing_time = np_bridge.ageing_time.map(devide_by_user_hz);
        options.multicast_last_member_count =
            np_bridge.multicast_last_member_count;
        options.multicast_last_member_interval =
            np_bridge.multicast_last_member_interval;
        options.multicast_membership_interval =
            np_bridge.multicast_membership_interval;
        options.multicast_querier = np_bridge.multicast_querier;
        options.multicast_querier_interval =
            np_bridge.multicast_querier_interval;
        options.multicast_query_interval = np_bridge.multicast_query_interval;
        options.multicast_query_response_interval =
            np_bridge.multicast_query_response_interval;
        options.multicast_query_use_ifaddr =
            np_bridge.multicast_query_use_ifaddr;
        options.multicast_router =
            np_bridge.multicast_router.as_ref().and_then(|r| match r {
                nispor::BridgePortMulticastRouterType::Disabled => {
                    Some(LinuxBridgeMulticastRouterType::Disabled)
                }
                nispor::BridgePortMulticastRouterType::TempQuery => {
                    Some(LinuxBridgeMulticastRouterType::Auto)
                }
                nispor::BridgePortMulticastRouterType::Perm => {
                    Some(LinuxBridgeMulticastRouterType::Enabled)
                }
                _ => {
                    warn!("Unsupported linux bridge multicast router {:?}", r);
                    None
                }
            });
        options.multicast_snooping = np_bridge.multicast_snooping;
        options.multicast_startup_query_count =
            np_bridge.multicast_startup_query_count;
        options.multicast_startup_query_interval =
            np_bridge.multicast_startup_query_interval;
    }
    options
}

// The kernel is multiplying these bridge properties by USER_HZ, we should
// divide into seconds:
//   * forward_delay
//   * ageing_time
//   * hello_time
//   * max_age
fn devide_by_user_hz(v: u32) -> u32 {
    let user_hz = unsafe { libc::sysconf(libc::_SC_CLK_TCK) } as u32;
    v / user_hz
}

fn get_stp_options(np_bridge: &nispor::BridgeInfo) -> LinuxBridgeStpOptions {
    let mut stp_opt = LinuxBridgeStpOptions::new();
    stp_opt.enabled = Some(
        [
            Some(nispor::BridgeStpState::KernelStp),
            Some(nispor::BridgeStpState::UserStp),
        ]
        .contains(&np_bridge.stp_state),
    );
    stp_opt.forward_delay = np_bridge.forward_delay.map(|v| {
        u8::try_from(devide_by_user_hz(v))
            .unwrap_or(LinuxBridgeStpOptions::FORWARD_DELAY_MAX)
    });
    stp_opt.max_age = np_bridge.max_age.map(|v| {
        u8::try_from(devide_by_user_hz(v))
            .unwrap_or(LinuxBridgeStpOptions::MAX_AGE_MAX)
    });
    stp_opt.hello_time = np_bridge.hello_time.map(|v| {
        u8::try_from(devide_by_user_hz(v))
            .unwrap_or(LinuxBridgeStpOptions::HELLO_TIME_MAX)
    });
    stp_opt.priority = np_bridge.priority;
    stp_opt
}
