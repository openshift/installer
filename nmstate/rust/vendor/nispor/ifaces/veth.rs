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

use std::collections::HashMap;

use rtnetlink::Handle;
use serde::{Deserialize, Serialize};

use crate::{Iface, IfaceType, NisporError};

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct VethInfo {
    // Interface name of peer.
    // Use interface index number when peer interface is in other namespace.
    pub peer: String,
}

pub type VethConf = VethInfo;

impl VethConf {
    pub(crate) async fn create(
        &self,
        handle: &Handle,
        name: &str,
    ) -> Result<(), NisporError> {
        match handle
            .link()
            .add()
            .veth(name.to_string(), self.peer.clone())
            .execute()
            .await
        {
            Ok(_) => Ok(()),
            Err(e) => Err(NisporError::bug(format!(
                "Failed to create new veth pair '{}' '{}': {}",
                &name, &self.peer, e
            ))),
        }
    }
}

pub(crate) fn veth_iface_tidy_up(iface_states: &mut HashMap<String, Iface>) {
    let mut index_to_name = HashMap::new();
    for iface in iface_states.values() {
        index_to_name.insert(format!("{}", iface.index), iface.name.clone());
    }

    for iface in iface_states.values_mut() {
        if iface.iface_type != IfaceType::Veth {
            continue;
        }

        if let Some(VethInfo { peer }) = &iface.veth {
            if let Some(peer_iface_name) = index_to_name.get(peer) {
                iface.veth = Some(VethInfo {
                    peer: peer_iface_name.clone(),
                })
            }
        }
    }
}
