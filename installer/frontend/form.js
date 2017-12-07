import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { dispatch as dispatch_ } from './store';
import { configActions, registerForm } from './actions';
import { toError, toIgnore, toExtraData, toInFly, toExtraDataInFly, toExtraDataError } from './utils';
import { ErrorComponent, ConnectedFieldList } from './components/ui';
import { TectonicGA } from './tectonic-ga';
import { PLATFORM_TYPE } from './cluster-config';

const { setIn, batchSetIn, append, removeAt } = configActions;
const nop = () => undefined;

let clock_ = 0;

class Node {
  constructor (id, opts) {
    if (!id) {
      throw new Error('I need an id');
    }
    this.clock_ = 0;
    this.id = id;
    this.name = opts.name || id;
    this.validator = opts.validator || nop;
    this.dependencies = opts.dependencies || [];
    this.ignoreWhen_ = opts.ignoreWhen;
    this.getExtraStuff_ = opts.getExtraStuff;
  }

  updateClock (now) {
    return this.clock_ = Math.max(now || clock_, this.clock_);
  }

  get isNow () {
    const now = this.clock_;
    return () => this.clock_ === now;
  }

  getExtraStuff (dispatch, cc, FIELDS, now) {
    if (!this.getExtraStuff_) {
      return Promise.resolve();
    }

    if (_.some(this.dependencies, d => FIELDS[d] && !FIELDS[d].isValid(cc))) {
      // Dependencies are not all satisfied yet
      return Promise.resolve();
    }

    this.updateClock(now);

    const inFlyPath = toExtraDataInFly(this.id);

    setIn(inFlyPath, true, dispatch);
    const isNow = this.isNow;

    return this.getExtraStuff_(dispatch, isNow).then(data => {
      if (!isNow()) {
        return;
      }
      batchSetIn(dispatch, [
        [inFlyPath, undefined],
        [toExtraData(this.id), data],
        [toExtraDataError(this.id), undefined],
      ]);
    }, e => {
      if (!isNow()) {
        return;
      }
      batchSetIn(dispatch, [
        [inFlyPath, undefined],
        [toExtraData(this.id), undefined],
        [toExtraDataError(this.id), e.message || e.toString()],
      ]);
    });
  }

  async validate (dispatch, getState, oldCC) {
    const id = this.id;
    const clusterConfig = getState().clusterConfig;
    const value = this.getData(clusterConfig);

    const syncErrorPath = toError(id);
    const inFlyPath = toInFly(id);

    const oldValue = this.getData(oldCC);

    const batches = [];

    if (_.get(clusterConfig, inFlyPath)) {
      batches.push([inFlyPath, false]);
    }

    console.debug(`validating ${this.name}`);
    const syncError = await this.validator(value, clusterConfig, oldValue, oldCC, dispatch);
    if (!_.isEmpty(syncError)) {
      console.info(`sync error ${this.name}: ${JSON.stringify(syncError)}`);
      batches.push([syncErrorPath, syncError]);
      batchSetIn(dispatch, batches);
      return false;
    }

    const oldError = _.get(oldCC, syncErrorPath);
    if (!_.isEmpty(oldError)) {
      batches.push([syncErrorPath, undefined]);
      batchSetIn(dispatch, batches);
    }

    const isValid = this.isValid(getState().clusterConfig);
    if (!isValid) {
      batchSetIn(dispatch, batches);
      return false;
    }

    batchSetIn(dispatch, batches);
    return true;
  }

  ignoreWhen (dispatch, clusterConfig) {
    if (!this.ignoreWhen_) {
      return false;
    }
    const value = !!this.ignoreWhen_(clusterConfig);
    console.debug(`ignoring ${this.id} value ${value}`);
    setIn(toIgnore(this.id), value, dispatch);
    return value;
  }

  isIgnored (clusterConfig) {
    return _.get(clusterConfig, toIgnore(this.id));
  }
}

async function promisify (dispatch, getState, oldCC, now, deps, FIELDS) {
  const { clusterConfig } = getState();

  // TODO: (kans) earlier return [] if not now?
  const promises = deps.map(field => {
    const { id } = field;
    field.ignoreWhen(dispatch, clusterConfig);
    return field.getExtraStuff(dispatch, clusterConfig, FIELDS, now)
      .then(() => field.validate(dispatch, getState, oldCC))
      .then(res => {
        if (!res) {
          console.debug(`${id} is invalid`);
        } else {
          console.debug(`${id} is valid`);
        }
        return res && id;
      }).catch(err => {
        console.error(err);
      });
  });

  return await Promise.all(promises).then(p => p.filter(id => id));
}

export class Field extends Node {
  constructor(id, opts = {}) {
    super(id, opts);
    if (!_.has(opts, 'default')) {
      throw new Error(`${id} needs a default`);
    }
    this.default = opts.default;
  }

  getData (clusterConfig) {
    return clusterConfig[this.id];
  }

  async update (dispatch, value, getState, FIELDS, FIELD_TO_DEPS, split) {
    const oldCC = getState().clusterConfig;

    const now = ++ clock_;

    let id = this.id;
    if (split && split.length) {
      id = `${id}.${split.join('.')}`;
    }

    console.info(`updating ${this.name}`);
    // TODO: (kans) - We need to lock the entire validation chain, not just validate proper
    setIn(id, value, dispatch);

    const isValid = await this.validate(dispatch, getState, oldCC);

    if (!isValid) {
      const dirty = getState().dirty;
      if (dirty[this.name]) {
        TectonicGA.sendEvent('Validation Error', 'user input', this.name, oldCC[PLATFORM_TYPE]);
      }

      console.debug(`${this.name} is invalid`);
      return;
    }

    const visited = new Set();
    const toVisit = [FIELD_TO_DEPS[this.id]];

    if (!toVisit[0].length) {
      console.debug(`no deps for ${this.name}`);
      return;
    }

    while (toVisit.length) {
      const deps = toVisit.splice(0, 1)[0];
      // TODO: check for relationship between deps
      const nextDepIDs = await promisify(dispatch, getState, oldCC, now, deps, FIELDS);
      nextDepIDs.forEach(depID => {
        const nextDeps = _.filter(FIELD_TO_DEPS[depID], d => !visited.has(d.id));
        if (!nextDeps.length) {
          return;
        }
        nextDeps.forEach(d => visited.add(d.id));
        toVisit.push(nextDeps);
      });
    }

    console.info(`finish validating ${this.name} ${isValid}`);
  }

  isValid (cc) {
    if (_.get(cc, toIgnore(this.id))) {
      return true;
    }
    const value = _.get(cc, this.id);
    return value !== '' &&
      value !== undefined &&
      _.isEmpty(_.get(cc, toError(this.id))) &&
      _.isEmpty(_.get(cc, toExtraDataError(this.id)));
  }

  inFly (clusterConfig) {
    return _.get(clusterConfig, toInFly(this.id)) || _.get(clusterConfig, toExtraDataInFly(this.id));
  }
}

export class Form extends Node {
  constructor(id, fields, opts = {}) {
    super(id, opts);
    this.isForm = true;
    this.fields = fields;
    this.fieldIDs = fields.map(f => f.id);

    this.dependencies = [...this.fieldIDs].concat(this.dependencies);

    this.errorComponent = connect(
      ({clusterConfig}) => ({error: _.get(clusterConfig, toError(id))})
    )(ErrorComponent);
    registerForm(this, fields);
  }

  isValid (clusterConfig) {
    const ignore = _.get(clusterConfig, toIgnore(this.id));
    if (ignore) {
      return true;
    }

    if (_.get(clusterConfig, toError(this.id))) {
      return false;
    }

    const invalidFields = this.fields.filter(field => !field.isValid(clusterConfig));
    return invalidFields.length === 0;
  }

  getData (clusterConfig) {
    return this.fields.filter(f => !f.isIgnored(clusterConfig)).reduce((acc, f) => {
      acc[f.name] = f.getData(clusterConfig);
      return acc;
    }, {});
  }

  inFly (clusterConfig) {
    return _.get(clusterConfig, toInFly(this.id)) || _.some(this.fields, f => f.inFly(clusterConfig));
  }

  get canNavigateForward () {
    return ({clusterConfig}) => !this.inFly(clusterConfig) && this.isValid(clusterConfig);
  }

  get Errors () {
    return this.errorComponent;
  }
}

const toValidator = (fields, listValidator) => (value, clusterConfig, oldValue) => {
  const errs = listValidator ? listValidator(value, clusterConfig, oldValue) : [];
  if (errs && !_.isObject(errs)) {
    throw new Error(`FieldLists validator must return an Array-like Object, not:\n${errs}`);
  }
  _.each(value, (child, i) => {
    errs[i] = errs[i] || {};
    _.each(child, (childValue, name) => {
      // TODO: check that the name is in the field...
      const validator = _.get(fields, [name, 'validator']);
      if (!validator) {
        return;
      }
      const err = validator(childValue, clusterConfig, _.get(oldValue, [i, name]));
      if (!err) {
        return;
      }
      errs[i][name] = err;
    });
  });

  return _.every(errs, err => _.isEmpty(err)) ? {} : errs;
};

const toDefaultOpts = opts => {
  const default_ = {};

  _.each(opts.fields, (v, k) => {
    default_[k] = v.default;
  });

  return Object.assign({}, opts, {default: [default_], validator: toValidator(opts.fields, opts.validator)});
};

export class FieldList extends Field {
  constructor(id, opts = {}) {
    super(id, toDefaultOpts(opts));
    this.fields = opts.fields;
  }

  get Map () {
    if (this.OuterListComponent_) {
      return this.OuterListComponent_;
    }
    const id = this.id;
    const fields = this.fields;

    this.OuterListComponent_ = function Outer ({children}) {
      return React.createElement(ConnectedFieldList, {id, fields}, children);
    };

    return this.OuterListComponent_;
  }

  get addOnClick () {
    return () => dispatch_(configActions.appendField(this.id));
  }

  get NonFieldErrors () {
    if (this.errorComponent_) {
      return this.errorComponent_;
    }

    const id = this.id;

    this.errorComponent_ = connect(
      ({clusterConfig}) => ({error: _.get(clusterConfig, toError(id), {})}),
    )(({error}) => React.createElement(ErrorComponent, {error: error[-1]}));

    return this.errorComponent_;
  }

  append (dispatch, getState) {
    const child = {};
    _.each(this.fields, (f, name) => {
      child[name] = _.cloneDeep(f.default);
    });
    append(this.id, child, dispatch);
    this.validate(dispatch, getState, getState().clusterConfig);
  }

  remove (dispatch, i, getState) {
    removeAt(this.id, i, dispatch);
    this.validate(dispatch, getState, getState().clusterConfig);
  }
}
