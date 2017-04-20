import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { awsActionTypes, configActions } from '../actions';
import { validate } from '../validate';
import { AsyncSelect } from './ui';

import * as awsActions from '../aws-actions';
import { AWS_KMS, AWS_SSH, PLATFORM_TYPE } from '../cluster-config';
import { AWS } from '../platforms';

const {setIn} = configActions;

export const AWS_SubmitKeys = connect(
  ({eventErrors, aws, clusterConfig}) => ({
    eventErrors,
    availableKms: aws.availableKms,
    availableSsh: aws.availableSsh,
    createdKms: aws.createdKms,
    kms: clusterConfig[AWS_KMS],
    ssh: clusterConfig[AWS_SSH],
    platform: clusterConfig[PLATFORM_TYPE],
  }),
  dispatch => ({
    getKms: () => dispatch(awsActions.getKms()),
    createKms: () => dispatch(awsActions.createKms()),
    getSsh: () => dispatch(awsActions.getSsh()),
    setKms: value => setIn(AWS_KMS, value, dispatch),
    setSsh: value => setIn(AWS_SSH, value, dispatch),
    addKms: (kms, availableKms) => {
      const kmses = _.cloneDeep(availableKms);
      kmses.value.push(kms);
      dispatch({
        type: awsActionTypes.SET,
        payload: { availableKms: kmses},
      });
    },
  })
)(class AWS_SubmitKeysComponent extends React.Component {
  constructor (props) {
    super(props);
    this.state = {
      error: null,
      creatingKey: false,
    };
  }
  componentDidMount() {
    this.isMounted_ = true;
  }

  maybeSetState (state) {
    if (this.isMounted_) {
      this.setState(state);
    }
  }
  render () {
    const {availableKms, setKms, getKms, createKms, kms, platform} = this.props;

    const createKey = () => {
      if (this.state.creatingKey) {
        return;
      }
      this.maybeSetState({creatingKey: true});
      createKms()
      .then(v => {
        this.props.addKms(v, availableKms);
        setKms(v.value);
      })
      .catch(error => this.maybeSetState({error}))
      .then(() => this.maybeSetState({creatingKey: false}));
    };

    return <div>
      <div className="row form-group">
        <div className="col-xs-12">
          Keys are used for encryption and connection. <a href="https://coreos.com/tectonic/docs/latest/install/aws/requirements.html#ssh-key" target="_blank">Generate new keys</a> if you don't have any existing ones.
        </div>
      </div>

      <div className="row form-group">
        <div className="col-xs-12">
          <h4>SSH Keys</h4>
          <AsyncSelect
            id="sshKey"
            availableValues={this.props.availableSsh}
            disabledValue="Please select SSH Key Pair"
            value={this.props.ssh}
            onChange={v => this.props.setSsh(v)}
            onRefresh={() => this.props.getSsh()}
          />
          {
            this.props.availableSsh.error &&
            <div className="wiz-error-message">Failed to fetch AWS SSH key pairs</div>
          }
        </div>
      </div>
      { platform === AWS && <div>
        <hr />
        <div className="row form-group">
          <div className="col-xs-12">
            <h4>KMS Keys</h4>
            <p>
              A KMS (Key Management Store) key is used to encrypt secrets sent to the control plane.
            </p>
            <div className="middler">
              <button className="btn btn-default" disabled={this.state.creatingKey} onClick={() => createKey()} style={{paddingRight: 35}}>
                <i style={{visibility: this.state.creatingKey ? 'visible' : 'hidden'}} className="fa fa-spin fa-circle-o-notch" />
                &nbsp;
                Generate a New Key
              </button>
              &nbsp;&nbsp;
              <span className="text-muted">or</span>
            </div>
            {this.state.error && <p className="wiz-error-message">
              {this.state.error}
            </p>}
            <AsyncSelect
              style={{marginTop: 20}}
              id="kmsKey"
              availableValues={availableKms}
              disabledValue="Please select KMS Key Pair"
              value={kms}
              onChange={v => setKms(v)}
              onRefresh={() => getKms()}>
            </AsyncSelect>
          </div>
        </div>
      </div>
    }
    </div>;
  }
});

AWS_SubmitKeys.canNavigateForward = ({clusterConfig}) => {
  const ssh = _.get(clusterConfig, AWS_SSH);
  if (_.get(clusterConfig, PLATFORM_TYPE) !== AWS) {
    return !validate.nonEmpty(ssh);
  }
  const kms = _.get(clusterConfig, AWS_KMS);
  return !validate.nonEmpty(ssh) && !validate.nonEmpty(kms);
};
