use crate::{
    ifaces::get_ignored_ifaces, BridgePortTunkTag, BridgePortVlanRange,
    InterfaceType, Interfaces, LinuxBridgeInterface,
    LinuxBridgeMulticastRouterType,
};

#[test]
fn test_linux_bridge_ignore_port() {
    let mut ifaces: Interfaces = serde_yaml::from_str(
        r#"---
- name: eth1
  type: ethernet
  state: ignore
- name: eth2
  type: ethernet
- name: br0
  type: linux-bridge
  state: up
  bridge:
    port:
    - name: eth2
"#,
    )
    .unwrap();
    let mut cur_ifaces: Interfaces = serde_yaml::from_str(
        r#"---
- name: eth1
  type: ethernet
  state: up
- name: eth2
  type: ethernet
  state: up
- name: br0
  type: linux-bridge
  state: up
  bridge:
    port:
    - name: eth1
    - name: eth2
"#,
    )
    .unwrap();

    let (ignored_kernel_ifaces, ignored_user_ifaces) =
        get_ignored_ifaces(&ifaces, &cur_ifaces);

    assert_eq!(ignored_kernel_ifaces, vec!["eth1".to_string()]);
    assert!(ignored_user_ifaces.is_empty());

    ifaces.remove_ignored_ifaces(&ignored_kernel_ifaces, &ignored_user_ifaces);
    cur_ifaces
        .remove_ignored_ifaces(&ignored_kernel_ifaces, &ignored_user_ifaces);

    let (add_ifaces, chg_ifaces, del_ifaces) =
        ifaces.gen_state_for_apply(&cur_ifaces, false).unwrap();

    assert!(!ifaces.kernel_ifaces.contains_key("eth1"));
    assert!(!cur_ifaces.kernel_ifaces.contains_key("eth1"));
    assert_eq!(ifaces.kernel_ifaces["br0"].ports(), Some(vec!["eth2"]));
    assert_eq!(cur_ifaces.kernel_ifaces["br0"].ports(), Some(vec!["eth2"]));
    assert!(!add_ifaces.kernel_ifaces.contains_key("eth1"));
    assert!(!chg_ifaces.kernel_ifaces.contains_key("eth1"));
    assert!(!del_ifaces.kernel_ifaces.contains_key("eth1"));
}

#[test]
fn test_linux_bridge_verify_ignore_port() {
    let ifaces: Interfaces = serde_yaml::from_str(
        r#"---
- name: eth1
  type: ethernet
  state: ignore
- name: eth2
  type: ethernet
- name: br0
  type: linux-bridge
  state: up
  bridge:
    port:
    - name: eth2
"#,
    )
    .unwrap();
    let cur_ifaces: Interfaces = serde_yaml::from_str(
        r#"---
- name: eth1
  type: ethernet
  state: up
- name: eth2
  type: ethernet
  state: up
- name: br0
  type: linux-bridge
  state: up
  bridge:
    port:
    - name: eth1
    - name: eth2
"#,
    )
    .unwrap();

    ifaces.verify(&cur_ifaces).unwrap();
}

#[test]
fn test_linux_bridge_stringlized_attributes() {
    let iface: LinuxBridgeInterface = serde_yaml::from_str(
        r#"---
name: br0
type: linux-bridge
state: up
bridge:
  options:
    group-forward-mask: "300"
    group-fwd-mask: "301"
    hash-max: "302"
    mac-ageing-time: "303"
    multicast-last-member-count: "304"
    multicast-last-member-interval: "305"
    multicast-membership-interval: "306"
    multicast-querier: "1"
    multicast-querier-interval: "307"
    multicast-query-interval: "308"
    multicast-query-response-interval: "309"
    multicast-query-use-ifaddr: "yes"
    multicast-snooping: "no"
    multicast-startup-query-count: "310"
    multicast-startup-query-interval: "311"
    stp:
      enabled: "false"
      forward-delay: "16"
      hello-time: "2"
      max-age: "20"
      priority: "32768"
  port:
  - name: eth1
    stp-hairpin-mode: "false"
    stp-path-cost: "100"
    stp-priority: "101"
    vlan:
      enable-native: "true"
      tag: "102"
      trunk-tags:
      - id: "103"
      - id-range:
          max: "1024"
          min: "105"
"#,
    )
    .unwrap();

    let br_conf = iface.bridge.unwrap();
    let port_conf = &br_conf.port.as_ref().unwrap()[0];
    let vlan_conf = port_conf.vlan.as_ref().unwrap();
    let opts = br_conf.options.as_ref().unwrap();
    let stp_opts = opts.stp.as_ref().unwrap();

    assert_eq!(port_conf.stp_hairpin_mode, Some(false));
    assert_eq!(port_conf.stp_path_cost, Some(100));
    assert_eq!(port_conf.stp_priority, Some(101));
    assert_eq!(vlan_conf.enable_native, Some(true));
    assert_eq!(vlan_conf.tag, Some(102));
    assert_eq!(
        &vlan_conf.trunk_tags.as_ref().unwrap()[0],
        &BridgePortTunkTag::Id(103)
    );
    assert_eq!(
        &vlan_conf.trunk_tags.as_ref().unwrap()[1],
        &BridgePortTunkTag::IdRange(BridgePortVlanRange {
            max: 1024,
            min: 105
        })
    );

    assert_eq!(stp_opts.enabled, Some(false));
    assert_eq!(stp_opts.forward_delay, Some(16));
    assert_eq!(stp_opts.hello_time, Some(2));
    assert_eq!(stp_opts.max_age, Some(20));
    assert_eq!(stp_opts.priority, Some(32768));

    assert_eq!(opts.group_forward_mask, Some(300));
    assert_eq!(opts.group_fwd_mask, Some(301));
    assert_eq!(opts.hash_max, Some(302));
    assert_eq!(opts.mac_ageing_time, Some(303));
    assert_eq!(opts.multicast_last_member_count, Some(304));
    assert_eq!(opts.multicast_last_member_interval, Some(305));
    assert_eq!(opts.multicast_membership_interval, Some(306));
    assert_eq!(opts.multicast_querier, Some(true));
    assert_eq!(opts.multicast_querier_interval, Some(307));
    assert_eq!(opts.multicast_query_interval, Some(308));
    assert_eq!(opts.multicast_query_response_interval, Some(309));
    assert_eq!(opts.multicast_query_use_ifaddr, Some(true));
    assert_eq!(opts.multicast_snooping, Some(false));
    assert_eq!(opts.multicast_startup_query_count, Some(310));
    assert_eq!(opts.multicast_startup_query_interval, Some(311));
}

#[test]
fn test_linux_bridge_partial_ignored() {
    let mut current = serde_yaml::from_str::<Interfaces>(
        r#"---
- name: eth1
  type: ethernet
  state: ignore
- name: eth2
  type: ethernet
  state: ignore
- name: br0
  type: linux-bridge
  state: ignore
  bridge:
    port:
    - name: eth1
    - name: eth2
"#,
    )
    .unwrap();
    let mut desired = serde_yaml::from_str::<Interfaces>(
        r#"---
- name: br0
  type: linux-bridge
  state: up
- name: eth1
  type: ethernet
  state: up
"#,
    )
    .unwrap();
    let (kernel_ifaces, user_ifaces) = get_ignored_ifaces(&desired, &current);

    assert_eq!(kernel_ifaces, vec!["eth2".to_string()]);
    assert_eq!(user_ifaces, vec![]);

    current.remove_ignored_ifaces(&kernel_ifaces, &user_ifaces);
    desired.remove_ignored_ifaces(&kernel_ifaces, &user_ifaces);

    let (add_ifaces, chg_ifaces, del_ifaces) =
        desired.gen_state_for_apply(&current, false).unwrap();

    assert!(add_ifaces.kernel_ifaces.is_empty());
    assert!(del_ifaces.kernel_ifaces.is_empty());
    assert_eq!(
        chg_ifaces.kernel_ifaces["eth1"].base_iface().controller,
        Some("br0".to_string())
    );
    assert_eq!(
        chg_ifaces.kernel_ifaces["eth1"]
            .base_iface()
            .controller_type,
        Some(InterfaceType::LinuxBridge)
    );
    assert!(chg_ifaces.kernel_ifaces.contains_key("br0"));
}

#[test]
fn test_linux_bridge_interger_multicast_router() {
    let iface: LinuxBridgeInterface = serde_yaml::from_str(
        r#"---
name: br0
type: linux-bridge
state: up
bridge:
  options:
    multicast-router: 0
"#,
    )
    .unwrap();

    assert_eq!(
        iface
            .bridge
            .unwrap()
            .options
            .as_ref()
            .unwrap()
            .multicast_router,
        Some(LinuxBridgeMulticastRouterType::Disabled)
    );
}
