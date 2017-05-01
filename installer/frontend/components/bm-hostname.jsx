import React from 'react';

import { validate } from '../validate';
import { CONTROLLER_DOMAIN, BM_TECTONIC_DOMAIN } from '../cluster-config';
import { WithClusterConfig, Input } from './ui';

export const BM_Hostname = () => {
  return (
    <div>
      <div className="form-group">Enter a DNS name which resolves to any master node.
        This is the name we'll use when configuring components.</div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor="rootDomain">Master DNS</label>
        </div>
        <div className="col-xs-9">
          <WithClusterConfig field={CONTROLLER_DOMAIN} validator={validate.domainName}>
            <Input
                id="rootDomain"
                placeholder="masters.example.com"
                autoFocus="true" />
          </WithClusterConfig>
        </div>
      </div>
      <div className="form-group">Enter a DNS name which resolves to any worker node(s) for Ingress.</div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor="rootDomain">Tectonic DNS</label>
        </div>
        <div className="col-xs-9">
          <WithClusterConfig field={BM_TECTONIC_DOMAIN} validator={validate.domainName}>
            <Input
                id="rootDomain"
                placeholder="tectonic.example.com"
                autoFocus="true" />
          </WithClusterConfig>
        </div>
      </div>
    </div>
  );
};
BM_Hostname.canNavigateForward = ({clusterConfig}) => {
  return (!validate.domainName(clusterConfig[CONTROLLER_DOMAIN]) &&
          !validate.domainName(clusterConfig[BM_TECTONIC_DOMAIN]));
};
