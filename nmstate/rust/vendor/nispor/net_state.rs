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
use crate::ifaces::{get_ifaces, Iface};
use crate::route::get_routes;
use crate::route::Route;
use crate::route_rule::get_route_rules;
use crate::route_rule::RouteRule;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::runtime;

#[derive(Serialize, Deserialize, Debug, PartialEq, Clone)]
#[non_exhaustive]
pub struct NetState {
    pub ifaces: HashMap<String, Iface>,
    pub routes: Vec<Route>,
    pub rules: Vec<RouteRule>,
}

impl NetState {
    pub fn retrieve() -> Result<NetState, NisporError> {
        let rt = runtime::Builder::new_current_thread().enable_io().build()?;
        let ifaces = rt.block_on(get_ifaces())?;
        let routes = rt.block_on(get_routes(&ifaces))?;
        let rules = rt.block_on(get_route_rules())?;
        Ok(NetState {
            ifaces,
            routes,
            rules,
        })
    }

    // TODO: autoconvert NetState to NetConf and provide apply() here
}
