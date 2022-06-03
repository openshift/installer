// Copyright 2021 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use nispor::{NetConf, NetState};
use pretty_assertions::assert_eq;
use serde_yaml;
use std::panic;

mod utils;

const IFACE_NAME: &str = "veth1";

fn with_veth_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("veth");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}

const ADD_IP_CONF: &str = r#"---
ifaces:
  - name: veth1
    ipv4:
      addresses:
        - address: "192.0.2.1"
          prefix_len: 24
    ipv6:
      addresses:
        - address: "2001:db8:a::9"
          prefix_len: 64"#;

const ADD_IP_CONF_DYNAMIC: &str = r#"---
ifaces:
  - name: veth1
    ipv4:
      addresses:
        - address: "192.0.2.1"
          prefix_len: 24
          valid_lft: 120sec
          preferred_lft: 60sec
    ipv6:
      addresses:
        - address: "2001:db8:a::9"
          prefix_len: 64
          valid_lft: 121sec
          preferred_lft: 61sec"#;

const DEL_IP_CONF: &str = r#"---
ifaces:
  - name: veth1
    ipv4:
      addresses:
        - address: "192.0.2.1"
          prefix_len: 24
          remove: true
    ipv6:
      addresses:
        - address: "2001:db8:a::9"
          prefix_len: 64
          remove: true"#;

const EXPECTED_IPV4_INFO: &str = r#"---
addresses:
  - address: 192.0.2.1
    prefix_len: 24
    valid_lft: forever
    preferred_lft: forever"#;

const EXPECTED_IPV4_DYNAMIC_INFO: &str = r#"---
addresses:
  - address: 192.0.2.1
    prefix_len: 24
    valid_lft: 120sec
    preferred_lft: 60sec"#;

const EXPECTED_IPV6_INFO: &str = r#"---
addresses:
  - address: "2001:db8:a::9"
    prefix_len: 64
    valid_lft: forever
    preferred_lft: forever
  - address: "fe80::223:45ff:fe67:891a"
    prefix_len: 64
    valid_lft: forever
    preferred_lft: forever"#;

const EXPECTED_IPV6_DYNAMIC_INFO: &str = r#"---
addresses:
  - address: "2001:db8:a::9"
    prefix_len: 64
    valid_lft: 121sec
    preferred_lft: 61sec
  - address: "fe80::223:45ff:fe67:891a"
    prefix_len: 64
    valid_lft: forever
    preferred_lft: forever"#;

const EXPECTED_EMPTY_IPV6_INFO: &str = r#"---
addresses:
  - address: "fe80::223:45ff:fe67:891a"
    prefix_len: 64
    valid_lft: forever
    preferred_lft: forever"#;

#[test]
fn test_add_and_remove_ip() {
    with_veth_iface(|| {
        let conf: NetConf = serde_yaml::from_str(ADD_IP_CONF).unwrap();
        conf.apply().unwrap();
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        let iface_type = &iface.iface_type;
        assert_eq!(iface_type, &nispor::IfaceType::Veth);
        assert_eq!(
            serde_yaml::to_string(&iface.ipv4).unwrap().trim(),
            EXPECTED_IPV4_INFO
        );
        assert_eq!(
            serde_yaml::to_string(&iface.ipv6).unwrap().trim(),
            EXPECTED_IPV6_INFO
        );
        let conf: NetConf = serde_yaml::from_str(DEL_IP_CONF).unwrap();
        conf.apply().unwrap();
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        let iface_type = &iface.iface_type;
        assert_eq!(iface_type, &nispor::IfaceType::Veth);
        assert_eq!(iface.ipv4, None);
        assert_eq!(
            serde_yaml::to_string(&iface.ipv6).unwrap().trim(),
            EXPECTED_EMPTY_IPV6_INFO
        );
    });
}

#[test]
fn test_add_and_remove_dynamic_ip() {
    with_veth_iface(|| {
        let conf: NetConf = serde_yaml::from_str(ADD_IP_CONF_DYNAMIC).unwrap();
        conf.apply().unwrap();
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        let iface_type = &iface.iface_type;
        assert_eq!(iface_type, &nispor::IfaceType::Veth);
        assert_eq!(
            serde_yaml::to_string(&iface.ipv4).unwrap().trim(),
            EXPECTED_IPV4_DYNAMIC_INFO
        );
        assert_eq!(
            serde_yaml::to_string(&iface.ipv6).unwrap().trim(),
            EXPECTED_IPV6_DYNAMIC_INFO
        );
        let conf: NetConf = serde_yaml::from_str(DEL_IP_CONF).unwrap();
        conf.apply().unwrap();
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        let iface_type = &iface.iface_type;
        assert_eq!(iface_type, &nispor::IfaceType::Veth);
        assert_eq!(iface.ipv4, None);
        assert_eq!(
            serde_yaml::to_string(&iface.ipv6).unwrap().trim(),
            EXPECTED_EMPTY_IPV6_INFO
        );
    });
}

#[test]
fn test_add_dynamic_ip_repeat() {
    with_veth_iface(|| {
        let conf: NetConf = serde_yaml::from_str(ADD_IP_CONF_DYNAMIC).unwrap();
        conf.apply().unwrap();
        conf.apply().unwrap();
        std::thread::sleep(std::time::Duration::from_secs(2));
        conf.apply().unwrap();
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        let iface_type = &iface.iface_type;
        assert_eq!(iface_type, &nispor::IfaceType::Veth);
        assert_eq!(
            serde_yaml::to_string(&iface.ipv4).unwrap().trim(),
            EXPECTED_IPV4_DYNAMIC_INFO
        );
        assert_eq!(
            serde_yaml::to_string(&iface.ipv6).unwrap().trim(),
            EXPECTED_IPV6_DYNAMIC_INFO
        );
    });
}
