import React from 'react';

import { BM_MASTERS_COUNT, BM_WORKERS_COUNT } from '../cluster-config';
import { validate } from '../validate';
import { WithClusterConfig, NumberInput } from './ui';

const TOO_MANY_MASTERS = 9;
const TOO_MANY_WORKERS = 1000;

const validateMasters = v => validate.int({min: 1, max: TOO_MANY_MASTERS})(v) || validate.isOdd(v);
const validateWorkers = validate.int({min: 1, max: TOO_MANY_WORKERS});


export const BM_Size = () => <div>
  <div className="form-group">
    Every node in your cluster is either a master or a worker.
  </div>
  <div className="form-group">
    Master nodes run essential cluster services and don't run end-user apps.
  </div>
  <div className="form-group">
    <div className="row">
      <div className="col-xs-3">
        <label htmlFor={BM_MASTERS_COUNT}>Masters</label>
      </div>
      <div className="col-xs-9">
        <WithClusterConfig field={BM_MASTERS_COUNT} validator={validateMasters}>
          <NumberInput
            id={BM_MASTERS_COUNT}
            className="wiz-super-short-input"
            min="1"
            max={TOO_MANY_MASTERS} />
        </WithClusterConfig>
        <p className="text-muted">An odd number of masters is required.</p>
      </div>
    </div>
  </div>
  <div className="form-group">
    Worker nodes run your end-user's apps. The cluster software automatically shares load
    between these nodes.
  </div>
  <div className="form-group">
    <div className="row">
      <div className="col-xs-3">
        <label htmlFor={BM_WORKERS_COUNT}>Workers</label>
      </div>
      <div className="col-xs-9">
        <WithClusterConfig field={BM_WORKERS_COUNT} validator={validateWorkers}>
          <NumberInput
            id={BM_WORKERS_COUNT}
            className="wiz-super-short-input"
            min="1"
            max={TOO_MANY_WORKERS} />
        </WithClusterConfig>
        <p className="text-muted">Workers can be added and removed at any time.</p>
      </div>
    </div>
  </div>
</div>;

BM_Size.canNavigateForward = ({clusterConfig}) => {
  return !validateMasters(clusterConfig[BM_MASTERS_COUNT]) && !validateWorkers(clusterConfig[BM_WORKERS_COUNT]);
};
