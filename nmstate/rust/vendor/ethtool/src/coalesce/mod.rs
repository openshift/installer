// SPDX-License-Identifier: MIT

mod attr;
mod get;
mod handle;

pub(crate) use attr::parse_coalesce_nlas;

pub use attr::EthtoolCoalesceAttr;
pub use get::EthtoolCoalesceGetRequest;
pub use handle::EthtoolCoalesceHandle;
