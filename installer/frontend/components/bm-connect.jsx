import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { SSHInstructions } from './ssh-instructions';
import { SaveAssets } from './save-assets';
import { WaitingLi } from './ui';

export const BM_Connect = connect(
  ({cluster, commitState}) => ({
    blob: commitState.response,
    controllerIP: _.get(cluster, 'status.kubelet.controllers[0].instance'),
    tectonicConsole: cluster.status.tectonicConsole,
  })
)(({tectonicConsole, controllerIP, blob}) => <div className="row">
  <div className="col-sm-12">
    <div>We're ready to connect your nodes! First, download your cluster assets.</div>
    <SaveAssets blob={blob} />
    <SSHInstructions controllerIP={controllerIP} />
    <div className="wiz-launch__progress">
      <ul className="wiz-launch-progress">
        <WaitingLi done={tectonicConsole.ready}>
          <span title={tectonicConsole.instance}>Tectonic Console</span>
        </WaitingLi>
      </ul>
    </div>
  </div>
</div>);

BM_Connect.canNavigateForward = ({cluster}) => {
  return cluster.ready && cluster.status.tectonicConsole.ready;
};
