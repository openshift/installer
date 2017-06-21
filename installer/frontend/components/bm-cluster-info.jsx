import React from 'react';

import { Input, Connect } from './ui';
import { TectonicLicense, licenseForm } from './tectonic-license';
import { ExperimentalFeatures } from './experimental-features';
import { CLUSTER_NAME } from '../cluster-config';
import { Form } from '../form';
import fields from '../fields';

const clusterInfoForm = new Form("BM_ClusterInfo", [
  licenseForm,
  fields[CLUSTER_NAME],
]);

export const BM_ClusterInfo = () => {
  return (
    <div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={CLUSTER_NAME}>Cluster Name</label>
        </div>
        <div className="col-xs-9">
          <Connect field={CLUSTER_NAME}>
            <Input placeholder="production" autoFocus="true" />
          </Connect>
          <p className="text-muted">Give this cluster a name that will help you identify it.</p>
        </div>
      </div>
      <ExperimentalFeatures />
      <TectonicLicense />
    </div>
  );
};

BM_ClusterInfo.canNavigateForward = clusterInfoForm.canNavigateForward;
