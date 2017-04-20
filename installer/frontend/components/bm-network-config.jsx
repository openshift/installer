import React from 'react';

import { validate } from '../validate';
import { KubernetesCIDRs } from './k8s-cidrs';
import {
  POD_CIDR,
  SERVICE_CIDR,
} from '../cluster-config';

export const BM_NetworkConfig = () => <KubernetesCIDRs validator={validate.CIDR} />;

BM_NetworkConfig.canNavigateForward = ({clusterConfig}) => {
  return !validate.CIDR(clusterConfig[POD_CIDR]) &&
         !validate.CIDR(clusterConfig[SERVICE_CIDR]);
};
