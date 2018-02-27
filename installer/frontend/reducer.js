import _ from 'lodash';
import { combineReducers } from 'redux';
import { fromJS } from 'immutable';

import { DEFAULT_CLUSTER_CONFIG } from './cluster-config';

import {
  awsActionTypes,
  clusterReadyActionTypes,
  configActionTypes,
  dirtyActionTypes,
  eventErrorsActionTypes,
  loadFactsActionTypes,
  restoreActionTypes,
  serverActionTypes,
  commitPhases,
} from './actions';

const UNLOADED_AWS_VALUE = {
  inFly: false,
  value: [],
  error: null,
};

const DEFAULT_AWS = {
  availableVpcs: UNLOADED_AWS_VALUE,
  availableVpcSubnets: UNLOADED_AWS_VALUE,
  availableSsh: UNLOADED_AWS_VALUE,
  availableIamRoles: UNLOADED_AWS_VALUE,
  availableRegions: UNLOADED_AWS_VALUE,
  availableKms: UNLOADED_AWS_VALUE,
  availableR53Zones: UNLOADED_AWS_VALUE,
  createdKms: UNLOADED_AWS_VALUE,
  domainInfo: UNLOADED_AWS_VALUE,
  subnets: UNLOADED_AWS_VALUE,
  validateSubnets: UNLOADED_AWS_VALUE,
  destroy: UNLOADED_AWS_VALUE,
};

const reducersTogether = combineReducers({

  // State machine associated with server submissions
  commitState: (state, action) => {
    if (!state) {
      return {
        phase: commitPhases.IDLE,
        response: null,
      };
    }

    switch (action.type) {
    case serverActionTypes.COMMIT_REQUESTED:
      return {
        phase: commitPhases.REQUESTED,
      };
    case serverActionTypes.COMMIT_SENT:
      return {
        phase: commitPhases.WAITING,
      };
    case serverActionTypes.COMMIT_SUCCESSFUL:
      console.log('COMMIT SUCCESSFUL: provisioner has been configured');
      console.dir(action.payload);
      return {
        phase: commitPhases.SUCCEEDED,
        response: action.payload,
      };
    case serverActionTypes.COMMIT_FAILED:
      console.log('COMMIT FAILED');
      console.dir(action.payload);
      return {
        phase: commitPhases.FAILED,
        response: action.payload,
      };
    default:
      return state;
    }
  },

  aws: (state, action) => {
    if (!state) {
      return DEFAULT_AWS;
    }

    switch (action.type) {
    case awsActionTypes.SET:
      Object.keys(action.payload).forEach(k => {
        if (!DEFAULT_AWS.hasOwnProperty(k)) {
          throw Error(`attempt to set aws property ${k} missing from defaults`);
        }
      });
      return Object.assign({}, state, action.payload);
    default:
      return state;
    }
  },

  // The user's intentions for their cluster
  clusterConfig: (state, action) => {
    if (!state) {
      return DEFAULT_CLUSTER_CONFIG;
    }

    switch (action.type) {
    case configActionTypes.RESET:
      return {};
    case configActionTypes.SET:
      Object.keys(action.payload).forEach(k => {
        if (!DEFAULT_CLUSTER_CONFIG.hasOwnProperty(k)) {
          throw Error(`attempt to set cluster property ${k} missing from defaults`);
        }
      });
      return Object.assign({}, state, action.payload);

    case configActionTypes.BATCH_SET_IN: {
      let object = fromJS(state);
      action.payload.forEach(batch => {
        const [path, value] = batch;
        const array = _.isString(path) ? path.split('.') : path;
        object = object.setIn(array, value);
      });
      return object.toJS();
    }
    case configActionTypes.SET_IN: {
      const {path, value} = action.payload;
      const pathArray = _.isString(path) ? path.split('.') : path;
      return fromJS(state).setIn(pathArray, value).toJS();
    }
    case configActionTypes.DELETE_IN:
      return fromJS(state).deleteIn(action.payload).toJS();
    default:
      return state;
    }
  },

  // Errors resulting from server states, user uploads, network issues, or other transient phenomena
  eventErrors: (state, action) => {
    if (!state) {
      return {};
    }

    switch (action.type) {
    case eventErrorsActionTypes.ERROR:
      return Object.assign({}, state, {
        [action.payload.name]: action.payload.error,
      });
    default:
      return state;
    }
  },

  // Facts the server knows at load time, that we have to get asynchronously
  serverFacts: (state, action) => {
    if (state === undefined) {
      return {
        loaded: false,
        error: null,
      };
    }

    switch (action.type) {
    case loadFactsActionTypes.LOADED:
      return Object.assign({}, action.payload, {
        loaded: true,
        error: null,
      });
    case loadFactsActionTypes.ERROR:
      return {
        loaded: true,
        error: action.payload,
      };
    default:
      return state;
    }
  },

  // The status of the cluster
  cluster: (state, action) => {
    if (!state) {
      return {
        loaded: false,
        ready: false,
        error: null,
        status: null,
      };
    }

    switch (action.type) {
    case clusterReadyActionTypes.NOT_READY:
      return Object.assign({}, state, {
        loaded: true,
        ready: false,
        status: null,
      });
    case clusterReadyActionTypes.STATUS:
      return Object.assign({}, state, {
        loaded: true,
        ready: true,
        status: action.payload,
        error: null,
      });
    case clusterReadyActionTypes.ERROR:
      return Object.assign({}, state, {
        loaded: true,
        error: action.payload,
      });
    default:
      return state;
    }
  },

  // Stores the "dirtiness" of UI fields
  dirty: (state, action) => {
    if (!state) {
      return {};
    }

    switch (action.type) {
    case dirtyActionTypes.ADD: {
      const pathArray = _.isString(action.payload) ? action.payload.split('.') : action.payload;
      return fromJS(state).setIn(pathArray, true).toJS();
    }
    case dirtyActionTypes.DELETE_IN:
      return fromJS(state).deleteIn(action.payload).toJS();
    default:
      return state;
    }
  },
});

// Shim for tests
export const reducer = (state, action) => {
  let restored;
  switch (action.type) {
  case restoreActionTypes.RESTORE_STATE:
    restored = Object.assign({}, action.payload, {
      serverFacts: state.serverFacts,
      cluster: state.cluster,
    });
    restored.aws = _.defaults(restored.aws, DEFAULT_AWS);
    if (restored.clusterConfig) {
      restored.clusterConfig = _.defaults(restored.clusterConfig, DEFAULT_CLUSTER_CONFIG);
    }
    break;
  default:
    restored = state;
  }

  return reducersTogether(restored, action);
};

// Returns only the bits of state that should be saved and restored
export const savable = (state) => {
  const {dirty, clusterConfig} = state;
  return {dirty, clusterConfig: _.omit(clusterConfig, ['error', 'extra', 'extraError', 'extraInFly', 'inFly'])};
};
