use crate::{
    unit_tests::testlib::new_eth_iface, EthernetConfig, Interface, Interfaces,
    SrIovConfig, SrIovVfConfig,
};

#[test]
fn test_sriov_vf_mac_mix_case() {
    let mut cur_ifaces = Interfaces::new();
    let mut cur_iface = new_eth_iface("eth1");
    if let Interface::Ethernet(ref mut eth_iface) = cur_iface {
        let mut eth_conf = EthernetConfig::new();
        let mut sriov_conf = SrIovConfig::new();
        let mut vf_conf = SrIovVfConfig::new();
        vf_conf.id = 0;
        vf_conf.mac_address = Some("00:11:22:33:44:FF".into());
        sriov_conf.vfs = Some(vec![vf_conf]);
        eth_conf.sr_iov = Some(sriov_conf);
        eth_iface.ethernet = Some(eth_conf);
    } else {
        panic!("Should be ethernet interface");
    }
    cur_ifaces.push(cur_iface);

    let mut des_ifaces = Interfaces::new();
    let mut des_iface = new_eth_iface("eth1");
    if let Interface::Ethernet(ref mut eth_iface) = des_iface {
        let mut eth_conf = EthernetConfig::new();
        let mut sriov_conf = SrIovConfig::new();
        let mut vf_conf = SrIovVfConfig::new();
        vf_conf.id = 0;
        vf_conf.mac_address = Some("00:11:22:33:44:Ff".into());
        sriov_conf.vfs = Some(vec![vf_conf]);
        eth_conf.sr_iov = Some(sriov_conf);
        eth_iface.ethernet = Some(eth_conf);
    } else {
        panic!("Should be ethernet interface");
    }
    des_ifaces.push(des_iface);

    des_ifaces.verify(&cur_ifaces).unwrap();
}
