use crate::nm::nm_dbus::{NmConnection, NmSettingConnection};
use crate::{
    nm::profile::use_uuid_for_controller_reference, Interface, InterfaceType,
    OvsBridgeBondConfig, OvsBridgeBondPortConfig, OvsBridgeConfig,
    OvsBridgeInterface, OvsBridgePortConfig,
};
use std::collections::HashMap;
const UUID1: &str = "8aca0200-accc-4d13-a62f-3c89a6da53c5";
const UUID2: &str = "1c646761-efcc-4d33-a0d9-cb3c1c2d3309";
const UUID3: &str = "06935474-b8d3-4e7c-be52-48e2e6e6b3b9";
const UUID4: &str = "3c80d8de-a6d7-47da-b0b3-47d2b1052fe5";

#[test]
fn test_use_uuid_for_controller_reference_with_ovs_bond() {
    let mut nm_conns: Vec<NmConnection> = Vec::new();
    let mut nm_conn = NmConnection::default();
    let mut nm_conn_set = NmSettingConnection::default();
    nm_conn_set.id = Some("ovs-br-br0".to_string());
    nm_conn_set.uuid = Some(UUID1.to_string());
    nm_conn_set.iface_type = Some("ovs-bridge".to_string());
    nm_conn_set.iface_name = Some("br0".to_string());
    nm_conn.connection = Some(nm_conn_set);
    nm_conns.push(nm_conn);

    let mut nm_conn = NmConnection::default();
    let mut nm_conn_set = NmSettingConnection::default();
    nm_conn_set.id = Some("ovs-port-bond1".to_string());
    nm_conn_set.uuid = Some(UUID2.to_string());
    nm_conn_set.iface_type = Some("ovs-port".to_string());
    nm_conn_set.iface_name = Some("bond1".to_string());
    nm_conn_set.controller = Some("br0".into());
    nm_conn_set.controller_type = Some("ovs-bridge".into());
    nm_conn.connection = Some(nm_conn_set);
    nm_conns.push(nm_conn);

    let mut nm_conn = NmConnection::default();
    let mut nm_conn_set = NmSettingConnection::default();
    nm_conn_set.id = Some("ovs-iface-p1".to_string());
    nm_conn_set.uuid = Some(UUID3.to_string());
    nm_conn_set.iface_type = Some("ovs-interface".to_string());
    nm_conn_set.iface_name = Some("p1".to_string());
    nm_conn_set.controller = Some("br0".into());
    nm_conn_set.controller_type = Some("ovs-port".into());
    nm_conn.connection = Some(nm_conn_set);
    nm_conns.push(nm_conn);

    let mut nm_conn = NmConnection::default();
    let mut nm_conn_set = NmSettingConnection::default();
    nm_conn_set.id = Some("ovs-iface-p2".to_string());
    nm_conn_set.uuid = Some(UUID4.to_string());
    nm_conn_set.iface_type = Some("ovs-interface".to_string());
    nm_conn_set.iface_name = Some("p2".to_string());
    nm_conn_set.controller = Some("br0".into());
    nm_conn_set.controller_type = Some("ovs-port".into());
    nm_conn.connection = Some(nm_conn_set);
    nm_conns.push(nm_conn);

    let mut br0 = OvsBridgeInterface::new();
    br0.base.iface_type = InterfaceType::OvsBridge;
    br0.base.name = "br0".to_string();
    let mut br_conf = OvsBridgeConfig::new();
    let mut p1_port_conf = OvsBridgeBondPortConfig::new();
    p1_port_conf.name = "p1".to_string();
    let mut p2_port_conf = OvsBridgeBondPortConfig::new();
    p2_port_conf.name = "p2".to_string();
    let mut bond_conf = OvsBridgeBondConfig::new();
    bond_conf.ports = Some(vec![p1_port_conf, p2_port_conf]);
    let mut port_conf = OvsBridgePortConfig::new();
    port_conf.name = "bond1".to_string();
    port_conf.bond = Some(bond_conf);
    br_conf.ports = Some(vec![port_conf]);
    br0.bridge = Some(br_conf);

    let mut user_ifaces: HashMap<(String, InterfaceType), Interface> =
        HashMap::new();

    user_ifaces.insert(
        ("br0".to_string(), InterfaceType::OvsBridge),
        Interface::OvsBridge(br0),
    );

    use_uuid_for_controller_reference(
        &mut nm_conns,
        &user_ifaces,
        &HashMap::new(),
        &[],
    )
    .unwrap();

    println!("nm_conns {:?}", nm_conns);

    let br0_nm_con_set = nm_conns[0].connection.as_ref().unwrap();
    println!("br0 {:?}", br0_nm_con_set);
    assert!(br0_nm_con_set.iface_name == Some("br0".to_string()));
    assert!(br0_nm_con_set.id == Some("ovs-br-br0".to_string()));

    let bond1_nm_con_set = nm_conns[1].connection.as_ref().unwrap();
    println!("bond1 {:?}", bond1_nm_con_set);
    assert!(bond1_nm_con_set.id == Some("ovs-port-bond1".to_string()));
    assert!(bond1_nm_con_set.iface_name == Some("bond1".to_string()));
    assert!(bond1_nm_con_set.iface_type == Some("ovs-port".to_string()));
    assert!(bond1_nm_con_set.controller == Some(UUID1.to_string()));
    assert!(bond1_nm_con_set.controller_type == Some("ovs-bridge".to_string()));

    let p1_nm_con_set = nm_conns[2].connection.as_ref().unwrap();
    println!("p1 {:?}", p1_nm_con_set);
    assert!(p1_nm_con_set.id == Some("ovs-iface-p1".to_string()));
    assert!(p1_nm_con_set.iface_name == Some("p1".to_string()));
    assert!(p1_nm_con_set.iface_type == Some("ovs-interface".to_string()));
    assert!(p1_nm_con_set.controller == Some(UUID2.to_string()));
    assert!(p1_nm_con_set.controller_type == Some("ovs-port".to_string()));

    let p2_nm_con_set = nm_conns[3].connection.as_ref().unwrap();
    println!("p2 {:?}", p2_nm_con_set);
    assert!(p2_nm_con_set.id == Some("ovs-iface-p2".to_string()));
    assert!(p2_nm_con_set.iface_name == Some("p2".to_string()));
    assert!(p2_nm_con_set.iface_type == Some("ovs-interface".to_string()));
    assert!(p2_nm_con_set.controller == Some(UUID2.to_string()));
    assert!(p2_nm_con_set.controller_type == Some("ovs-port".to_string()));
}
