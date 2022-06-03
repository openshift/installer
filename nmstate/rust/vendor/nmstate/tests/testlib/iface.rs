use crate::testlib::fs::read_folder;
use std::path::Path;

pub(crate) fn assert_iface_missing(iface_name: &str) {
    assert!(!is_iface_exits(iface_name));
}

pub(crate) fn assert_iface_exists(iface_name: &str) {
    assert!(is_iface_exits(iface_name));
}

fn is_iface_exits(iface_name: &str) -> bool {
    Path::new(&format!("/sys/class/net/{}", iface_name)).exists()
}

pub(crate) fn assert_iface_bridge(iface_name: &str, ports: &[&str]) {
    assert!(
        Path::new(&format!("/sys/class/net/{}/bridge/", iface_name)).exists()
    );
    let cur_ports = read_folder(&format!("/sys/class/net/{}/brif", iface_name));

    println!("Current bridge ports {:?}", cur_ports);
    println!("Desired bridge ports {:?}", ports);

    assert_eq!(ports.len(), cur_ports.len());

    for port in ports {
        assert!(cur_ports.contains(&port.to_string()));
    }
}
