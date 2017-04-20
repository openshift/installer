import React from 'react';

import { validate } from '../validate';
import { Input, WithClusterConfig } from './ui';
import { TectonicLicense, licenseForm } from './tectonic-license';
import { ExperimentalFeatures } from './experimental-features';
import { CLUSTER_NAME } from '../cluster-config';

export const BM_ClusterInfo = () => {
  return (
    <div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={CLUSTER_NAME}>Cluster Name</label>
        </div>
        <div className="col-xs-9">
          <WithClusterConfig field={CLUSTER_NAME} validator={validate.k8sName}>
            <Input
                id={CLUSTER_NAME}
                placeholder="production"
                autoFocus="true" />
          </WithClusterConfig>
          <p className="text-muted">Give this cluster a name that will help you identify it.</p>
        </div>
      </div>
      <TectonicLicense />
      <br />
      <ExperimentalFeatures />
    </div>
  );
};
BM_ClusterInfo.canNavigateForward = ({clusterConfig}) => {
  return !validate.k8sName(clusterConfig[CLUSTER_NAME]) &&
    licenseForm.isValid(clusterConfig);
};
