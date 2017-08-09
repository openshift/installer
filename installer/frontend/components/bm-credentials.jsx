import React from 'react';
import { validate } from '../validate';
import {
  BM_MATCHBOX_CA,
  BM_MATCHBOX_CLIENT_CERT,
  BM_MATCHBOX_CLIENT_KEY,
} from '../cluster-config';

import { CertArea, PrivateKeyArea, WithClusterConfig } from './ui';

export const BM_Credentials = () => {
  return (
    <div>
      <div className="form-group">
        Credentials were generated during the matchbox installation. Provide the certificates and keys here:
      </div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={BM_MATCHBOX_CA}>CA Certificate</label></div>
        <div className="col-xs-9">
          <WithClusterConfig field={BM_MATCHBOX_CA}>
            <CertArea id={BM_MATCHBOX_CA}
              autoFocus="true" />
          </WithClusterConfig>
        </div>
      </div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={BM_MATCHBOX_CLIENT_CERT}>Client Certificate</label></div>
        <div className="col-xs-9">
          <WithClusterConfig field={BM_MATCHBOX_CLIENT_CERT}>
            <CertArea id={BM_MATCHBOX_CLIENT_CERT} />
          </WithClusterConfig>
        </div>
      </div>
      <div className="row form-group">
        <div className="col-xs-3">
          <label htmlFor={BM_MATCHBOX_CLIENT_KEY}>Client Key</label></div>
        <div className="col-xs-9">
          <WithClusterConfig field={BM_MATCHBOX_CLIENT_KEY} >
            <PrivateKeyArea id={BM_MATCHBOX_CLIENT_KEY} />
          </WithClusterConfig>
        </div>
      </div>
    </div>
  );
};
BM_Credentials.canNavigateForward = ({clusterConfig}) => {
  return !validate.certificate(clusterConfig[BM_MATCHBOX_CA]) &&
         !validate.certificate(clusterConfig[BM_MATCHBOX_CLIENT_CERT]) &&
         !validate.privateKey(clusterConfig[BM_MATCHBOX_CLIENT_KEY]);
};
