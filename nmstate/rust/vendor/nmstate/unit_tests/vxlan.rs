use crate::VxlanInterface;

#[test]
fn test_vxlan_stringlized_attributes() {
    let iface: VxlanInterface = serde_yaml::from_str(
        r#"---
name: vxlan1
type: vxlan
state: up
vxlan:
  base-iface: "eth1"
  id: "101"
  destination-port: "3389"
"#,
    )
    .unwrap();
    let vxlan_conf = iface.vxlan.unwrap();

    assert_eq!(vxlan_conf.id, 101);
    assert_eq!(vxlan_conf.dst_port, Some(3389));
}
