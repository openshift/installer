import React from 'react';
import { CheckBox, WithClusterConfig } from './ui';

import { UPDATER_ENABLED } from '../cluster-config';

export const ExperimentalFeatures = () => <div className="row form-group">
  <div className="col-xs-3">
    <label htmlFor="tectonicOperator">Automated Updates</label>
  </div>
  <div className="col-xs-9">
    <WithClusterConfig field={UPDATER_ENABLED}>
      <CheckBox id="tectonicOperator" suffix={
        <label htmlFor="tectonicOperator">
          Enable one-click updates for Tectonic, etcd, Prometheus.
        </label>
      } />
    </WithClusterConfig>
    <p className="text-muted checkbox-helper-text">Only use with non-production clusters.</p>
  </div>
</div>;
