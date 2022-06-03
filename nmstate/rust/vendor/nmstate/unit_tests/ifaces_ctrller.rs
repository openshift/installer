use crate::{
    ifaces::inter_ifaces_controller::check_overbook_ports,
    unit_tests::testlib::{
        bond_with_ports, bridge_with_ports, new_br_iface, new_eth_iface,
        new_nested_4_ifaces, new_ovs_br_iface, new_ovs_iface,
    },
    ErrorKind, Interface, InterfaceState, InterfaceType, Interfaces,
    OvsBridgeInterface,
};

#[test]
fn test_ifaces_up_order_no_ctrler_reserse_order() {
    let mut ifaces = Interfaces::new();
    ifaces.push(new_eth_iface("eth2"));
    ifaces.push(new_eth_iface("eth1"));

    let (add_ifaces, _, _) = ifaces
        .gen_state_for_apply(&Interfaces::new(), false)
        .unwrap();

    assert_eq!(ifaces.kernel_ifaces["eth1"].base_iface().up_priority, 0);
    assert_eq!(ifaces.kernel_ifaces["eth2"].base_iface().up_priority, 0);

    let ordered_ifaces = add_ifaces.to_vec();
    assert_eq!(ordered_ifaces[0].name(), "eth1".to_string());
    assert_eq!(ordered_ifaces[1].name(), "eth2".to_string());
}

#[test]
fn test_ifaces_up_order_nested_4_depth_worst_case() {
    let mut ifaces = Interfaces::new();

    let [br0, br1, br2, br3, p1, p2] = new_nested_4_ifaces();

    // Push with reverse order which is the worst case
    ifaces.push(p2);
    ifaces.push(p1);
    ifaces.push(br3);
    ifaces.push(br2);
    ifaces.push(br1);
    ifaces.push(br0);

    let (add_ifaces, _, _) = ifaces
        .gen_state_for_apply(&Interfaces::new(), false)
        .unwrap();

    assert_eq!(ifaces.kernel_ifaces["br0"].base_iface().up_priority, 0);
    assert_eq!(ifaces.kernel_ifaces["br1"].base_iface().up_priority, 1);
    assert_eq!(ifaces.kernel_ifaces["br2"].base_iface().up_priority, 2);
    assert_eq!(ifaces.kernel_ifaces["br3"].base_iface().up_priority, 3);
    assert_eq!(ifaces.kernel_ifaces["p1"].base_iface().up_priority, 4);
    assert_eq!(ifaces.kernel_ifaces["p2"].base_iface().up_priority, 4);

    let ordered_ifaces = add_ifaces.to_vec();

    assert_eq!(ordered_ifaces[0].name(), "br0".to_string());
    assert_eq!(ordered_ifaces[1].name(), "br1".to_string());
    assert_eq!(ordered_ifaces[2].name(), "br2".to_string());
    assert_eq!(ordered_ifaces[3].name(), "br3".to_string());
    assert_eq!(ordered_ifaces[4].name(), "p1".to_string());
    assert_eq!(ordered_ifaces[5].name(), "p2".to_string());
}

#[test]
fn test_ifaces_up_order_nested_5_depth_worst_case() {
    let mut ifaces = Interfaces::new();
    let [_, br1, br2, br3, p1, p2] = new_nested_4_ifaces();

    let br4 = new_br_iface("br4");
    let mut br0 = new_br_iface("br0");

    br0.base_iface_mut().controller = Some("br4".to_string());
    br0.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);

    // Push with reverse order which is the worst case
    ifaces.push(p1);
    ifaces.push(p2);
    ifaces.push(br3);
    ifaces.push(br2);
    ifaces.push(br1);
    ifaces.push(br0);
    ifaces.push(br4);

    let result = ifaces.gen_state_for_apply(&Interfaces::new(), false);
    assert!(result.is_err());

    if let Err(e) = result {
        assert_eq!(e.kind(), ErrorKind::InvalidArgument);
    }
}

#[test]
fn test_ifaces_up_order_nested_5_depth_good_case() {
    let mut ifaces = Interfaces::new();
    let [_, br1, br2, br3, p1, p2] = new_nested_4_ifaces();

    let br4 = new_br_iface("br4");
    let mut br0 = new_br_iface("br0");

    br0.base_iface_mut().controller = Some("br4".to_string());
    br0.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);

    ifaces.push(br4);
    ifaces.push(br0);
    ifaces.push(br1);
    ifaces.push(br2);
    ifaces.push(br3);
    ifaces.push(p2);
    ifaces.push(p1);

    let (add_ifaces, _, _) = ifaces
        .gen_state_for_apply(&Interfaces::new(), false)
        .unwrap();

    assert_eq!(ifaces.kernel_ifaces["br4"].base_iface().up_priority, 0);
    assert_eq!(ifaces.kernel_ifaces["br0"].base_iface().up_priority, 1);
    assert_eq!(ifaces.kernel_ifaces["br1"].base_iface().up_priority, 2);
    assert_eq!(ifaces.kernel_ifaces["br2"].base_iface().up_priority, 3);
    assert_eq!(ifaces.kernel_ifaces["br3"].base_iface().up_priority, 4);
    assert_eq!(ifaces.kernel_ifaces["p1"].base_iface().up_priority, 5);
    assert_eq!(ifaces.kernel_ifaces["p2"].base_iface().up_priority, 5);

    let ordered_ifaces = add_ifaces.to_vec();

    assert_eq!(ordered_ifaces[0].name(), "br4".to_string());
    assert_eq!(ordered_ifaces[1].name(), "br0".to_string());
    assert_eq!(ordered_ifaces[2].name(), "br1".to_string());
    assert_eq!(ordered_ifaces[3].name(), "br2".to_string());
    assert_eq!(ordered_ifaces[4].name(), "br3".to_string());
    assert_eq!(ordered_ifaces[5].name(), "p1".to_string());
    assert_eq!(ordered_ifaces[6].name(), "p2".to_string());
}

#[test]
fn test_auto_include_ovs_interface() {
    let mut ifaces = Interfaces::new();
    ifaces.push(new_ovs_br_iface("br0", &["p1", "p2"]));

    let (add_ifaces, _, _) = ifaces
        .gen_state_for_apply(&Interfaces::new(), false)
        .unwrap();

    println!("{:?}", ifaces);

    assert_eq!(ifaces.kernel_ifaces["p1"].base_iface().up_priority, 1);
    assert_eq!(ifaces.kernel_ifaces["p1"].base_iface().name, "p1");
    assert_eq!(
        ifaces.kernel_ifaces["p1"].base_iface().iface_type,
        InterfaceType::OvsInterface
    );
    assert_eq!(
        ifaces.kernel_ifaces["p1"].base_iface().controller,
        Some("br0".to_string())
    );
    assert_eq!(
        ifaces.kernel_ifaces["p1"].base_iface().controller_type,
        Some(InterfaceType::OvsBridge)
    );
    assert_eq!(ifaces.kernel_ifaces["p2"].base_iface().up_priority, 1);
    assert_eq!(ifaces.kernel_ifaces["p2"].base_iface().name, "p2");
    assert_eq!(
        ifaces.kernel_ifaces["p2"].base_iface().iface_type,
        InterfaceType::OvsInterface
    );
    assert_eq!(
        ifaces.kernel_ifaces["p2"].base_iface().controller,
        Some("br0".to_string())
    );
    assert_eq!(
        ifaces.kernel_ifaces["p2"].base_iface().controller_type,
        Some(InterfaceType::OvsBridge)
    );
    assert_eq!(
        ifaces.user_ifaces[&("br0".to_string(), InterfaceType::OvsBridge)]
            .base_iface()
            .up_priority,
        0
    );

    let ordered_ifaces = add_ifaces.to_vec();

    assert_eq!(ordered_ifaces[0].name(), "br0".to_string());
    assert_eq!(ordered_ifaces[1].name(), "p1".to_string());
    assert_eq!(ordered_ifaces[2].name(), "p2".to_string());
}

#[test]
fn test_auto_absent_ovs_interface() {
    let mut cur_ifaces = Interfaces::new();
    cur_ifaces.push(new_ovs_br_iface("br0", &["p1", "p2"]));
    cur_ifaces.push(new_ovs_iface("p1", "br0"));
    cur_ifaces.push(new_ovs_iface("p2", "br0"));

    let mut absent_br0 = OvsBridgeInterface::new();
    absent_br0.base.name = "br0".to_string();
    absent_br0.base.state = InterfaceState::Absent;
    let mut ifaces = Interfaces::new();
    ifaces.push(Interface::OvsBridge(absent_br0));

    let (_, _, del_ifaces) =
        ifaces.gen_state_for_apply(&cur_ifaces, false).unwrap();

    println!("{:?}", ifaces);

    assert_eq!(ifaces.kernel_ifaces["p1"].base_iface().name, "p1");
    assert_eq!(
        ifaces.kernel_ifaces["p1"].base_iface().iface_type,
        InterfaceType::OvsInterface
    );
    assert_eq!(
        ifaces.kernel_ifaces["p1"].base_iface().state,
        InterfaceState::Absent
    );

    assert_eq!(ifaces.kernel_ifaces["p2"].base_iface().name, "p2");
    assert_eq!(
        ifaces.kernel_ifaces["p2"].base_iface().iface_type,
        InterfaceType::OvsInterface
    );
    assert_eq!(
        ifaces.kernel_ifaces["p2"].base_iface().state,
        InterfaceState::Absent
    );
    assert_eq!(
        ifaces.user_ifaces[&("br0".to_string(), InterfaceType::OvsBridge)]
            .base_iface()
            .state,
        InterfaceState::Absent
    );

    let del_ifaces = del_ifaces.to_vec();
    assert_eq!(del_ifaces[0].base_iface().name, "br0");
    assert_eq!(del_ifaces[1].base_iface().name, "p1");
    assert_eq!(del_ifaces[2].base_iface().name, "p2");
}

#[test]
fn test_overbook_port_in_single_bridge() {
    let mut desired = Interfaces::new();

    desired.push(bridge_with_ports("br0", &["eth0"]));
    desired.push(new_eth_iface("eth0"));

    let current = Interfaces::new();

    assert!(check_overbook_ports(&desired, &current).is_ok());
}

#[test]
fn test_overbook_port_in_two_bridges() {
    let mut desired = Interfaces::new();

    desired.push(bridge_with_ports("br0", &["eth0"]));
    desired.push(bridge_with_ports("br1", &["eth0"]));

    let mut current = Interfaces::new();
    current.push(new_eth_iface("eth0"));

    let result = check_overbook_ports(&desired, &current);
    assert!(result.is_err());
    assert_eq!(result.err().unwrap().kind(), ErrorKind::InvalidArgument);
}

#[test]
fn test_overbook_port_moves_between_bridges() {
    let mut current = Interfaces::new();
    current.push(bridge_with_ports("br0", &["eth0"]));
    let mut eth0 = new_eth_iface("eth0");
    eth0.base_iface_mut().controller = Some("br0".to_string());
    eth0.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    current.push(eth0);

    let mut desired = Interfaces::new();
    desired.push(bridge_with_ports("br0", &[]));
    desired.push(bridge_with_ports("br1", &["eth0"]));

    assert!(check_overbook_ports(&desired, &current).is_ok());
}

#[test]
fn test_overbook_current_bridge_is_deleted() {
    let mut current = Interfaces::new();
    current.push(bridge_with_ports("br0", &["eth0"]));

    let mut desired = Interfaces::new();
    desired.push(bridge_with_ports("br1", &["eth0"]));
    let mut absent_iface = new_br_iface("br0");
    absent_iface.base_iface_mut().state = InterfaceState::Absent;
    desired.push(absent_iface);

    assert!(check_overbook_ports(&desired, &current).is_ok());
}

#[test]
fn test_overbook_port_used_in_current_bridge() {
    let mut current = Interfaces::new();
    current.push(bridge_with_ports("br0", &["eth0"]));
    let mut eth0 = new_eth_iface("eth0");
    eth0.base_iface_mut().controller = Some("br0".to_string());
    eth0.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    current.push(eth0);

    let mut desired = Interfaces::new();
    desired.push(bridge_with_ports("br1", &["eth0"]));

    let result = check_overbook_ports(&desired, &current);
    assert!(result.is_err());
    assert_eq!(result.err().unwrap().kind(), ErrorKind::InvalidArgument);
}

#[test]
fn test_overbook_port_used_in_current_bond() {
    let mut current = Interfaces::new();
    current.push(bond_with_ports("bond0", &["eth0"]));
    let mut eth0 = new_eth_iface("eth0");
    eth0.base_iface_mut().controller = Some("bond0".to_string());
    eth0.base_iface_mut().controller_type = Some(InterfaceType::Bond);
    current.push(eth0);

    let mut desired = Interfaces::new();
    desired.push(bond_with_ports("bond1", &["eth0"]));

    let result = check_overbook_ports(&desired, &current);
    assert!(result.is_err());
    assert_eq!(result.err().unwrap().kind(), ErrorKind::InvalidArgument);
}

#[test]
fn test_overbook_swap_port_of_bond() {
    let mut current = Interfaces::new();
    current.push(bond_with_ports("bond0", &["eth0"]));
    current.push(bond_with_ports("bond1", &["eth1"]));
    let mut eth0 = new_eth_iface("eth0");
    eth0.base_iface_mut().controller = Some("bond0".to_string());
    eth0.base_iface_mut().controller_type = Some(InterfaceType::Bond);
    current.push(eth0);
    let mut eth1 = new_eth_iface("eth1");
    eth1.base_iface_mut().controller = Some("bond1".to_string());
    eth1.base_iface_mut().controller_type = Some(InterfaceType::Bond);
    current.push(eth1);

    let mut desired = Interfaces::new();
    desired.push(bond_with_ports("bond1", &["eth0"]));
    desired.push(bond_with_ports("bond0", &["eth1"]));

    assert!(check_overbook_ports(&desired, &current).is_ok());
}
