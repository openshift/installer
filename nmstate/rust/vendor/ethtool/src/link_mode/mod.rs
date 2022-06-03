// SPDX-License-Identifier: MIT

mod attr;
mod get;
mod handle;

pub(crate) use attr::parse_link_mode_nlas;
pub use attr::{EthtoolLinkModeAttr, EthtoolLinkModeDuplex};
pub use get::EthtoolLinkModeGetRequest;
pub use handle::EthtoolLinkModeHandle;
