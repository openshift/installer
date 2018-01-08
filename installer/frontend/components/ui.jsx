import _ from 'lodash';
import classNames from 'classnames';
import { Set as ImmutableSet } from 'immutable';
import React from 'react';
import { connect } from 'react-redux';

import { withNav } from '../nav';
import { validate } from '../validate';
import { readFile } from '../readfile';
import { toError, toExtraData, toInFly, toExtraDataInFly, toExtraDataError } from '../utils';

import { dirtyActions, configActions } from '../actions';
import { DESELECTED_FIELDS } from '../cluster-config.js';

import { Alert } from './alert';

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

export const ErrorComponent = props => {
  const error = props.error;
  if (props.ErrorComponent) {
    return <props.ErrorComponent error={error} />;
  }
  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }
  return <span />;
};

const Field = withNav(connect(
  (state, {id}) => ({isDirty: _.get(state.dirty, id)}),
  (dispatch, {id}) => ({makeDirty: () => dispatch(dirtyActions.add(id))}),
)(props => {
  const tag = props.tag || 'input';
  const isInvalid = props.invalid && (props.forceDirty || props.isDirty);
  const fieldClasses = classNames(props.className, {'wiz-invalid': isInvalid});
  const errorClasses = classNames('wiz-error-message', {hidden: !isInvalid});

  const elementProps = {};
  Object.keys(props).filter(k => FIELD_PROPS.has(k)).forEach(k => {
    elementProps[k] = props[k];
  });

  const onEnterKeyNavigateNext = e => {
    if (e.keyCode === 13) {
      e.preventDefault();
      props.navNext();
    }
  };

  const nextProps = Object.assign({
    className: fieldClasses,
    value: props.value || '',
    autoCorrect: 'off',
    autoComplete: 'off',
    spellCheck: 'false',
    children: undefined,
    onPaste: props.makeDirty,
    onChange: e => {
      if (props.onValue) {
        props.onValue(e.target.value);
      }
    },
    onBlur: e => {
      if (e.target.value) {
        props.makeDirty();
      }
    },
    onKeyDown: ['number', 'password', 'text'].includes(props.type) ? onEnterKeyNavigateNext : undefined,
  }, elementProps);

  return (
    <div>
      {props.prefix}
      {props.renderField
        ? props.renderField(props, elementProps, fieldClasses)
        : React.createElement(tag, nextProps)
      }
      {props.suffix}
      {props.children}
      <div className={errorClasses}>
        {props.invalid}
      </div>
    </div>
  );
}));

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
        onBlur={injectedProps.makeDirty}
      />;
    };
    return <Field {...props} renderField={renderField} />;
  };
};

// component for uninteresting input[type="text"] fields.
// Handles error displays and boilerplate attributes.
// <Input id:REQUIRED invalid="error message" placeholder value onValue />
export const Input = props => <Field tag="input" type="text" {...props}>{props.children}</Field>;

// If props.validator is specified, use it to override any existing props.invalid error
export const CIDR = props => props.validator
  ? <Input {...props} invalid={props.validator(props.value)} />
  : <Input {...props} />;

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
      onBlur={injectedProps.makeDirty}
    />;
  };
  return <Field {...props} renderField={renderField} />;
};

export const CheckBox = makeBooleanField('checkbox');

export const ToggleButton = props => <button className={props.className} style={props.style} onClick={() => props.onValue(!props.value)}>
  {props.value ? 'Hide' : 'Show'}&nbsp;{props.children}
  <i style={{marginLeft: 7}} className={classNames('fa', {'fa-chevron-up': props.value, 'fa-chevron-down': !props.value})}></i>
</button>;

// A textarea/file-upload combo
// <FileArea id:REQUIRED invalid="error message" placeholder value onValue>
export const FileArea = connect(
  () => ({tag: 'textarea'}),
  (dispatch, {id}) => ({makeDirty: () => dispatch(dirtyActions.add(id))}),
)((props) => {
  const {onValue, makeDirty, uploadButtonLabel} = props;
  const handleUpload = (e) => {
    readFile(e.target.files.item(0))
      .then((value) => {
        onValue(value);
      })
      .catch((msg) => {
        console.error(msg);
      })
      .then(() => {
        makeDirty();
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
      options = [{label: value, value}].concat(options);
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
      <select id={id} value={value} disabled={disabled} style={value === '' ? {color: '#aaa'} : undefined} onChange={e => {
        makeDirty();
        onValue(e.target.value);
      }}>
        {children}
        {optionElems}
      </select>
      {invalid && isDirty &&
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
  invalid: _.get(clusterConfig, toError(field)) || _.get(clusterConfig, toExtraDataError(field)),
  isDirty: _.get(dirty, field),
  extraData: _.get(clusterConfig, toExtraData(field)),
  inFly: _.get(clusterConfig, toInFly(field)) || _.get(clusterConfig, toExtraDataInFly(field)),
});

const dispatchToProps = (dispatch, {field}) => ({
  updateField: (path, value) => dispatch(configActions.updateField(path, value)),
  makeDirty: () => dispatch(dirtyActions.add(field)),
  refreshExtraData: () => dispatch(configActions.refreshExtraData(field)),
  removeField: (i) => dispatch(configActions.removeField(field, i)),
  appendField: () => dispatch(configActions.appendField(field)),
});

class Connect_ extends React.Component {
  handleValue (v) {
    const { children, field, updateField } = this.props;
    const child = React.Children.only(children);

    updateField(field, v);
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
    const { field, value, invalid, children, isDirty, makeDirty,
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

// if undefined, default to true
const stateToIsDeselected = ({clusterConfig}, {field}) => {
  field = `${DESELECTED_FIELDS}.${field}`;
  return {
    field,
    isDeselected: !!_.get(clusterConfig, field),
  };
};

export const Deselect = connect(
  stateToIsDeselected,
  {updateField: configActions.updateField}
)(({field, isDeselected, updateField}) => <span className="deselect">
  <CheckBox id={field} value={!isDeselected} onValue={v => updateField(field, !v)} />
</span>);

export const DeselectField = connect(stateToIsDeselected)(({children, isDeselected}) => React.cloneElement(
  React.Children.only(children),
  {disabled: isDeselected, selectable: true}
));

const certPlaceholder = `Paste your certificate here. It should start with:

-----BEGIN CERTIFICATE-----

It should end with:

-----END CERTIFICATE-----`;

export const CertArea = (props) => <FileArea
  {...props}
  className="wiz-tls-asset-field"
  invalid={validate.certificate(props.value)}
  placeholder={certPlaceholder}
/>;

const privateKeyPlaceholder = `Paste your private key here. It should start with:

-----BEGIN RSA PRIVATE KEY-----

It should end with:

-----END RSA PRIVATE KEY-----`;

export const PrivateKeyArea = (props) => <FileArea
  {...props}
  className="wiz-tls-asset-field"
  invalid={validate.privateKey(props.value)}
  placeholder={privateKeyPlaceholder}
/>;

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
      options = [{label: value, value}].concat(options);
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
          {onRefresh &&
            <button className="btn btn-default" disabled={availableValues.inFly} onClick={onRefresh} title="Refresh">
              <i className={iClassNames}></i>
            </button>
          }
        </div>
        {props.invalid &&
          <div className="wiz-error-message">
            {props.invalid}
          </div>
        }
      </div>
    );
  }
}

class InnerFieldList_ extends React.Component {
  render() {
    const {value, removeField, children, fields, id} = this.props;
    const onlyChild = React.Children.only(children);
    const newChildren = _.map(value, (unused, i) => {
      const row = {};
      _.keys(fields).forEach(k => row[k] = `${id}.${i}.${k}`);
      const childProps = { row, i, key: i, remove: () => removeField(id, i) };
      if (i === value.length - 1) {
        childProps.autoFocus = true;
      }
      return React.cloneElement(onlyChild, childProps);
    });
    return <div>{newChildren}</div>;
  }
}

export const ConnectedFieldList = connect(
  ({clusterConfig}, {id}) => ({value: clusterConfig[id]}),
  (dispatch) => ({removeField: (id, i) => dispatch(configActions.removeField(id, i))})
)((props) => <InnerFieldList_ {...props} />);

export class DropdownMixin extends React.PureComponent {
  constructor(props) {
    super(props);
    this.listener = this._onWindowClick.bind(this);
    this.state = {active: !!props.active};
    this.toggle = this.toggle.bind(this);
    this.hide = this.hide.bind(this);
  }

  _onWindowClick ( event ) {
    if (!this.state.active ) {
      return;
    }

    if (event.target === this.dropdownElement || this.dropdownElement.contains(event.target)) {
      return;
    }
    this.hide();
  }

  componentDidMount () {
    window.addEventListener('click', this.listener);
  }

  componentWillUnmount () {
    window.removeEventListener('click', this.listener);
  }

  onClick_ (key, e) {
    e.stopPropagation();
    this.setState({active: false});
  }

  toggle () {
    this.setState({active: !this.state.active});
  }

  hide (e) {
    e && e.stopPropagation();
    this.setState({active: false});
  }
}

export class Dropdown extends DropdownMixin {
  render() {
    const {active} = this.state;
    const {items, header} = this.props;

    const children = _.map(items, (href, title) => {
      return <li className="tectonic-dropdown-menu-item" key={title}>
        <a className="tectonic-dropdown-menu-item__link" href={href} rel="noopener" target="_blank">{title}</a>
      </li>;
    });

    return (
      <div ref={el => this.dropdownElement = el} className="dropdown" onClick={this.toggle}>
        <a className="tectonic-dropdown-menu-title">{header}&nbsp;&nbsp;<i className="fa fa-angle-down"></i></a>
        <ul className="dropdown-menu tectonic-dropdown-menu" style={{display: active ? 'block' : 'none'}}>
          {children}
        </ul>
      </div>
    );
  }
}

export class DropdownInline extends DropdownMixin {
  render() {
    const {active} = this.state;
    const {items, header} = this.props;

    return (
      <div ref={el => this.dropdownElement = el} className="dropdown" onClick={this.toggle} style={{display: 'inline-block'}}>
        <a>{header}&nbsp;&nbsp;<i className="fa fa-caret-down"></i></a>
        <ul className="dropdown-menu--dark" style={{display: active ? 'block' : 'none'}}>
          {items.map(([title, cb], i) => <li className="dropdown-menu--dark__item" key={i} onClick={cb}>{title}</li>)}
        </ul>
      </div>
    );
  }
}
