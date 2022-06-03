use crate::{
    BondConfig, BondInterface, BondMode, EthernetInterface, Interface,
    InterfaceType, LinuxBridgeConfig, LinuxBridgeInterface,
    LinuxBridgePortConfig, OvsBridgeConfig, OvsBridgeInterface,
    OvsBridgePortConfig, OvsInterface, UnknownInterface, VlanConfig,
    VlanInterface,
};

pub(crate) fn new_eth_iface(name: &str) -> Interface {
    let mut iface = EthernetInterface::new();
    iface.base.name = name.to_string();
    Interface::Ethernet(iface)
}

pub(crate) fn new_unknown_iface(name: &str) -> Interface {
    let mut iface = UnknownInterface::new();
    iface.base.name = name.to_string();
    Interface::Unknown(iface)
}

pub(crate) fn new_br_iface(name: &str) -> Interface {
    let mut iface = LinuxBridgeInterface::new();
    iface.base.name = name.to_string();
    Interface::LinuxBridge(iface)
}

fn new_bond_iface(name: &str) -> Interface {
    let mut iface = BondInterface::new();
    iface.base.name = name.to_string();
    Interface::Bond(iface)
}

pub(crate) fn new_ovs_br_iface(name: &str, port_names: &[&str]) -> Interface {
    let mut br0 = OvsBridgeInterface::new();
    br0.base.iface_type = InterfaceType::OvsBridge;
    br0.base.name = name.to_string();
    let mut br_conf = OvsBridgeConfig::new();
    let mut br_port_confs = Vec::new();
    for port_name in port_names {
        let mut br_port_conf = OvsBridgePortConfig::new();
        br_port_conf.name = port_name.to_string();
        br_port_confs.push(br_port_conf);
    }
    br_conf.ports = Some(br_port_confs);
    br0.bridge = Some(br_conf);
    Interface::OvsBridge(br0)
}

pub(crate) fn new_ovs_iface(name: &str, ctrl_name: &str) -> Interface {
    let mut iface = OvsInterface::new();
    iface.base.iface_type = InterfaceType::OvsInterface;
    iface.base.name = name.to_string();
    iface.base.controller = Some(ctrl_name.to_string());
    iface.base.controller_type = Some(InterfaceType::OvsBridge);
    Interface::OvsInterface(iface)
}

pub(crate) fn new_vlan_iface(name: &str, parent: &str, id: u16) -> Interface {
    let mut iface = VlanInterface::new();
    iface.base.name = name.to_string();
    iface.base.iface_type = InterfaceType::Vlan;
    iface.vlan = Some(VlanConfig {
        base_iface: parent.to_string(),
        id,
    });
    Interface::Vlan(iface)
}

pub(crate) fn new_nested_4_ifaces() -> [Interface; 6] {
    let br0 = new_br_iface("br0");
    let mut br1 = new_br_iface("br1");
    let mut br2 = new_br_iface("br2");
    let mut br3 = new_br_iface("br3");
    let mut p1 = new_eth_iface("p1");
    let mut p2 = new_eth_iface("p2");

    br1.base_iface_mut().controller = Some("br0".to_string());
    br1.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    br2.base_iface_mut().controller = Some("br1".to_string());
    br2.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    br3.base_iface_mut().controller = Some("br2".to_string());
    br3.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    p1.base_iface_mut().controller = Some("br3".to_string());
    p1.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);
    p2.base_iface_mut().controller = Some("br3".to_string());
    p2.base_iface_mut().controller_type = Some(InterfaceType::LinuxBridge);

    // Place the ifaces in mixed order to complex the work
    [br0, br1, br2, br3, p1, p2]
}

pub(crate) fn bridge_with_ports(name: &str, ports: &[&str]) -> Interface {
    let ports = ports
        .iter()
        .map(|port| LinuxBridgePortConfig {
            name: port.to_string(),
            ..Default::default()
        })
        .collect::<Vec<_>>();

    let mut br0 = new_br_iface(name);
    if let Interface::LinuxBridge(br) = &mut br0 {
        br.bridge = Some(LinuxBridgeConfig {
            port: Some(ports),
            ..Default::default()
        })
    };
    br0
}

pub(crate) fn bond_with_ports(name: &str, ports: &[&str]) -> Interface {
    let ports = ports.iter().map(|p| p.to_string()).collect::<Vec<String>>();
    let mut iface = new_bond_iface(name);
    if let Interface::Bond(bond_iface) = &mut iface {
        bond_iface.bond = Some(BondConfig {
            mode: Some(BondMode::RoundRobin),
            port: Some(ports),
            ..Default::default()
        });
    }
    iface
}
