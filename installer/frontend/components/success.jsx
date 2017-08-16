import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { AWS_HOSTED_ZONE_ID, CLUSTER_SUBDOMAIN, PLATFORM_TYPE} from '../cluster-config';
import { TectonicGA } from '../tectonic-ga';

const handleAllDone = (platformType) => TectonicGA.sendEvent('Installer Button', 'click', 'User clicks over to the console', platformType);

const stateToProps = ({cluster, clusterConfig}) => {
  let tectonicConsole = _.get(cluster, ['status', 'tectonicConsole', 'instance']);
  if (!tectonicConsole) {
    tectonicConsole = clusterConfig.tectonicDomain;
  }
  if (!tectonicConsole) {
    // TODO: (kans) add this to the terraform status response
    const hostedZoneID = clusterConfig[AWS_HOSTED_ZONE_ID];
    const domain = _.get(clusterConfig, ['extra', AWS_HOSTED_ZONE_ID, 'zoneToName', hostedZoneID]);
    tectonicConsole = clusterConfig[CLUSTER_SUBDOMAIN] + (clusterConfig[CLUSTER_SUBDOMAIN].endsWith('.') ? '' : '.') + domain;
  }
  return {tectonicConsole, platformType: clusterConfig[PLATFORM_TYPE]};
};

export const Success = connect(stateToProps)(
  ({tectonicConsole, platformType}) => <div>
    <div className="row">
      <div className="col-xs-12">
        <h4>1. Cluster assets are important, save them now!</h4>
        <p>Download and keep your cluster assets in a safe place. These are needed to destroy, replicate or quickly reinstall.</p>
        <a href="/terraform/assets" download>
          <button className="btn btn-default" style={{marginTop: 10}}>
            <i className="fa fa-download"></i>&nbsp;&nbsp;Download Assets
          </button>
        </a>
      </div>
    </div>

    <hr className="spacer" />

    <div className="row">
      <div className="col-xs-12">
        <h4>2. Check out Tectonic Console</h4>
        <p>The Tectonic Console gives you an easy-to-navigate view of your cluster.</p>
        <a href={`https://${tectonicConsole}`} target="_blank">
          <button className="btn btn-primary wiz-giant-button"
            style={{marginTop: 20, marginBottom: 80}}
            onClick={() => handleAllDone(platformType)}>Go to my Tectonic Console&nbsp;&nbsp;<i className="fa fa-external-link"></i></button>
        </a>
      </div>
    </div>
  </div>
);

Success.canNavigateForward = () => false;
