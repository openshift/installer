import _ from 'lodash';
import classNames from 'classnames';
import { Set as ImmutableSet } from 'immutable';
import React from 'react';
import { connect } from 'react-redux';

import { validate } from '../validate';
import { readFile } from '../readfile';
import { toError, toAsyncError, toExtraData, toInFly, toExtraDataInFly, toExtraDataError } from '../utils';

import { configActionTypes, dirtyActionTypes, configActions } from '../actions';
import { DESELECTED_FIELDS } from '../cluster-config.js';


// Use this function to dirty a field due to
// non-user interaction (like uploading a config file)
export const markIDDirty = (dispatch, id) => {
  if (!id) {
    throw new Error('ID required!');
  }
  dispatch({
    type: dirtyActionTypes.ADD,
    payload: id,
  });
};

// Taken more-or-less from https://www.w3.org/TR/html5/forms.html#the-input-element
const FIELD_PROPS = ImmutableSet([
  'accept',
  'accesskey',
  'alt',
  'autocomplete',
  'autoFocus', // Funny case required by react input management magic
  'checked',
  'dirname',
  'disabled',
  'form',
  'formaction',
  'formenctype',
  'formmethod',
  'formnovalidate',
  'formtarget',
  'height',
  'hidden',
  'id',
  'inputmode',
  'lang',
  'list',
  'max',
  'maxlength',
  'min',
  'minlength',
  'multiple',
  'name',
  'pattern',
  'placeholder',
  'onChange',
  'readonly',
  'required',
  'size',
  'spellcheck',
  'src',
  'step',
  'style',
  'tabindex',
  'title',
  'type',
  'width',
]);

export const ErrorComponent = props => props.error ? <div className="wiz-error-message">{props.error}</div> : <span/>;

const Field = connect(
  (state, {id}) => ({isDirty: _.get(state.dirty, id)}),
  dispatch => ({
    markDirty: id => markIDDirty(dispatch, id),
    markClean: id => dispatch({type: dirtyActionTypes.CLEAN, payload: id }),
  })
)(class Field extends React.Component {
  componentWillUnmount() {
    if (this.props.autoClean) {
      this.props.markClean(this.props.id);
    }
  }
  render () {
    const props = this.props;
    const tag = props.tag || 'input';
    const dirty = props.forceDirty || props.isDirty;
    const fieldClasses = classNames(props.className, {
      'wiz-dirty': dirty,
      'wiz-invalid': props.invalid,
    });
    const errorClasses = classNames('wiz-error-message', {
      hidden: !(dirty && props.invalid),
    });

    const elementProps = {};
    Object.keys(props).filter(k => FIELD_PROPS.has(k)).forEach(k => {
      elementProps[k] = props[k];
    });

    const nextProps = Object.assign({
      className: fieldClasses,
      value: props.value || '',
      autoCorrect: 'off',
      autoComplete: 'off',
      spellCheck: 'false',
      children: undefined,
      onPaste: () => props.markDirty(props.id),
      onChange: e => {
        if (props.onValue) {
          props.onValue(e.target.value);
        }
      },
      onBlur: e => {
        if (props.blurry || e.target.value) {
          props.markDirty(props.id);
        }
      },
    }, elementProps);

    return (
      <div>
        {props.prefix}
        {props.renderField
          ? props.renderField(props, elementProps, fieldClasses)
          : React.createElement(tag, nextProps)
        }
        {props.suffix && <span>&nbsp;&nbsp;{props.suffix}</span>}
        {props.children}
        <div className={errorClasses}>
          {props.invalid}
        </div>
      </div>
    );
  }
});

const makeBooleanField = type => {
  return function booleanField(props) {
    const renderField = (injectedProps, cleanedProps, classes) => {
      return <input type={type} checked={injectedProps.inverted ? !injectedProps.value : injectedProps.value} className={classes} {...cleanedProps}
        onChange={e => {
          const value = injectedProps.inverted ? !e.target.checked : !!e.target.checked;
          injectedProps.onValue(value);
          if (props.onChange) {
            props.onChange(value);
          }
        }}
        onBlur={() => injectedProps.markDirty(injectedProps.id)}
      />;
    };
    return <Field {...props} renderField={renderField} />;
  };
};

// component for uninteresting input[type="text"] fields.
// Handles error displays and boilerplate attributes.
// <Input id:REQUIRED invalid="error message" placeholder value onValue />
export const Input = props => <Field tag="input" type="text" {...props}>{props.children}</Field>;
export const NumberInput = props => <Field tag="input" type="number" onChange={e => {
  const number = parseInt(e.target.value, 10);
  props.onValue(isNaN(number) ? 0 : number);
}} {...props} />;
export const Password = props => <Field tag="input" type="password" {...props} />;
export const RadioBoolean = makeBooleanField('radio');
export const Radio = props => {
  const renderField = (injectedProps, cleanedProps, classes) => {
    return <input type="radio" className={classes} {...cleanedProps}
      onChange={() => {
        injectedProps.onValue(props.value);
        if (props.onChange) {
          props.onChange(props.value);
        }
      }}
      onBlur={() => injectedProps.markDirty(injectedProps.id)}
    />;
  };
  return <Field {...props} renderField={renderField} />;
};
export const CheckBox = makeBooleanField('checkbox');
export const ToggleButton = props => <button className={props.className} style={props.style} onClick={() => props.onValue(!props.value)}>
  {props.value ? 'Hide' : 'Show'}&nbsp;{props.children}
  <i style={{marginLeft: 7}} className={classNames("fa", {"fa-chevron-up": props.value, "fa-chevron-down": !props.value})}></i>
</button>;

// A textarea/file-upload combo
// <FileArea id:REQUIRED invalid="error message" placeholder value onValue>
export const FileArea = connect(
  () => {
    return {
      tag: 'textarea',
    };
  },
  (dispatch) => {
    return {
      markDirtyUpload: (id) => {
        markIDDirty(dispatch, id);
      },
    };
  }
)((props) => {
  const {id, onValue, markDirtyUpload, uploadButtonLabel} = props;
  const handleUpload = (e) => {
    readFile(e.target.files.item(0))
    .then((value) => {
      onValue(value);
    })
    .catch((msg) => {
      console.error(msg);
    })
    .then(() => {
      markDirtyUpload(id);
    });
    // Reset value so that onChange fires if you pick the same file again.
    e.target.value = null;
  };

  return (
    <div>
      <label className="btn btn-sm btn-link">
        <span className="fa fa-upload"></span>&nbsp;&nbsp;{uploadButtonLabel || 'Upload'} {' '}
        <input style={{display: 'none'}}
               type="file"
               onChange={handleUpload} />
      </label>
      <Field {...props} />
    </div>
  );
});

// <Select id:REQUIRED value onValue>
//   <option....>
// </Select>
export const Select = ({id, children, value, onValue, invalid, isDirty, makeDirty, availableValues, className, disabled, style}) => {
  const optionElems = [];
  if (availableValues) {
    let options = availableValues.value;
    if (value && !options.map(r => r.value).includes(value)) {
      options = [{label: value, value: value}].concat(options);
    }

    const optgroups = new Map();

    _.each(options, o => {
      const elem = <option key={o.value} value={o.value}>{o.label}</option>;
      if (!o.optgroup) {
        optionElems.push(elem);
        return;
      }

      if (!optgroups.get(o.optgroup)) {
        optgroups.set(o.optgroup, []);
      }
      optgroups.get(o.optgroup).push(elem);
    });

    optgroups.forEach((child, label) => optionElems.push(<optgroup key={label} label={label}>{child}</optgroup>));
  }

  return (
    <div className={className} style={style}>
      <select id={id} value={value} disabled={disabled} onChange={e => {
        makeDirty();
        onValue(e.target.value);
      }}>
        {children}
        {optionElems}
      </select>
      { invalid && isDirty &&
        <div className="wiz-error-message">
          {invalid}
        </div>
      }
    </div>
  );
};

export const Selector = props => {
  const value = props.value;
  const options = _.get(props, 'extraData.options', []);

  if (value && !options.map(r => r.value).includes(value)) {
    options.splice(0, 0, {value, label: value});
  }

  const optionsElems = options.map(o => <option value={o.value} key={o.value}>{o.label}</option>);
  if (props.disabledValue) {
    optionsElems.splice(0, 0, <option disabled={true} key="disabled" value="">{props.disabledValue}</option>);
  }

  const style = Object.assign({}, props.style || {});
  if (!props.refreshBtn) {
    style.marginRight = 0;
  }
  const iClassNames = classNames('fa', 'fa-refresh', {
    'fa-spin': props.inFly,
  });

  return <div className="async-select">
    <Select className="async-select--select" {...props} style={style}>{optionsElems}</Select>
    {props.refreshBtn && <button className="btn btn-default" onClick={props.refreshExtraData} title="Refresh">
      <i className={iClassNames}></i>
    </button>}
  </div>;
};

const stateToProps = ({clusterConfig, dirty}, {field}) => ({
  value: _.get(clusterConfig, field),
  invalid: _.get(clusterConfig, toError(field))
    || _.get(clusterConfig, toAsyncError(field))
    || _.get(clusterConfig, toExtraDataError(field)),
  isDirty:  _.get(dirty, field),
  extraData: _.get(clusterConfig, toExtraData(field)),
  inFly: _.get(clusterConfig, toInFly(field)) || _.get(clusterConfig, toExtraDataInFly(field)),
});

const dispatchToProps = (dispatch, {field}) => ({
  setField: (path, value, invalid) => dispatch(configActions.updateField(path, value, invalid)),
  makeDirty: () => markIDDirty(dispatch, field),
  makeClean: () => dispatch({type: dirtyActionTypes.CLEAN, payload: field }),
  refreshExtraData: () => dispatch(configActions.refreshExtraData(field)),
  removeField: (i) => dispatch(configActions.removeField(field, i)),
  appendField: () => dispatch(configActions.appendField(field)),
});

class Connect_ extends React.Component {
  handleValue (v) {
    const { children, field, setField } = this.props;
    const child = React.Children.only(children);

    setField(field, v);
    if (child.props.onValue) {
      child.props.onValue(v);
    }
  }
  componentDidMount () {
    const { getDefault } = this.props;

    if (_.isFunction(getDefault)) {
      this.handleValue(getDefault());
    }
  }

  render () {
    const { field, value, invalid, children, isDirty, makeDirty, makeClean,
      extraData, refreshExtraData, inFly, removeField, appendField,
    } = this.props;

    const child = React.Children.only(children);
    const id = child.props.id || field;

    const props = {
      extraData,
      id,
      inFly,
      invalid,
      isDirty,
      makeClean,
      makeDirty,
      refreshExtraData,
      removeField,
      appendField,
      onValue: v => this.handleValue(v),
    };

    switch (child.type) {
    case Radio:
      props.checked = child.props.value === value;
      break;
    case Select:
    case Selector:
      props.value = value || '';
      break;
    default:
      props.value = value;
      break;
    }

    return React.cloneElement(child, props);
  }
}

export const Connect = connect(stateToProps, dispatchToProps)(Connect_);

// <WithClusterConfig field="banana" validator={isBanana}>
//     <Input id="jones" placeholder="Enter your banana" />
// </WithClusterConfig>
//
//
export const WithClusterConfig = connect(
  ({clusterConfig, dirty}, {field}) => {
    return {
      value: _.get(clusterConfig, field),
      invalid: _.get(clusterConfig, toError(field)),
      isDirty: _.get(dirty, field),
    };
  },
  (dispatch, {field}) => {
    return {
      setField: (k, v) => {
        dispatch({
          type: configActionTypes.SET_IN,
          payload: {value: v, path: k},
        });
      },
      makeDirty: () => markIDDirty(dispatch, field),
      makeClean: () => dispatch({type: dirtyActionTypes.CLEAN, payload: field }),
    };
  }
)(class WithClusterConfig extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      inFly: false,
      syncInvalid: undefined,
      asyncInvalid: undefined,
    };
    const field = props.field;
    this.fieldName = _.isArray(field) ? field.join('.') : field;
    this.isMounted_ = false;
    this.safeSetState = (...args) => {
      if (!this.isMounted_) {
        return;
      }
      return this.setState(...args);
    };
  }

  get invalidKey () {
    return toError(this.fieldName);
  }

  componentDidMount () {
    this.isMounted_ = true;
    this.asyncValidate(this.props);
  }

  componentWillMount() {
    if (!this.props.default) {
      return;
    }
    const { field, setField, value } = this.props;
    if (value) {
      return;
    }
    setField(field, this.props.default);
  }

  componentWillUnmount() {
    this.isMounted_ = false;
  }

  setValidation () {
    const { syncInvalid, asyncInvalid } = this.state;
    this.props.setField(this.invalidKey, syncInvalid || asyncInvalid);
  }

  syncValidate ({value, validator}) {
    if (!validator) {
      return;
    }
    const syncInvalid = validator(value);
    this.safeSetState({syncInvalid}, () => this.setValidation());
    return syncInvalid;
  }

  asyncValidate ({value, asyncValidator}) {
    if (!asyncValidator) {
      return;
    }
    if (this.state.syncInvalid) {
      return;
    }
    this.safeSetState({inFly: true});
    asyncValidator(value)
    .then(() => {
      this.safeSetState({asyncInvalid: undefined}, () => this.setValidation());
    })
    .catch(err => {
      this.safeSetState({asyncInvalid: err}, () => this.setValidation());
    })
    .then(() => this.safeSetState({inFly: false}));
  }

  componentWillReceiveProps (nextProps) {
    if (nextProps.value === this.props.value && nextProps.invalid === this.props.invalid) {
      return;
    }
    this.syncValidate(nextProps);
    this.asyncValidate(nextProps);
  }

  render() {
    const { field, setField, invalid, isDirty, makeDirty, makeClean, value } = this.props;
    const child = React.Children.only(this.props.children);
    const id = child.props.id || this.fieldName;

    const props = {
      id,
      invalid,
      isDirty,
      makeClean,
      makeDirty,
      onValue: (v) => {
        setField(field, v);
        if (child.props.onValue) {
          child.props.onValue(v);
        }
      },
    };

    switch (child.type) {
    case Radio:
      props.checked = child.props.value === value;
      break;
    case Select:
      props.value = value || '';
      break;
    default:
      props.value = value;
      break;
    }

    return React.cloneElement(child, props);
  }
});

// if undefined, default to true
const stateToIsDeselected = ({clusterConfig}, {field}) => {
  field = `${DESELECTED_FIELDS}.${field}`;
  return {
    field: field,
    isDeselected: !!_.get(clusterConfig, field),
  };
};

export const Deselect = connect(
  stateToIsDeselected,
  dispatch => ({
    setField: (k, v) => {
      dispatch({
        type: configActionTypes.SET_IN,
        payload: {value: v, path: k},
      });
    },
  })
)(({field, isDeselected, setField}) => <span className="deselect">
  <CheckBox id={field} value={!isDeselected} onValue={v => setField(field, !v)}/>
</span>);

export const DeselectField = connect(stateToIsDeselected)(({children, isDeselected}) => React.cloneElement(React.Children.only(children), {disabled: isDeselected, selectable: true}));

const certPlaceholder = `Paste your certificate here. It should start with:

-----BEGIN CERTIFICATE-----

It should end with:

-----END CERTIFICATE-----`;

export const CertArea = (props) => {
  const invalid = validate.certificate(props.value);
  const areaProps = Object.assign({}, props, {
    className: props.className + ' wiz-tls-asset-field',
    invalid: invalid,
    placeholder: certPlaceholder,
  });
  return <FileArea {...areaProps} />;
};

const privateKeyPlaceholder = `Paste your private key here. It should start with:

-----BEGIN RSA PRIVATE KEY-----

It should end with:

-----END RSA PRIVATE KEY-----`;

export const PrivateKeyArea = (props) => {
  const invalid = validate.privateKey(props.value);
  const areaProps = Object.assign({}, props, {
    className: props.className + ' wiz-tls-asset-field',
    invalid: invalid,
    placeholder: privateKeyPlaceholder,
  });
  return <FileArea {...areaProps} />;
};

export const WaitingLi = ({done, error, children, substep}) => {
  const progressClasses = classNames({
    'wiz-launch-progress__step': !substep,
    'wiz-launch-progress__substep': substep,
    'wiz-error-fg': error,
    'wiz-success-fg': done && !error,
    'wiz-running-fg': !done && !error,
  });
  const iconClasses = classNames('fa', 'fa-fw', {
    'fa-exclamation-circle': error,
    'fa-check-circle': done && !error,
    'fa-spin fa-circle-o-notch': !done && !error,
  });

  return <li className={progressClasses}>
    <i className={iconClasses}></i>&nbsp;{children}
  </li>;
};

export class AsyncSelect extends React.Component {
  componentDidMount () {
    const { onChange, onRefresh, value } = this.props;
    onRefresh && onRefresh();
    value && onChange && onChange(value);
  }

  render () {
    const {id, availableValues, disabledValue, value, onChange, onRefresh} = this.props;
    const iClassNames = classNames('fa', 'fa-refresh', {
      'fa-spin': availableValues.inFly,
    });

    let options = availableValues.value;
    if (value && !options.map(r => r.value).includes(value)) {
      options = [{label: value, value: value}].concat(options);
    }
    const style = {};
    if (!onRefresh) {
      style.marginRight = 0;
    }

    const optionElems = [];
    const optgroups = new Map();

    _.each(options, o => {
      const elem = <option key={o.value} value={o.value}>{o.label}</option>;
      if (!o.optgroup) {
        optionElems.push(elem);
        return;
      }

      if (!optgroups.get(o.optgroup)) {
        optgroups.set(o.optgroup, []);
      }
      optgroups.get(o.optgroup).push(elem);
    });

    optgroups.forEach((children, label) => optionElems.push(<optgroup key={label} label={label}>{children}</optgroup>));

    const props = this.props;

    return (
      <div>
        <div className={classNames('async-select', props.className)} style={props.style}>
          {props.children}
          <select style={style}
              id={id}
              className="async-select--select"
              value={value}
              disabled={availableValues.inFly}
              onChange={e => {
                const v = e.target.value;
                props.onValue && props.onValue(v);
                onChange && onChange(v);
              }}>
            {disabledValue && <option value="" disabled>{disabledValue}</option>}
            {optionElems}
          </select>
          {
            onRefresh &&
            <button className="btn btn-default" disabled={availableValues.inFly} onClick={onRefresh} title="Refresh">
              <i className={iClassNames}></i>
            </button>
          }
        </div>
        { props.invalid &&
          <div className="wiz-error-message">
            {props.invalid}
          </div>
        }
      </div>
    );
  }
}
