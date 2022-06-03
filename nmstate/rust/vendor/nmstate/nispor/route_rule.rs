use crate::{RouteRuleEntry, RouteRules};

pub(crate) fn get_route_rules(np_rules: &[nispor::RouteRule]) -> RouteRules {
    let mut ret = RouteRules::new();

    let mut rules = Vec::new();
    for np_rule in np_rules {
        let mut rule = RouteRuleEntry::new();
        // We only support route rules with 'table' action
        if np_rule.action != nispor::RuleAction::Table {
            continue;
        }
        // Neither ip_from or ip_to should be defeind
        if np_rule.dst.is_none() && np_rule.src.is_none() {
            continue;
        }
        if np_rule.dst.as_deref() == Some("")
            && np_rule.src.as_deref() == Some("")
        {
            continue;
        }
        rule.ip_to = np_rule.dst.clone();
        rule.ip_from = np_rule.src.clone();
        rule.table_id = np_rule.table;
        rule.priority = np_rule.priority.map(i64::from);
        rules.push(rule);
    }
    ret.config = Some(rules);

    ret
}
