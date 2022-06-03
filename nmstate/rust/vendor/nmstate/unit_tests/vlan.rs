use crate::VlanInterface;

#[test]
fn test_vlan_stringlized_attributes() {
    let iface: VlanInterface = serde_yaml::from_str(
        r#"---
name: vlan1
type: vlan
state: up
vlan:
  base-iface: "eth1"
  id: "101"
"#,
    )
    .unwrap();

    assert_eq!(iface.vlan.unwrap().id, 101);
}
