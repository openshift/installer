import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';

import { commitToServer } from '../server';
import { commitPhases } from '../actions';
import { PLATFORM_TYPE } from '../cluster-config';
import { AWS, BARE_METAL, isTerraform } from '../platforms';

export const SubmitDefinition = connect(
  state => ({
    phase: state.commitState.phase,
    response: state.commitState.response,
    ready: state.cluster.ready,
    platformType: state.clusterConfig[PLATFORM_TYPE],
  }),
  dispatch => ({
    onFinish: (dryRun=false) => dispatch(commitToServer(dryRun)),
  })
)(({phase, platformType, response, ready, onFinish, navigatePrevious, navigateNext}) => {
  let feature =
    <div className="wiz-giant-button-container">
      <button className="btn btn-primary wiz-giant-button"
              onClick={() => onFinish(false)}>
        Submit
      </button>
    </div>;

  let pager = '';

  const inProgress = (phase === commitPhases.REQUESTED ||
                      phase === commitPhases.WAITING ||
                      phase === commitPhases.SUCCEEDED && !ready);

  if (inProgress) {
    feature =
      <div className="wiz-giant-button-container">
        <button className="btn btn-primary wiz-giant-button disabled">
          <i className="fa fa-spin fa-circle-o-notch"></i>{' '}
          Checking Your Definition...
        </button>
      </div>;
  }

  if (phase === commitPhases.SUCCEEDED && ready) {
    feature = (
      <div>
        <div className="wiz-herotext wiz-herotext--success">
          <span className="fa fa-check-circle wiz-herotext-icon"></span> {' '}
          High Fives! <br /> Your matchbox server was configured successfully!
        </div>
      </div>
    );
    pager = (
      <div className="wiz-form__actions">
        <button className="btn btn-primary wiz-form__actions__next"
           onClick={navigateNext}
           >Next Step</button>
      </div>
    );
  }

  const errorMessage = response ? response.toString() : '';
  const errorClasses = classNames('wiz-error-message', {
    hidden: phase !== commitPhases.FAILED,
  });

  let msg;
  if (isTerraform(platformType)) {
    msg = <span>Congratulations! Your cluster has been defined and will be submitted to Terraform.</span>;
  } else if (platformType === AWS) {
    msg = <span>Congratulations! Your cluster has been defined and will be submitted to AWS.</span>;
  } else if (platformType === BARE_METAL) {
    msg = <span>Now that you've defined your cluster, it will be submitted to your <a href="https://github.com/coreos/matchbox/blob/master/Documentation/matchbox.md"
             target="_blank">matchbox</a> service.</span>;
  }

  return (
    <div>
      <p>
        { msg }
        <span> After submission, the definition cannot be updated. Go <a onClick={!inProgress && navigatePrevious} className={inProgress && 'disabled'}>back</a> to update or make changes.</span>
      </p>
      <p>
        You'll be able to download your assets zip file after the definition is submitted. These assets will be used in step 3.
      </p>
      <br />
      {feature}
      <p>
        <b>Advanced mode: </b>
        <a onClick={() => onFinish(true)} href="#">Manually boot</a> your own cluster. Validate configuration, generate assets, but don't create the cluster.
      </p>
      <div className={errorClasses}>
        {errorMessage}
      </div>
      {pager}
    </div>
  );
});
SubmitDefinition.canNavigateForward = ({cluster}) => {
  return cluster.ready;
};
