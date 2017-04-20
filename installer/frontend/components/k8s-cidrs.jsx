import React from 'react';

import { CIDR } from './cidr';
import {
  POD_CIDR,
  SERVICE_CIDR,
} from '../cluster-config';

export const KubernetesCIDRs = ({validator}) => <div className="row form-group">
  <div className="col-xs-12">
    <h4>Kubernetes</h4>
    <CIDR name="Pod Range" field={POD_CIDR} placeholder="10.2.0.0/16" validator={validator} />
    <CIDR name="Service Range" field={SERVICE_CIDR} placeholder="10.3.0.0/24" validator={validator} />
  </div>
</div>;
