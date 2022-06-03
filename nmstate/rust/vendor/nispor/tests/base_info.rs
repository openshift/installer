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

#[test]
fn test_iface_info_loopback() {
    let state = NetState::retrieve().unwrap();
    let iface = &state.ifaces["lo"];
    assert_eq!(iface.iface_type, nispor::IfaceType::Loopback);
    assert_eq!(iface.state, nispor::IfaceState::Unknown);
    assert_eq!(iface.mtu, 65536);
    assert_eq!(&iface.mac_address, "00:00:00:00:00:00");
    assert_eq!(iface.max_mtu, None);
    assert_eq!(iface.min_mtu, None);
    assert_eq!(
        iface.flags,
        &[
            nispor::IfaceFlags::Loopback,
            nispor::IfaceFlags::LowerUp,
            nispor::IfaceFlags::Running,
            nispor::IfaceFlags::Up,
        ]
    );
}
