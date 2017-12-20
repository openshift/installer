import _ from 'lodash';
import pluralize from 'pluralize';
import React from 'react';
import { connect } from 'react-redux';

import { Field, Form } from '../form';
import { AWS_TF, BARE_METAL_TF } from '../platforms';
import { validate } from '../validate';
import { Alert } from './alert';
import { CIDRRow, cidrSize } from './cidr';
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

const DockerBridgeWarning = connect(
  ({clusterConfig}) => ({podCidr: clusterConfig[POD_CIDR]})
)(({podCidr}) => {
  if (!validate.CIDR(podCidr)) {
    const ip = podCidr.split('/')[0];
    const octets = ip.split('.').map(octet => parseInt(octet, 10));
    if ((ip.startsWith('172.') && 17 <= octets[1] && octets[1] <= 31) ||
        (ip.startsWith('192.168.') && octets[2] <= 15)) {
      return <Alert>
        <b>Pod range may conflict with Docker Bridge subnet</b><br />
        Docker Bridge may use any IP range from 172.17-31.x.x or 192.168.0-15.x.
      </Alert>;
    }
  }
  return null;
});

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
      <b>Pod range too small</b><br />
      {maxNodes === 0 ? 'No nodes' : `Only ${maxNodes} of your ${nodes} ${pluralize('node', nodes)}`} can fit within the pod range, since each node requires a minimum of 256 IP addresses.
    </Alert>;
  }
  return <Alert>
    <b>Pod range mostly assigned</b><br />
    Only {maxNodes} {pluralize('node', maxNodes)} can fit within the pod range, since each node requires a minimum of 256 IP addresses. You have selected {nodes} {pluralize('node', nodes)}.
  </Alert>;
});

const DEFAULT_POD_CIDR = '10.2.0.0/16';
const DEFAULT_SERVICE_CIDR = '10.3.0.0/16';

const validator = (s, cc) => cc[PLATFORM_TYPE] === AWS_TF ? validate.AWSsubnetCIDR(s) : validate.CIDR(s);

const form = new Form('K8S_CIDRS_FORM', [
  new Field(POD_CIDR, {default: DEFAULT_POD_CIDR, validator}),
  new Field(SERVICE_CIDR, {default: DEFAULT_SERVICE_CIDR, validator}),
]);

export const KubernetesCIDRs = ({autoFocus = true}) => <div id="k8sCIDRs" className="row form-group">
  <div className="col-xs-12">
    <h4>Kubernetes</h4>
    <CIDRRow autoFocus={autoFocus} name="Pod Range" field={POD_CIDR} placeholder={DEFAULT_POD_CIDR} />
    <CIDRRow name="Service Range" field={SERVICE_CIDR} placeholder={DEFAULT_SERVICE_CIDR} />
    <DockerBridgeWarning />
    <PodRangeWarning />
  </div>
</div>;

KubernetesCIDRs.canNavigateForward = form.canNavigateForward;
