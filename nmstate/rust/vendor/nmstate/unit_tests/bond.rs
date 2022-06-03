use crate::{
    BondAdSelect, BondAllPortsActive, BondArpValidate, BondFailOverMac,
    BondInterface, BondLacpRate, BondPrimaryReselect, ErrorKind, Interface,
};

#[test]
fn test_bond_validate_mac_restricted_with_mac_undefined() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: active
"#,
    )
    .unwrap();
    iface.validate(None).unwrap();
}

#[test]
fn test_bond_validate_mac_restricted_with_mac_defined() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
mac-address: 00:01:02:03:04:05
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: active
"#,
    )
    .unwrap();
    let result = iface.validate(None);
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_bond_validate_mac_restricted_with_mac_defined_for_exist_bond() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
mac-address: 00:01:02:03:04:05
"#,
    )
    .unwrap();
    let current_iface: Interface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: active
"#,
    )
    .unwrap();
    let result = iface.validate(Some(&current_iface));
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_bond_validate_mac_restricted_with_mac_defined_changing_mode() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
mac-address: 00:01:02:03:04:05
link-aggregation:
  mode: 802.3ad
"#,
    )
    .unwrap();
    let current_iface: Interface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: active
"#,
    )
    .unwrap();
    iface.validate(Some(&current_iface)).unwrap();
}

#[test]
fn test_bond_validate_mac_restricted_with_mac_defined_changing_opt() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
mac-address: 00:01:02:03:04:05
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: follow
"#,
    )
    .unwrap();
    let current_iface: Interface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: active-backup
  options:
    fail_over_mac: active
"#,
    )
    .unwrap();
    iface.validate(Some(&current_iface)).unwrap();
}

#[test]
fn test_bond_validate_bond_mode_not_defined_for_new_iface() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
"#,
    )
    .unwrap();
    let result = iface.validate(None);
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_bond_validate_ad_actor_system_mac_address() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: 802.3ad
  options:
    ad_actor_system: "01:00:5E:00:0f:01"
"#,
    )
    .unwrap();
    let result = iface.validate(None);
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_bond_validate_miimon_and_arp_interval() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: 802.3ad
  options:
    miimon: 100
    arp_interval: 60
"#,
    )
    .unwrap();
    let result = iface.validate(None);
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_bond_stringlized_attributes() {
    let iface: BondInterface = serde_yaml::from_str(
        r#"---
name: bond99
type: bond
state: up
link-aggregation:
  mode: 802.3ad
  options:
    ad_actor_sys_prio: "555"
    ad_select: "1"
    ad_user_port_key: "16"
    all_slaves_active: "1"
    arp_interval: "100"
    arp_validate: "1"
    downdelay: "50"
    fail_over_mac: "1"
    lacp_rate: "0"
    lp_interval: "101"
    miimon: "200"
    min_links: "2"
    num_grat_arp: "3"
    num_unsol_na: "4"
    packets_per_slave: "1000"
    primary_reselect: "1"
    resend_igmp: "103"
    tlb_dynamic_lb: "true"
    updelay: "104"
    use_carrier: "false"
"#,
    )
    .unwrap();
    let bond_opts = iface.bond.unwrap().options.unwrap();
    assert_eq!(bond_opts.ad_actor_sys_prio, Some(555));
    assert_eq!(bond_opts.ad_select, Some(BondAdSelect::Bandwidth));
    assert_eq!(bond_opts.ad_user_port_key, Some(16));
    assert_eq!(
        bond_opts.all_slaves_active,
        Some(BondAllPortsActive::Delivered)
    );
    assert_eq!(bond_opts.arp_interval, Some(100));
    assert_eq!(bond_opts.arp_validate, Some(BondArpValidate::Active));
    assert_eq!(bond_opts.downdelay, Some(50));
    assert_eq!(bond_opts.fail_over_mac, Some(BondFailOverMac::Active));
    assert_eq!(bond_opts.lacp_rate, Some(BondLacpRate::Slow));
    assert_eq!(bond_opts.lp_interval, Some(101));
    assert_eq!(bond_opts.miimon, Some(200));
    assert_eq!(bond_opts.min_links, Some(2));
    assert_eq!(bond_opts.num_grat_arp, Some(3));
    assert_eq!(bond_opts.num_unsol_na, Some(4));
    assert_eq!(bond_opts.packets_per_slave, Some(1000));
    assert_eq!(
        bond_opts.primary_reselect,
        Some(BondPrimaryReselect::Better)
    );
    assert_eq!(bond_opts.resend_igmp, Some(103));
    assert_eq!(bond_opts.tlb_dynamic_lb, Some(true));
    assert_eq!(bond_opts.updelay, Some(104));
    assert_eq!(bond_opts.use_carrier, Some(false));
}
