use crate::MacVtapInterface;

#[test]
fn test_mac_vtap_stringlized_attributes() {
    let iface: MacVtapInterface = serde_yaml::from_str(
        r#"---
name: mac1
type: mac-vtap
state: up
mac-vtap:
  base-iface: "eth1"
  mode: "vepa"
  accept-all-mac: "true"
"#,
    )
    .unwrap();

    let mac_conf = iface.mac_vtap.unwrap();
    assert_eq!(mac_conf.accept_all_mac, Some(true));
}
