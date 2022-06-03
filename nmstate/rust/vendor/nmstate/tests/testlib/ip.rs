use crate::testlib::cmd::cmd_exec_check;
use serde::Deserialize;

#[derive(Deserialize, Debug, PartialEq, Eq, Clone)]
struct IpAddrShowResult {
    ifname: String,
    addr_info: Vec<IpAddrShowResultAddrInfo>,
}

#[derive(Deserialize, Debug, PartialEq, Eq, Clone)]
struct IpAddrShowResultAddrInfo {
    local: String,
    prefixlen: u32,
}

pub(crate) fn assert_ip(iface_name: &str, ip_addrs: &[&str]) {
    let stdout = cmd_exec_check(&["ip", "-j", "addr", "show", iface_name]);
    let info: Vec<IpAddrShowResult> =
        serde_json::from_str(&stdout).expect("Failed to parse ip addr output");
    let mut cur_ip_addrs = Vec::new();
    for addr in &info[0].addr_info {
        cur_ip_addrs.push(format!("{}/{}", addr.local, addr.prefixlen));
    }

    println!("Current IP addresses {:?}", cur_ip_addrs);
    println!("Desired IP addresses {:?}", ip_addrs);
    for desire_ip_addr in ip_addrs {
        assert!(cur_ip_addrs.contains(&desire_ip_addr.to_string()));
    }
}
