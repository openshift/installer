use crate::NetworkState;

#[test]
fn test_dns_ignore_dns_purge_on_absent_iface() {
    let desired: NetworkState = serde_yaml::from_str(
        r#"---
dns-resolver:
  config:
    server: []
interfaces:
  - name: dummy0
    type: dummy
    state: absent
"#,
    )
    .unwrap();
    let current: NetworkState = serde_yaml::from_str(
        r#"---
dns-resolver:
  config:
    search:
    - example.com
    - example.org
    server:
    - 8.8.8.8
    - 2001:4860:4860::8888
interfaces:
  - name: dummy0
    type: dummy
    state: up
    ipv4:
      enabled: true
      dhcp: true
      auto-dns: false
    ipv6:
      enabled: true
      dhcp: true
      autoconf: true
      auto-dns: false
"#,
    )
    .unwrap();
    let (_, chg_state, del_state) =
        desired.gen_state_for_apply(&current).unwrap();

    assert!(chg_state.interfaces.to_vec().is_empty());
    let iface = del_state.interfaces.to_vec()[0];
    assert_eq!(iface.name(), "dummy0");
    assert!(iface.is_absent());
}
