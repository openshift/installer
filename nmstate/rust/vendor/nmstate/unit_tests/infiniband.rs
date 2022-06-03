use crate::{ErrorKind, InfiniBandInterface, Interfaces};

#[test]
fn test_ib_autoremove_pkey_if_base_iface_removed() {
    let mut desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: absent
"#,
    )
    .unwrap();
    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
- name: mlx5_ib2.8001
  type: infiniband
  state: up
  infiniband:
    pkey: "0x8001"
    mode: "connected"
    base-iface: "mlx5_ib2"
"#,
    )
    .unwrap();
    let (_, _, del_ifaces) =
        desire.gen_state_for_apply(&current, false).unwrap();

    assert!(del_ifaces.kernel_ifaces["mlx5_ib2.8001"].is_absent());
    assert!(del_ifaces.kernel_ifaces["mlx5_ib2"].is_absent());
}

#[test]
fn test_ib_support_integer_pkey() {
    let iface: InfiniBandInterface = serde_yaml::from_str(
        r#"---
name: mlx5_ib2.8001
type: infiniband
state: up
infiniband:
  pkey: 32769
  mode: "connected"
  base-iface: "mlx5_ib2"
"#,
    )
    .unwrap();
    assert_eq!(iface.ib.unwrap().pkey, Some(0x8001));
}

#[test]
fn test_ib_support_string_pkey() {
    let iface: InfiniBandInterface = serde_yaml::from_str(
        r#"---
name: mlx5_ib2.8001
type: infiniband
state: up
infiniband:
  pkey: "32769"
  mode: "connected"
  base-iface: "mlx5_ib2"
"#,
    )
    .unwrap();
    assert_eq!(iface.ib.unwrap().pkey, Some(0x8001));
}

#[test]
fn test_ib_support_hex_string_pkey() {
    let iface: InfiniBandInterface = serde_yaml::from_str(
        r#"---
name: mlx5_ib2.8001
type: infiniband
state: up
infiniband:
  pkey: "0x8001"
  mode: "connected"
  base-iface: "mlx5_ib2"
"#,
    )
    .unwrap();
    assert_eq!(iface.ib.unwrap().pkey, Some(0x8001));
}

#[test]
fn test_ib_port_of_bridge_in_desire() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
- name: br0
  type: linux-bridge
  bridge:
    port:
    - name: mlx5_ib2
"#,
    )
    .unwrap();

    let result =
        crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
            &desire,
            &Interfaces::default(),
        );
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_ib_port_of_bridge_in_current() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: br0
  type: linux-bridge
  bridge:
    port:
    - name: mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
"#,
    )
    .unwrap();

    let result =
        crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
            &desire, &current,
        );
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_ib_port_of_bond_mode_in_desire() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: balance-rr
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
"#,
    )
    .unwrap();

    let result =
        crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
            &desire, &current,
        );
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_ib_port_of_bond_mode_in_current() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: bond0
  type: bond
  state: up
  link-aggregation:
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: balance-rr
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let result =
        crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
            &desire, &current,
        );
    assert!(result.is_err());
    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_ib_port_of_active_backup_bond_mode_in_current() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: bond0
  type: bond
  state: up
  link-aggregation:
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: active-backup
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
        &desire, &current,
    )
    .unwrap();
}

#[test]
fn test_ib_port_of_active_backup_bond_mode_in_both() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: active-backup
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: balance-rr
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
        &desire, &current,
    )
    .unwrap();
}

#[test]
fn test_ib_port_of_active_backup_bond_mode_in_desire() {
    let desire: Interfaces = serde_yaml::from_str(
        r#"---
- name: bond0
  type: bond
  state: up
  link-aggregation:
    mode: active-backup
    port:
    - mlx5_ib2
"#,
    )
    .unwrap();

    let current: Interfaces = serde_yaml::from_str(
        r#"---
- name: mlx5_ib2
  type: infiniband
  state: up
  infiniband:
    pkey: "0xffff"
    mode: "connected"
"#,
    )
    .unwrap();

    crate::ifaces::inter_ifaces_controller::check_infiniband_as_ports(
        &desire, &current,
    )
    .unwrap();
}
