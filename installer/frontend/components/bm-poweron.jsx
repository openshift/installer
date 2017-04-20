import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { SaveAssets } from './save-assets';
import { WaitingLi } from './ui';

export const BM_PowerOn = connect(
  ({cluster, commitState}) => {
    let controllers = [];
    let etcdByIP = {};
    let workers = [];

    if (cluster.status) {
      if (cluster.status.kubelet) {
        controllers = cluster.status.kubelet.controllers;
        workers = cluster.status.kubelet.workers;
      }
      if (cluster.status.etcd) {
        cluster.status.etcd.forEach(node => {
          etcdByIP[node.instance] = node;
        });
      }
    }

    return {
      blob: commitState.response,
      error: cluster.error,
      controllers: _.map(controllers, node => {
        return {
          instance: node.instance,
          ready: node.ready,
          etcd: etcdByIP[node.instance],
        };
      }),
      workers: workers,
    };
  }
)(({blob, error, controllers, workers}) => {
  const controllerLis = _.map(controllers, (node, ix) => {
    return (
      <WaitingLi done={node.ready} key={`controller-${ix}`}>
        <span title={node.instance}>Controller #{ix + 1}</span>
        <ul>
          {
            node.etcd && <WaitingLi substep={true} done={node.etcd.ready}>
              {node.etcd.ready ? 'etcd is running' : 'Waiting for etcd'}
            </WaitingLi>
          }
          <WaitingLi substep={true} done={node.ready}>
            {node.ready ? 'Kubernetes is running' : 'Waiting for kubelet to be ready'}
          </WaitingLi>
        </ul>
      </WaitingLi>
    );
  });

  const workerLis = _.map(workers, (node, ix) => {
    return (
      <WaitingLi done={node.ready} key={`worker-${ix}`}>
        <span title={node.instance}>Worker #{ix + 1}</span>
        <ul>
          <WaitingLi substep={true} done={node.ready}>
            {node.ready ? 'Kubernetes is running' : 'Waiting for kubelet to be ready'}
          </WaitingLi>
        </ul>
      </WaitingLi>
    );
  });

  return (
    <div>
      <div className="row">
        <div className="col-sm-12">
          <div>We're ready to boot the cluster.</div>
          <div className="wiz-herotext">
            <span className="fa fa-power-off wiz-herotext-icon"></span> Power on the nodes
          </div>
          <div className="form-group">After powering up, your nodes will provision themselves automatically.
            This process can take up to 30 minutes, while the following happens.</div>
          <div className="form-group">
            <ul>
              <li>Container Linux is downloaded and installed to disk (about 200 MB)</li>
              <li>Cluster software is downloaded (about 500 MB)</li>
              <li>One or two reboots may occur</li>
            </ul>
          </div>
          <div className="form-group">
            <div className="wiz-launch__progress">
              <ul className="wiz-launch-progress">
                {controllerLis}
                {workerLis}
              </ul>
            </div>
          </div>
          { error &&
            <div className="form-group">
              <div className='wiz-error-message'>{error.toString()}</div>
            </div>
          }
          <div className="form-group">
            <SaveAssets blob={blob} />
          </div>
        </div>
      </div>
    </div>
  );
});

// horrible hack to prevent a flapping cluster from redirecting user from connect page back to poweron page
let ready = false;

BM_PowerOn.canNavigateForward = ({cluster}) => {
  if (!cluster.ready) {
    return false;
  }

  const {controllers, workers} = cluster.status.kubelet;
  ready = ready || (_.every(controllers, n => n.ready) && _.some(workers, n => n.ready));
  return ready;
};
