use crate::BaseInterface;

#[test]
fn test_base_iface_stringlized_attributes() {
    let iface: BaseInterface = serde_yaml::from_str(
        r#"
name: "eth1"
mtu: "1500"
accept-all-mac-addresses: "true"
"#,
    )
    .unwrap();
    assert_eq!(iface.accept_all_mac_addresses, Some(true));
}

#[test]
fn test_base_iface_mac_address_uppercase_before_verification() {
    let mut iface: BaseInterface = serde_yaml::from_str(
        r#"
name: "eth1"
mtu: "1500"
mac-address: "d4:ee:07:25:42:5a"
"#,
    )
    .unwrap();
    iface.pre_verify_cleanup();
    assert_eq!(iface.mac_address, Some(String::from("D4:EE:07:25:42:5A")));
}
