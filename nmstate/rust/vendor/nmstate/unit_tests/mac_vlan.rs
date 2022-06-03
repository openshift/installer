use crate::MacVlanInterface;

#[test]
fn test_mac_vlan_stringlized_attributes() {
    let iface: MacVlanInterface = serde_yaml::from_str(
        r#"---
name: mac1
type: mac-vlan
state: up
mac-vlan:
  base-iface: "eth1"
  mode: "vepa"
  accept-all-mac: "true"
"#,
    )
    .unwrap();

    let mac_conf = iface.mac_vlan.unwrap();
    assert_eq!(mac_conf.accept_all_mac, Some(true));
}
