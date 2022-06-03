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

use crate::error::NisporError;
use crate::ifaces::{
    change_ifaces, create_ifaces, delete_ifaces, get_iface_name2index,
    get_ifaces, IfaceConf, IfaceState,
};
use crate::route::{apply_routes_conf, RouteConf};

use serde::{Deserialize, Serialize};
use tokio::runtime;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone, Default)]
#[non_exhaustive]
pub struct NetConf {
    pub ifaces: Option<Vec<IfaceConf>>,
    pub routes: Option<Vec<RouteConf>>,
}

impl NetConf {
    pub fn apply(&self) -> Result<(), NisporError> {
        let rt = runtime::Builder::new_current_thread().enable_io().build()?;
        if let Some(ref ifaces) = &self.ifaces {
            let cur_iface_name_2_index = rt.block_on(get_iface_name2index())?;
            let mut new_ifaces = Vec::new();
            let mut del_ifaces = Vec::new();
            let mut chg_ifaces = Vec::new();
            for iface in ifaces {
                let cur_iface_index = cur_iface_name_2_index.get(&iface.name);
                if iface.state == IfaceState::Absent {
                    if let Some(cur_iface_index) = cur_iface_index {
                        del_ifaces
                            .push((iface.name.as_str(), *cur_iface_index));
                    }
                } else if cur_iface_index == None {
                    new_ifaces.push(iface);
                    chg_ifaces.push(iface);
                } else {
                    chg_ifaces.push(iface);
                }
            }
            rt.block_on(delete_ifaces(&del_ifaces))?;
            rt.block_on(create_ifaces(&new_ifaces, &cur_iface_name_2_index))?;

            let cur_ifaces = rt.block_on(get_ifaces())?;
            rt.block_on(change_ifaces(&chg_ifaces, &cur_ifaces))?;
        }

        if let Some(routes) = self.routes.as_ref() {
            let cur_iface_name_2_index = rt.block_on(get_iface_name2index())?;
            rt.block_on(apply_routes_conf(routes, &cur_iface_name_2_index))?;
        }
        Ok(())
    }
}
