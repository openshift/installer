use log::warn;

use crate::{
    BaseInterface, BondAdSelect, BondAllPortsActive, BondArpAllTargets,
    BondArpValidate, BondConfig, BondFailOverMac, BondInterface, BondLacpRate,
    BondMode, BondOptions, BondPrimaryReselect, BondXmitHashPolicy,
};

pub(crate) fn np_bond_to_nmstate(
    np_iface: &nispor::Iface,
    base_iface: BaseInterface,
) -> BondInterface {
    let mut bond_iface = BondInterface::new();
    let mut bond_conf = BondConfig::new();

    bond_iface.base = base_iface;
    bond_conf.options = Some(np_bond_options_to_nmstate(np_iface));
    if let Some(np_bond) = &np_iface.bond {
        bond_conf.port = Some(
            np_bond
                .subordinates
                .as_slice()
                .iter()
                .map(|iface_name| iface_name.to_string())
                .collect(),
        );
        bond_conf.mode = match np_bond.mode {
            nispor::BondMode::BalanceRoundRobin => Some(BondMode::RoundRobin),
            nispor::BondMode::ActiveBackup => Some(BondMode::ActiveBackup),
            nispor::BondMode::BalanceXor => Some(BondMode::XOR),
            nispor::BondMode::Broadcast => Some(BondMode::Broadcast),
            nispor::BondMode::Ieee8021AD => Some(BondMode::LACP),
            nispor::BondMode::BalanceTlb => Some(BondMode::TLB),
            nispor::BondMode::BalanceAlb => Some(BondMode::ALB),
            _ => {
                warn!("Unsupported bond mode");
                Some(BondMode::Unknown)
            }
        };
    }
    bond_iface.bond = Some(bond_conf);
    bond_iface
}

fn np_bond_options_to_nmstate(np_iface: &nispor::Iface) -> BondOptions {
    let mut options = BondOptions::default();
    if let Some(ref np_bond) = &np_iface.bond {
        options.ad_actor_sys_prio = np_bond.ad_actor_sys_prio;
        options.ad_actor_system = np_bond.ad_actor_system.clone();
        options.ad_select = np_bond.ad_select.as_ref().and_then(|r| match r {
            nispor::BondAdSelect::Stable => Some(BondAdSelect::Stable),
            nispor::BondAdSelect::Bandwidth => Some(BondAdSelect::Bandwidth),
            nispor::BondAdSelect::Count => Some(BondAdSelect::Count),
            _ => {
                warn!("Unsupported bond ad_select option {:?}", r);
                None
            }
        });
        options.ad_user_port_key = np_bond.ad_user_port_key;
        options.all_slaves_active = np_bond
            .all_subordinates_active
            .as_ref()
            .and_then(|r| match r {
                nispor::BondAllSubordinatesActive::Dropped => {
                    Some(BondAllPortsActive::Dropped)
                }
                nispor::BondAllSubordinatesActive::Delivered => {
                    Some(BondAllPortsActive::Delivered)
                }
                _ => {
                    warn!("Unsupported bond all ports active options {:?}", r);
                    None
                }
            });
        options.arp_all_targets =
            np_bond.arp_all_targets.as_ref().and_then(|r| match r {
                nispor::BondModeArpAllTargets::Any => {
                    Some(BondArpAllTargets::Any)
                }
                nispor::BondModeArpAllTargets::All => {
                    Some(BondArpAllTargets::All)
                }
                _ => {
                    warn!("Unsupported bond arp_all_targets option {:?}", r);
                    None
                }
            });
        options.arp_interval = np_bond.arp_interval;
        options.arp_ip_target = np_bond.arp_ip_target.clone();
        options.arp_validate =
            np_bond.arp_validate.as_ref().and_then(|r| match r {
                nispor::BondArpValidate::None => Some(BondArpValidate::None),
                nispor::BondArpValidate::Active => {
                    Some(BondArpValidate::Active)
                }
                nispor::BondArpValidate::Backup => {
                    Some(BondArpValidate::Backup)
                }
                nispor::BondArpValidate::All => Some(BondArpValidate::All),
                nispor::BondArpValidate::FilterActive => {
                    Some(BondArpValidate::FilterActive)
                }
                nispor::BondArpValidate::FilterBackup => {
                    Some(BondArpValidate::FilterBackup)
                }
                _ => {
                    warn!("Unsupported bond arp_validate options {:?}", r);
                    None
                }
            });
        options.downdelay = np_bond.downdelay;
        options.fail_over_mac =
            np_bond.fail_over_mac.as_ref().and_then(|r| match r {
                nispor::BondFailOverMac::None => Some(BondFailOverMac::None),
                nispor::BondFailOverMac::Active => {
                    Some(BondFailOverMac::Active)
                }
                nispor::BondFailOverMac::Follow => {
                    Some(BondFailOverMac::Follow)
                }
                _ => {
                    warn!("Unsupported bond fail_over_mac options {:?}", r);
                    None
                }
            });
        options.lacp_rate = np_bond.lacp_rate.as_ref().and_then(|r| match r {
            nispor::BondLacpRate::Slow => Some(BondLacpRate::Slow),
            nispor::BondLacpRate::Fast => Some(BondLacpRate::Fast),
            _ => {
                warn!("Unsupported bond lacp_rate options {:?}", r);
                None
            }
        });
        options.lp_interval = np_bond.lp_interval;
        options.miimon = np_bond.miimon;
        options.min_links = np_bond.min_links;
        options.num_grat_arp = np_bond.num_grat_arp;
        options.num_unsol_na = np_bond.num_unsol_na;
        options.packets_per_slave = np_bond.packets_per_subordinate;
        options.primary = np_bond.primary.clone();
        options.primary_reselect =
            np_bond.primary_reselect.as_ref().and_then(|r| match r {
                nispor::BondPrimaryReselect::Always => {
                    Some(BondPrimaryReselect::Always)
                }
                nispor::BondPrimaryReselect::Better => {
                    Some(BondPrimaryReselect::Better)
                }
                nispor::BondPrimaryReselect::Failure => {
                    Some(BondPrimaryReselect::Failure)
                }
                _ => {
                    warn!("Unsupported bond primary_reselect options {:?}", r);
                    None
                }
            });
        options.resend_igmp = np_bond.resend_igmp;
        options.tlb_dynamic_lb = np_bond.tlb_dynamic_lb;
        options.updelay = np_bond.updelay;
        options.use_carrier = np_bond.use_carrier;
        options.xmit_hash_policy =
            np_bond.xmit_hash_policy.as_ref().and_then(|r| match r {
                nispor::BondXmitHashPolicy::Layer2 => {
                    Some(BondXmitHashPolicy::Layer2)
                }
                nispor::BondXmitHashPolicy::Layer34 => {
                    Some(BondXmitHashPolicy::Layer34)
                }
                nispor::BondXmitHashPolicy::Layer23 => {
                    Some(BondXmitHashPolicy::Layer23)
                }
                nispor::BondXmitHashPolicy::Encap23 => {
                    Some(BondXmitHashPolicy::Encap23)
                }
                nispor::BondXmitHashPolicy::Encap34 => {
                    Some(BondXmitHashPolicy::Encap34)
                }
                nispor::BondXmitHashPolicy::VlanSrcMac => {
                    Some(BondXmitHashPolicy::VlanSrcMac)
                }
                _ => {
                    warn!("Unsupported bond xmit_hash_policy options {:?}", r);
                    None
                }
            });
    }
    options
}
