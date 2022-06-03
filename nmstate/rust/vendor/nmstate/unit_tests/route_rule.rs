use crate::{
    unit_tests::testlib::new_eth_iface, InterfaceType, Interfaces,
    NetworkState, RouteEntry, RouteRuleEntry, RouteRules, Routes,
};

const TEST_NIC: &str = "eth1";
const TEST_IPV4_NET1: &str = "192.0.2.0/24";
const TEST_IPV4_ADDR1: &str = "198.51.100.1";
const TEST_IPV6_NET1: &str = "2001:db8:1::/64";
const TEST_IPV6_ADDR1: &str = "2001:db8:0::1";
const TEST_TABLE_ID1: u32 = 101;
const TEST_TABLE_ID2: u32 = 102;

const TEST_RULE_IPV6_FROM: &str = "2001:db8:1::2";
const TEST_RULE_IPV4_FROM: &str = "192.0.2.1";
const TEST_RULE_IPV6_TO: &str = "2001:db8:2::2";
const TEST_RULE_IPV4_TO: &str = "198.51.100.1";
const TEST_RULE_PRIORITY1: i64 = 201;
const TEST_RULE_PRIORITY2: i64 = 202;

const TEST_ROUTE_METRIC: i64 = 100;

#[test]
fn test_sort_uniqe_route_rules() {
    let mut test_routes = gen_test_rule_entries();
    test_routes.reverse();
    test_routes.extend(gen_test_rule_entries());
    test_routes.sort_unstable();
    test_routes.dedup();

    assert_eq!(test_routes, gen_test_rule_entries());
}

#[test]
fn test_add_rules_to_new_interface() {
    let cur_net_state = NetworkState::new();

    let des_iface = new_eth_iface(TEST_NIC);
    let mut des_ifaces = Interfaces::new();
    des_ifaces.push(des_iface);
    let mut des_net_state = NetworkState::new();
    des_net_state.interfaces = des_ifaces;
    des_net_state.routes = gen_test_routes_conf();
    des_net_state.rules = gen_test_rules_conf();

    let (add_net_state, chg_net_state, del_net_state) =
        des_net_state.gen_state_for_apply(&cur_net_state).unwrap();

    assert_eq!(chg_net_state, NetworkState::new());
    assert_eq!(del_net_state, NetworkState::new());

    let add_ifaces = add_net_state.interfaces.to_vec();

    assert_eq!(add_ifaces.len(), 1);
    assert_eq!(add_ifaces[0].name(), TEST_NIC);
    let config_rules = add_ifaces[0].base_iface().rules.as_ref().unwrap();

    println!("{:?}", config_rules);
    assert_eq!(config_rules.len(), 2);

    assert_eq!(
        config_rules[0].ip_from.as_ref().unwrap().as_str(),
        TEST_RULE_IPV6_FROM,
    );
    assert_eq!(
        config_rules[0].ip_to.as_ref().unwrap().as_str(),
        TEST_RULE_IPV6_TO,
    );
    assert_eq!(config_rules[0].priority.unwrap(), TEST_RULE_PRIORITY1);
    assert_eq!(config_rules[0].table_id.unwrap(), TEST_TABLE_ID1);
    assert_eq!(
        config_rules[1].ip_from.as_ref().unwrap().as_str(),
        TEST_RULE_IPV4_FROM,
    );
    assert_eq!(
        config_rules[1].ip_to.as_ref().unwrap().as_str(),
        TEST_RULE_IPV4_TO,
    );
    assert_eq!(config_rules[1].priority.unwrap(), TEST_RULE_PRIORITY2);
    assert_eq!(config_rules[1].table_id.unwrap(), TEST_TABLE_ID2);
}

fn gen_test_routes_conf() -> Routes {
    let mut ret = Routes::new();
    ret.running = Some(gen_test_route_entries());
    ret.config = Some(gen_test_route_entries());
    ret
}

fn gen_test_route_entries() -> Vec<RouteEntry> {
    vec![
        gen_route_entry(
            TEST_IPV6_NET1,
            TEST_NIC,
            TEST_IPV6_ADDR1,
            TEST_TABLE_ID1,
        ),
        gen_route_entry(
            TEST_IPV4_NET1,
            TEST_NIC,
            TEST_IPV4_ADDR1,
            TEST_TABLE_ID2,
        ),
    ]
}

fn gen_route_entry(
    dst: &str,
    next_hop_iface: &str,
    next_hop_addr: &str,
    table_id: u32,
) -> RouteEntry {
    let mut ret = RouteEntry::new();
    ret.destination = Some(dst.to_string());
    ret.next_hop_iface = Some(next_hop_iface.to_string());
    ret.next_hop_addr = Some(next_hop_addr.to_string());
    ret.metric = Some(TEST_ROUTE_METRIC);
    ret.table_id = Some(table_id);
    ret
}

fn gen_test_rules_conf() -> RouteRules {
    RouteRules {
        config: Some(gen_test_rule_entries()),
    }
}

fn gen_test_rule_entries() -> Vec<RouteRuleEntry> {
    vec![
        gen_rule_entry(
            TEST_RULE_IPV6_FROM,
            TEST_RULE_IPV6_TO,
            TEST_RULE_PRIORITY1,
            TEST_TABLE_ID1,
        ),
        gen_rule_entry(
            TEST_RULE_IPV4_FROM,
            TEST_RULE_IPV4_TO,
            TEST_RULE_PRIORITY2,
            TEST_TABLE_ID2,
        ),
    ]
}

fn gen_rule_entry(
    ip_from: &str,
    ip_to: &str,
    priority: i64,
    table_id: u32,
) -> RouteRuleEntry {
    RouteRuleEntry {
        state: None,
        ip_from: Some(ip_from.to_string()),
        ip_to: Some(ip_to.to_string()),
        table_id: Some(table_id),
        priority: Some(priority),
    }
}

#[test]
fn test_route_rule_stringlized_attributes() {
    let rule: RouteRuleEntry = serde_yaml::from_str(
        r#"
priority: "500"
route-table: "129"
"#,
    )
    .unwrap();
    assert_eq!(rule.table_id, Some(129));
    assert_eq!(rule.priority, Some(500));
}

#[test]
fn test_route_rule_use_auto_route_table_id() {
    let current: NetworkState = serde_yaml::from_str(
        r#"
---
interfaces:
  - name: br0
    type: ovs-interface
    state: up
    ipv4:
      enabled: true
      dhcp: true
      auto-dns: false
      auto-routes: true
      auto-gateway: true
      auto-route-table-id: 500
    ipv6:
      enabled: false
  - name: br0
    type: ovs-bridge
    state: up
    bridge:
      port:
        - name: br0
"#,
    )
    .unwrap();

    let desire: NetworkState = serde_yaml::from_str(
        r#"
---
route-rules:
  config:
    - route-table: 500
      priority: 3200
      ip-to: 192.0.3.0/24
    - route-table: 500
      priority: 3200
      ip-from: 192.0.3.0/24
"#,
    )
    .unwrap();

    let expected_rules: Vec<RouteRuleEntry> = serde_yaml::from_str(
        r#"
- route-table: 500
  priority: 3200
  ip-to: 192.0.3.0/24
- route-table: 500
  priority: 3200
  ip-from: 192.0.3.0/24
"#,
    )
    .unwrap();

    let (_, chg_net_state, _) = desire.gen_state_for_apply(&current).unwrap();

    let ovs_iface = &chg_net_state.interfaces.kernel_ifaces["br0"];

    assert_eq!(ovs_iface.iface_type(), InterfaceType::OvsInterface);
    assert_eq!(ovs_iface.base_iface().rules, Some(expected_rules));
}

#[test]
fn test_route_rule_ignore_absent_ifaces() {
    let desired: NetworkState = serde_yaml::from_str(
        r#"
interfaces:
- name: br0
  state: absent
  type: linux-bridge
route-rules:
  config:
  - route-table: 200
    state: absent
"#,
    )
    .unwrap();

    let current: NetworkState = serde_yaml::from_str(
        r#"
interfaces:
- name: eth1
  type: ethernet
  state: up
- name: br0
  type: linux-bridge
  state: up
  ipv4:
    address:
    - ip: 192.0.2.251
      prefix-length: 24
    dhcp: false
    enabled: true
  bridge:
    options:
      stp:
        enabled: false
    port:
    - name: eth1
routes:
  config:
    - destination: 198.51.100.0/24
      metric: 150
      next-hop-address: 192.0.2.1
      next-hop-interface: br0
      table-id: 200
route-rules:
  config:
    - ip-from: 192.51.100.2/32
      route-table: 200
"#,
    )
    .unwrap();

    let (add_net_state, chg_net_state, del_net_state) =
        desired.gen_state_for_apply(&current).unwrap();

    println!("add_net_state {:?}", add_net_state);
    println!("chg_net_state {:?}", chg_net_state);
    println!("del_net_state {:?}", del_net_state);

    assert!(add_net_state.interfaces.to_vec().is_empty());
    assert!(chg_net_state.interfaces.to_vec().len() == 1);
    assert!(chg_net_state.interfaces.kernel_ifaces["eth1"].is_up());

    assert!(del_net_state.interfaces.to_vec().len() == 1);
    assert!(del_net_state.interfaces.kernel_ifaces["br0"].is_absent());
}
