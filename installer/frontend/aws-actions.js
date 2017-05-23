import _ from 'lodash';

import * as awsApis from './aws-api';
import { awsActionTypes, configActions } from './actions';
import {
  AWS_ACCESS_KEY_ID,
  AWS_CONTROLLER_SUBNET_IDS,
  AWS_CONTROLLER_SUBNETS,
  AWS_REGION,
  AWS_SECRET_ACCESS_KEY,
  AWS_SESSION_TOKEN,
  AWS_WORKER_SUBNET_IDS,
  AWS_WORKER_SUBNETS,
  PLATFORM_TYPE,
  STS_ENABLED,
} from './cluster-config';
import { BARE_METAL_TF } from './platforms';
import { TectonicGA } from './tectonic-ga';

const { batchSetIn } = configActions;

const createAction = (name, fn, shouldReject=false) => (body, creds, isNow) => (dispatch, getState) => {
  const { clusterConfig } = getState();

  creds = Object.assign({
    AccessKeyID: clusterConfig[AWS_ACCESS_KEY_ID],
    SecretAccessKey: clusterConfig[AWS_SECRET_ACCESS_KEY],
    SessionToken: clusterConfig[STS_ENABLED] && clusterConfig[AWS_SESSION_TOKEN] || '',
    Region: clusterConfig[AWS_REGION],
  }, creds);

  const obj = {
    inFly: true,
  };

  // Don't unset values and errors if we can track causality
  if (!isNow) {
    obj.value = [];
    obj.error = null;
  }

  dispatch({
    type: awsActionTypes.SET,
    payload: {[name]: obj},
  });

  let platform;
  if (clusterConfig[PLATFORM_TYPE] === BARE_METAL_TF) {
    platform = 'metal';
  }
  return fn(body, creds, platform)
    .then(value => {
      if (isNow && !isNow()) {
        return;
      }
      dispatch({
        type: awsActionTypes.SET,
        payload: {
          [name]: {
            inFly: false,
            value: value,
            error: null,
          },
        },
      });
      return value;
    })
    .catch(error => {
      if (isNow && !isNow()) {
        console.log("not now. returning", error);
        return;
      }

      console.error(error.stack || error);

      dispatch({
        type: awsActionTypes.SET,
        payload: {
          [name]: {
            inFly: false,
            value: [],
            error: error,
          },
        },
      });
      if (!shouldReject) {
        return error;
      }

      let message = error.message || error;
      if (message === 'Failed to fetch') {
        message += ` ${name}`;
      }
      TectonicGA.sendEvent('Async Error', 'user input', message);
      return Promise.reject(message);
    });
};

export const getVpcs = createAction('availableVpcs', awsApis.getVpcs);
export const getVpcSubnets = createAction('availableVpcSubnets', awsApis.getVpcSubnets);
export const getSsh = createAction('availableSsh', awsApis.getSsh, true);
export const getRegions = createAction('availableRegions', awsApis.getRegions, true);
export const getZones = createAction('availableR53Zones', awsApis.getZones, true);
export const getDomainInfo = createAction('domainInfo', awsApis.getDomainInfo);
export const validateSubnets = createAction('validateSubnets', awsApis.validateSubnets);
export const TFDestroy = createAction('destroy', awsApis.TFDestroy, true);

const getDefaultSubnets_ = createAction('subnets', awsApis.getDefaultSubnets);

export const getDefaultSubnets = (body, creds, isNow) => (dispatch, getState) =>
  getDefaultSubnets_({vpcCIDR: "10.0.0.0/16"}, creds)(dispatch, getState)
  .then(subnets => {
    if (isNow && !isNow()) {
      return;
    }
    const batches = [];
    _.each(subnets.public, ({availabilityZone, instanceCIDR, id}) => {
      batches.push([`${AWS_CONTROLLER_SUBNETS}.${availabilityZone}`, instanceCIDR]);
      // TODO: (ggreer) stop resetting this? (ditto for worker subnet ids)
      batches.push([`${AWS_CONTROLLER_SUBNET_IDS}.${availabilityZone}`, id]);
    });
    _.each(subnets.private, ({availabilityZone, instanceCIDR, id}) => {
      batches.push([`${AWS_WORKER_SUBNETS}.${availabilityZone}`, instanceCIDR]);
      batches.push([`${AWS_WORKER_SUBNET_IDS}.${availabilityZone}`, id]);
    });
    batchSetIn(dispatch, batches);
  });
