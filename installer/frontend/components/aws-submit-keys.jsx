import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';

import { configActions } from '../actions';
import { validate } from '../validate';
import { AsyncSelect } from './ui';

import * as awsActions from '../aws-actions';
import { AWS_SSH, PLATFORM_TYPE } from '../cluster-config';

const {setIn} = configActions;

export const AWS_SubmitKeys = connect(
  ({eventErrors, aws, clusterConfig}) => ({
    eventErrors,
    availableSsh: aws.availableSsh,
    ssh: clusterConfig[AWS_SSH],
    platform: clusterConfig[PLATFORM_TYPE],
  }),
  dispatch => ({
    getSsh: () => dispatch(awsActions.getSsh()),
    setSsh: value => setIn(AWS_SSH, value, dispatch),
  })
)(class AWS_SubmitKeysComponent extends React.Component {
  constructor (props) {
    super(props);
    this.state = {
      error: null,
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
    </div>;
  }
});

AWS_SubmitKeys.canNavigateForward = ({clusterConfig}) => {
  const ssh = _.get(clusterConfig, AWS_SSH);
  return !validate.nonEmpty(ssh);
};
