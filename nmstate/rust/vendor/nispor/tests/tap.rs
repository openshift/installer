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

use nispor::NetState;
use pretty_assertions::assert_eq;
use serde_yaml;
use std::panic;

mod utils;

const IFACE_NAME: &str = "tap1";

const EXPECTED_TAP_INFO: &str = r#"---
mode: tap
owner: 1001
group: 0
pi: false
vnet_hdr: true
multi_queue: true
persist: true
num_queues: 0
num_disabled_queues: 0"#;

#[test]
fn test_get_tap_iface_yaml() {
    with_tap_iface(|| {
        let state = NetState::retrieve().unwrap();
        let iface = &state.ifaces[IFACE_NAME];
        assert_eq!(iface.iface_type, nispor::IfaceType::Tun);
        assert_eq!(
            serde_yaml::to_string(&iface.tun).unwrap().trim(),
            EXPECTED_TAP_INFO
        );
    });
}

fn with_tap_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("tap");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
