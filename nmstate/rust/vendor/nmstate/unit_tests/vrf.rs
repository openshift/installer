use crate::VrfInterface;

#[test]
fn test_vrf_stringlized_attributes() {
    let iface: VrfInterface = serde_yaml::from_str(
        r#"---
name: vrf1
type: vrf
state: up
vrf:
  route-table-id: "101"
"#,
    )
    .unwrap();

    assert_eq!(iface.vrf.unwrap().table_id, 101);
}
