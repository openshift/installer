mod testlib;

use crate::testlib::{
    context::with_clean_up_afterwords,
    iface::{assert_iface_bridge, assert_iface_exists, assert_iface_missing},
    ip::assert_ip,
};
use nmstate::NetworkState;

const BRIDGE_CREATE_WITH_STATIC_IP_YML: &str = r#"---
interfaces:
- name: br0
  type: linux-bridge
  state: up
  bridge:
    port:
      - name: eth2
      - name: eth1
    options:
      stp:
        enabled: false
  ipv4:
    address:
    - ip: 192.0.2.251
      prefix-length: 24
    dhcp: false
    enabled: true
  ipv6:
    address:
      - ip: 2001:db8:1::1
        prefix-length: 64
    autoconf: false
    dhcp: false
    enabled: true
- name: eth1
  type: veth
  state: up
  veth:
    peer: eth1.ep
- name: eth2
  type: veth
  state: up
  veth:
    peer: eth2.ep"#;

const CREATE_VETH_AND_ATTACH_TO_BRIDGE: &str = r#"---
interfaces:
- name: eth3
  type: veth
  state: up
  veth:
    peer: eth3.ep
  controller: br0"#;

const CLEAN_UP_YML: &str = r#"---
interfaces:
- name: br0
  state: absent
- name: eth1
  state: absent
- name: eth2
  state: absent
- name: eth3
  state: absent"#;

#[test]
fn test_create_bridge_with_veth_port_and_static_ip() {
    with_clean_up_afterwords(
        || {
            let mut state: NetworkState =
                serde_yaml::from_str(BRIDGE_CREATE_WITH_STATIC_IP_YML).unwrap();
            state.set_kernel_only(true);
            state.apply().unwrap();
            assert_iface_exists("br0");
            assert_iface_exists("eth1");
            assert_iface_exists("eth1.ep");
            assert_iface_exists("eth2");
            assert_iface_exists("eth2.ep");
            assert_iface_bridge("br0", &["eth1", "eth2"]);
            assert_ip("br0", &["192.0.2.251/24", "2001:db8:1::1/64"]);
        },
        clean_up,
    );
}

#[test]
fn test_attach_veth_to_bridge_via_controller_prop() {
    with_clean_up_afterwords(
        || {
            let mut state: NetworkState =
                serde_yaml::from_str(BRIDGE_CREATE_WITH_STATIC_IP_YML).unwrap();
            state.set_kernel_only(true);
            state.apply().unwrap();
            let mut state: NetworkState =
                serde_yaml::from_str(CREATE_VETH_AND_ATTACH_TO_BRIDGE).unwrap();
            state.set_kernel_only(true);
            state.apply().unwrap();
            assert_iface_exists("br0");
            assert_iface_exists("eth1");
            assert_iface_exists("eth1.ep");
            assert_iface_exists("eth2");
            assert_iface_exists("eth2.ep");
            assert_iface_exists("eth3");
            assert_iface_exists("eth3.ep");
            assert_iface_bridge("br0", &["eth1", "eth2", "eth3"]);
            assert_ip("br0", &["192.0.2.251/24", "2001:db8:1::1/64"]);
        },
        clean_up,
    );
}

fn clean_up() {
    let mut state: NetworkState = serde_yaml::from_str(CLEAN_UP_YML).unwrap();
    state.set_kernel_only(true);
    state.apply().unwrap();
    assert_iface_missing("br0");
    assert_iface_missing("eth1");
    assert_iface_missing("eth1.ep");
    assert_iface_missing("eth2");
    assert_iface_missing("eth2.ep");
}
