use crate::LldpConfig;

#[test]
fn test_lldp_stringlized_attributes() {
    let confs: Vec<LldpConfig> = serde_yaml::from_str(
        r#"
- enabled: "true"
- enabled: true
- enabled: "yes"
- enabled: "y"
- enabled: 1
- enabled: "1"
"#,
    )
    .unwrap();
    for conf in &confs {
        assert!(conf.enabled);
    }
    let confs: Vec<LldpConfig> = serde_yaml::from_str(
        r#"
- enabled: "false"
- enabled: false
- enabled: "no"
- enabled: "n"
- enabled: 0
- enabled: "0"
"#,
    )
    .unwrap();
    for conf in &confs {
        assert!(!conf.enabled);
    }
}
