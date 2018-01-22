import _ from 'lodash';

import { DEFAULT_CLUSTER_CONFIG } from './cluster-config';

export const awsActionTypes = {
  SET: 'AWS_SET',
};

export const configActionTypes = {
  ADD_IN: 'CONFIG_ACTION_ADD_IN',
  APPEND: 'CONFIG_ACTION_APPEND',
  REMOVE_FIELD_LIST_ROW: 'CONFIG_ACTION_REMOVE_FIELD_LIST_ROW',
  SET: 'CONFIG_ACTION_SIMPLE_SET',
  SET_IN: 'CONFIG_ACTION_SET_IN',
  BATCH_SET_IN: 'CONFIG_ACTION_BATCH_SET_IN',
  MERGE: 'CONFIG_ACTION_MERGE',
  RESET: 'CONFIG_ACTION_RESET',
};

export const clusterReadyActionTypes = {
  ERROR: 'clusterReadyActionTypes.ERROR',
  STATUS: 'clusterReadyActionTypes.CLUSTER_STATUS',
  NOT_READY: 'clusterReadyActionTypes.NOT_READY',
};

export const dirtyActionTypes = {
  ADD: 'DIRTY_ADD',
  CLEAN: 'DIRTY_CLEAN',
};
export const dirtyActions = {
  add: field => ({type: dirtyActionTypes.ADD, payload: field}),
  clean: field => ({type: dirtyActionTypes.CLEAN, payload: field}),
};

export const eventErrorsActionTypes = {
  ERROR: 'EVENT_ERRORS_ERROR',
};

export const loadFactsActionTypes = {
  LOADED: 'LOAD_FACTS_LOADED',
  ERROR: 'LOAD_FACTS_ERROR',
};

export const restoreActionTypes = {
  RESTORE_STATE: 'RESTORE_RESTORE_STATE',
};

export const serverActionTypes = {
  COMMIT_REQUESTED: 'COMMIT_REQUESTED',
  COMMIT_SENT: 'COMMIT_SENT',
  COMMIT_SUCCESSFUL: 'COMMIT_SUCCESSFUL',
  COMMIT_FAILED: 'COMMIT_FAILED',
};

// Commit state machine
//
// IDLE|FAILED -> REQUESTED -> WAITING -> SUCCEEDED|FAILED

export const commitPhases = {
  IDLE: 'COMMIT_IDLE',
  REQUESTED: 'COMMIT_REQUESTED',
  WAITING: 'COMMIT_WAITING',
  SUCCEEDED: 'COMMIT_SUCCEEDED',
  FAILED: 'COMMIT_FAILED',
};
export const FIELDS = {};
const FIELD_TO_DEPS = {};
export const FORMS = {};

const getField = name => {
  if (!FIELDS[name]) {
    throw new Error(`Field ${name} not found`);
  }
  return FIELDS[name];
};

// TODO (ggreer) standardize on order of params. is dispatch first or last?
export const configActions = {
  addIn: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.ADD_IN}),
  set: (payload, dispatch) => dispatch({type: configActionTypes.SET, payload}),
  setIn: (path, value, dispatch) => dispatch({payload: {path, value}, type: configActionTypes.SET_IN}),
  batchSetIn: (dispatch, payload) => {
    const values = dispatch({payload, type: configActionTypes.BATCH_SET_IN});
    payload.splice(0, payload.length - 1);
    return values;
  },
  merge: payload => dispatch => dispatch({payload, type: configActionTypes.MERGE}),

  // TODO: (kans) move below to form actions...
  removeField: (fieldListId, index) => (dispatch, getState) => {
    const fieldList = getField(fieldListId);
    dispatch({payload: {fieldListId, index}, type: configActionTypes.REMOVE_FIELD_LIST_ROW});
    fieldList.validate(dispatch, getState);
  },
  appendField: fieldListId => (dispatch, getState) => {
    const fieldList = getField(fieldListId);
    const value = _.mapValues(fieldList.rowFields, 'default');
    dispatch({payload: {path: fieldListId, value}, type: configActionTypes.APPEND});
    fieldList.validate(dispatch, getState);
  },
  refreshExtraData: fieldName => (dispatch, getState) => {
    const field = getField(fieldName);
    field.getExtraStuff(dispatch, getState, FIELDS);
  },
  updateField: (fieldName, inputValue) => (dispatch, getState) => {
    const [name, ...split] = fieldName.split('.');
    const field = getField(name);
    return field.update(dispatch, inputValue, getState, FIELDS, FIELD_TO_DEPS, split);
  },
};

export const __deleteEverything__ = () => {
  [FIELDS, FIELD_TO_DEPS, FORMS, DEFAULT_CLUSTER_CONFIG]
    .forEach(o => _.keys(o).forEach(k => delete o[k]));

  ['error', 'inFly', 'extra'].forEach(k => DEFAULT_CLUSTER_CONFIG[k] = {});

  return {type: configActionTypes.RESET};
};

export const validateFields = async (ids, getState, dispatch, updatedId, isNow) => {
  const unvisitedIds = ids;

  // Just shake the array really hard until all the nodes fall out...
  while (unvisitedIds.length) {
    // All the fields that have already had their dependencies validated
    const toVisit = unvisitedIds.filter(id => !_.intersection(unvisitedIds, FIELDS[id].dependencies).length);
    if (!toVisit.length) {
      throw new Error(`Unresolvable fields: ${unvisitedIds}`);
    }
    await Promise.all(toVisit.map(
      id => FIELDS[id].getExtraStuff(dispatch, getState, FIELDS, isNow)
        .then(() => FIELDS[id].validate(dispatch, getState, updatedId, isNow))
    ));
    _.pullAll(unvisitedIds, toVisit);
  }
};

const addDep = (field, dep) => {
  if (!FIELD_TO_DEPS[dep]) {
    FIELD_TO_DEPS[dep] = [field];
    return;
  }
  FIELD_TO_DEPS[dep].push(field);
};

export const registerForm = (form, fields) => {
  const formName = form.id;

  if (FORMS[formName]) {
    throw new Error(`form ${formName} already exists`);
  }

  FORMS[formName] = form;

  // forms can be fields too
  FIELDS[form.id] = form;

  _.each(fields, f => {
    const fieldName = f.id;
    if (f.isForm) {
      return;
    }
    if (!fieldName) {
      throw new Error(`form ${formName}: field has no name!`);
    }
    if (DEFAULT_CLUSTER_CONFIG[fieldName]) {
      throw new Error(`form ${formName}: field ${fieldName} already exists`);
    }

    DEFAULT_CLUSTER_CONFIG[fieldName] = f.default;
    FIELDS[fieldName] = f;

    _.each(f.dependencies, d => addDep(f.id, d));
  });

  // HACK to avoid figuring out the "correct" order in FIELD_TO_DEPS
  // ... peers can have deps on the same branch
  _.each(form.dependencies, d => addDep(form.id, d));
};
