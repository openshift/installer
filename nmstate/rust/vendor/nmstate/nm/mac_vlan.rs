use crate::nm::nm_dbus::NmSettingMacVlan;

use crate::{MacVlanConfig, MacVtapConfig};

impl From<&MacVlanConfig> for NmSettingMacVlan {
    fn from(config: &MacVlanConfig) -> Self {
        let mut settings = NmSettingMacVlan::default();
        settings.mode = Some(config.mode.into());
        settings.parent = Some(config.base_iface.clone());
        settings.tap = Some(false);
        if let Some(v) = config.accept_all_mac {
            settings.accept_all_mac = Some(v);
        }
        settings
    }
}

impl From<&MacVtapConfig> for NmSettingMacVlan {
    fn from(config: &MacVtapConfig) -> Self {
        let mut settings = NmSettingMacVlan::default();
        settings.mode = Some(config.mode.into());
        settings.parent = Some(config.base_iface.clone());
        settings.tap = Some(true);
        if let Some(v) = config.accept_all_mac {
            settings.accept_all_mac = Some(v);
        }
        settings
    }
}
