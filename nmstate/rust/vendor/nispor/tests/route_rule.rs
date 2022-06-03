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

const TEST_TABLE_ID: u32 = 100;

const EXPECTED_YAML_OUTPUT: &str = r#"---
- action: table
  address_family: ipv6
  flags: 0
  tos: 16
  table: 100
  dst: "2001:db8:f::253/128"
  src: "2001:db8:f::254/128"
  iif: eth1
  oif: eth2
  priority: 999
- action: table
  address_family: ipv4
  flags: 0
  tos: 16
  table: 100
  dst: 192.0.2.2/32
  src: 192.0.2.1/32
  iif: eth1
  oif: eth2
  priority: 999"#;

#[test]
fn test_get_route_rule_yaml() {
    with_route_rule_test_iface(|| {
        let state = NetState::retrieve().unwrap();
        let mut expected_rules = Vec::new();
        for mut rule in state.rules {
            if Some(TEST_TABLE_ID) == rule.table {
                // Travis CI Ubuntu 18.04 does not support protocol.
                rule.protocol = None;
                expected_rules.push(rule)
            }
        }
        assert_eq!(
            serde_yaml::to_string(&expected_rules).unwrap().trim(),
            EXPECTED_YAML_OUTPUT
        );
    });
}

fn with_route_rule_test_iface<T>(test: T) -> ()
where
    T: FnOnce() -> () + panic::UnwindSafe,
{
    utils::set_network_environment("rule");

    let result = panic::catch_unwind(|| {
        test();
    });

    utils::clear_network_environment();
    assert!(result.is_ok())
}
