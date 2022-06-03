use serde::{Deserialize, Serialize};

use crate::{BaseInterface, InterfaceType};

#[derive(Debug, Clone, PartialEq, Eq, Serialize, Deserialize)]
#[non_exhaustive]
pub struct DummyInterface {
    #[serde(flatten)]
    pub base: BaseInterface,
}

impl Default for DummyInterface {
    fn default() -> Self {
        let mut base = BaseInterface::new();
        base.iface_type = InterfaceType::Dummy;
        Self { base }
    }
}

impl DummyInterface {
    pub fn new() -> Self {
        Self::default()
    }
}
