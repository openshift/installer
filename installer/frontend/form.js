import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { dispatch as dispatch_ } from './store';
import { configActions, registerForm, validateFields } from './actions';
import { toError, toExtraData, toInFly, toExtraDataInFly, toExtraDataError } from './utils';
import { ErrorComponent, ConnectedFieldList } from './components/ui';
import { TectonicGA } from './tectonic-ga';
import { PLATFORM_TYPE } from './cluster-config';

const { setIn, batchSetIn, append, removeAt } = configActions;
const nop = () => undefined;

class Node {
  constructor (id, opts) {
    if (!id) {
      throw new Error('I need an id');
    }
    this.clock = 0;
    this.id = id;
    this.name = opts.name || id;
    this.validator = opts.validator || nop;
    this.dependencies = opts.dependencies || [];
    this.ignoreWhen = opts.ignoreWhen;
    this.getExtraStuff_ = opts.getExtraStuff;
  }

  getExtraStuff (dispatch, getState, FIELDS, isNow = () => true) {
    if (!this.getExtraStuff_) {
      return Promise.resolve();
    }

    const cc = getState().clusterConfig;
    if (_.some(this.dependencies, d => FIELDS[d] && !FIELDS[d].isValid(cc))) {
      // Dependencies are not all satisfied yet
      return Promise.resolve();
    }

    const inFlyPath = toExtraDataInFly(this.id);
    setIn(inFlyPath, true, dispatch);

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

  async validate (dispatch, getState, updatedId = undefined, isNow = () => true) {
    const inFlyPath = toInFly(this.id);
    setIn(inFlyPath, true, dispatch);

    const cc = getState().clusterConfig;
    const error = await this.validator(this.getData(cc), cc, updatedId, dispatch);
    if (isNow()) {
      await setIn(toError(this.id), _.isEmpty(error) ? undefined : error, dispatch);
    }
    setIn(inFlyPath, false, dispatch);
    return _.isEmpty(error);
  }

  noError (cc) {
    return _.isEmpty(_.get(cc, toError(this.id))) && _.isEmpty(_.get(cc, toExtraDataError(this.id)));
  }

  isIgnored (cc) {
    return this.ignoreWhen && !!this.ignoreWhen(cc);
  }
}

export class Field extends Node {
  constructor(id, opts = {}) {
    super(id, opts);
    if (!_.has(opts, 'default')) {
      throw new Error(`${id} needs a default`);
    }
    this.default = opts.default;
  }

  getData (cc) {
    return cc[this.id];
  }

  async update (dispatch, value, getState, FIELDS, FIELD_TO_DEPS, split) {
    // Create an isNow() function that only returns true until this.update() is called again. This allows async
    // callbacks to confirm that we are still dealing with the same Field update event.
    this.clock = this.clock + 1;
    const now = this.clock;
    const isNow = () => this.clock === now;

    setIn([this.id, ...split], value, dispatch);
    const isFieldValid = await this.validate(dispatch, getState, this.id, isNow);

    if (!isFieldValid) {
      if (getState().dirty[this.name]) {
        TectonicGA.sendEvent('Validation Error', 'user input', this.name, getState().clusterConfig[PLATFORM_TYPE]);
      }
      return;
    }

    try {
      // Recursively find all fields that are dependent on this field
      const depsDeep = id => {
        const deps = FIELD_TO_DEPS[id];
        return deps ? _.uniq([...deps, ..._.flatMap(deps, d => depsDeep(d))]) : [];
      };
      await validateFields(depsDeep(this.id), getState, dispatch, this.id, isNow);
    } catch (e) {
      console.error(e);
    }
  }

  isValid (cc) {
    return this.isIgnored(cc) || this.noError(cc);
  }

  inFly (cc) {
    return _.get(cc, toInFly(this.id)) || _.get(cc, toExtraDataInFly(this.id));
  }
}

export class Form extends Node {
  constructor(id, fields, opts = {}) {
    super(id, opts);
    this.isForm = true;
    this.fields = fields;
    this.dependencies = fields.map(f => f.id).concat(this.dependencies);

    this.errorComponent = connect(
      ({clusterConfig}) => ({error: _.get(clusterConfig, toError(id))})
    )(ErrorComponent);
    registerForm(this, fields);
  }

  isValid (cc) {
    return this.isIgnored(cc) || (this.noError(cc) && _.every(this.fields, f => f.isValid(cc)));
  }

  getData (cc) {
    return this.fields.filter(f => !f.isIgnored(cc)).reduce((acc, f) => {
      acc[f.name] = f.getData(cc);
      return acc;
    }, {});
  }

  inFly (cc) {
    return _.get(cc, toInFly(this.id)) || _.some(this.fields, f => f.inFly(cc));
  }

  get canNavigateForward () {
    return ({clusterConfig}) => !this.inFly(clusterConfig) && this.isValid(clusterConfig);
  }

  get Errors () {
    return this.errorComponent;
  }
}

const toValidator = (fields, listValidator) => (value, cc) => {
  const errs = listValidator ? listValidator(value, cc) : [];
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
      const err = validator(childValue, cc);
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
    this.validate(dispatch, getState);
  }

  remove (dispatch, i, getState) {
    removeAt(this.id, i, dispatch);
    this.validate(dispatch, getState);
  }
}
