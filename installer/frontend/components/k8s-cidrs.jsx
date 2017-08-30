import _ from 'lodash';
import pluralize from 'pluralize';
import React from 'react';
import { connect } from 'react-redux';

import { AWS_TF, BARE_METAL_TF } from '../platforms';
import { Alert } from './alert';
import { CIDR, cidrSize } from './cidr';
import {
  AWS_CONTROLLERS,
  AWS_WORKERS,
  BM_MASTERS,
  BM_WORKERS,
  NUMBER_OF_INSTANCES,
  PLATFORM_TYPE,
  POD_CIDR,
  SERVICE_CIDR,
} from '../cluster-config';

const PodRangeWarning = connect(
  ({clusterConfig: cc}) => {
    let nodes;
    if (cc[PLATFORM_TYPE] === AWS_TF) {
      nodes = _.get(cc, `${AWS_CONTROLLERS}-${NUMBER_OF_INSTANCES}`, 0) + _.get(cc, `${AWS_WORKERS}-${NUMBER_OF_INSTANCES}`, 0);
    }
    if (cc[PLATFORM_TYPE] === BARE_METAL_TF) {
      nodes = _.get(cc, `${BM_MASTERS}.length`, 0) + _.get(cc, `${BM_WORKERS}.length`, 0);
    }
    return {nodes, size: cidrSize(_.get(cc, POD_CIDR))};
  }
)(({nodes, size}) => {
  if (nodes === undefined || size === undefined) {
    return null;
  }

  // Flannel assigns a minimum network size of /24 (256 IP addresses)
  const maxNodes = Math.floor(size / 256);

  const utilization = nodes / maxNodes;

  if (utilization < 0.75) {
    return null;
  }

  if (utilization > 1) {
    return <Alert severity="error">
      <b>Pod Range Too Small</b><br />
      {maxNodes === 0 ? 'No nodes' : `Only ${maxNodes} of your ${nodes} ${pluralize('node', nodes)}`} can fit within the pod range, since each node requires a minimum of 256 IP addresses.
    </Alert>;
  }
  return <Alert>
    <b>Pod Range Mostly Assigned</b><br />
    Only {maxNodes} {pluralize('node', maxNodes)} can fit within the pod range, since each node requires a minimum of 256 IP addresses. You have selected {nodes} {pluralize('node', nodes)}.
  </Alert>;
});

export const KubernetesCIDRs = ({validator}) => <div className="row form-group">
  <div className="col-xs-12">
    <h4>Kubernetes</h4>
    <CIDR name="Pod Range" field={POD_CIDR} placeholder="10.2.0.0/16" validator={validator} />
    <CIDR name="Service Range" field={SERVICE_CIDR} placeholder="10.3.0.0/16" validator={validator} />
    <PodRangeWarning />
  </div>
</div>;
