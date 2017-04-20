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
  STS_ENABLED,
} from './cluster-config';
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

  dispatch({
    type: awsActionTypes.SET,
    payload: {
      [name]: {
        inFly: true,
        value: [],
        error: null,
      },
    },
  });

  // TODO: pass in platform
  return fn(body, creds)
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
export const getKms = createAction('availableKms', awsApis.getKms);
export const createKms = createAction('createdKms', awsApis.createKms);
export const getSsh = createAction('availableSsh', awsApis.getSsh);
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
