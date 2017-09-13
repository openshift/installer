import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';

import { commitToServer } from '../server';
import { commitPhases } from '../actions';
import { withNav } from '../nav';

export const SubmitDefinition = withNav(connect(
  state => ({
    phase: state.commitState.phase,
    response: state.commitState.response,
    ready: state.cluster.ready,
  }),
  {onFinish: commitToServer}
)(({navPrevious, phase, response, ready, onFinish}) => {
  const inProgress = (phase === commitPhases.REQUESTED ||
                      phase === commitPhases.WAITING ||
                      phase === commitPhases.SUCCEEDED && !ready);

  let btn = <button className="btn btn-primary wiz-giant-button" onClick={() => onFinish(false)}>Submit</button>;
  if (inProgress) {
    btn = <button className="btn btn-primary wiz-giant-button disabled">
      <i className="fa fa-spin fa-circle-o-notch"></i> Checking Your Definition...
    </button>;
  }

  const errorMessage = response ? response.toString() : '';
  const errorClasses = classNames('wiz-error-message', {
    hidden: phase !== commitPhases.FAILED,
  });

  return (
    <div>
      <p>
        Congratulations! Your cluster has been defined and will be submitted to Terraform. After submission, the definition cannot be updated. Go <a onClick={!inProgress && navPrevious} className={inProgress && 'disabled'}>back</a> to update or make changes.
      </p>
      <p>
        {/* eslint-disable react/jsx-no-target-blank */}
        You'll be able to download your <a href="https://coreos.com/tectonic/docs/latest/admin/assets-zip.html" rel="noopener" target="_blank">assets zip file</a> after the definition is submitted.
        {/* eslint-enable react/jsx-no-target-blank */}
      </p>
      <br />
      <div className="wiz-giant-button-container">
        {btn}
      </div>
      <p>
        <b>Advanced mode: </b>
        <a onClick={() => onFinish(true)}>Manually boot</a> your own cluster. Validate configuration, generate assets, but don't create the cluster.
      </p>
      <div className={errorClasses}>
        {errorMessage}
      </div>
    </div>
  );
}));
