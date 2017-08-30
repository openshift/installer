import _ from 'lodash';
import React from 'react';
import ReactDom from 'react-dom';
import Modal from 'react-modal';
import { connect } from 'react-redux';

import { store } from '../store';
import { eventErrorsActionTypes, restoreActionTypes, validateAllFields } from '../actions';
import { isEnabled, isSupported } from '../platforms';
import { readFile } from '../readfile';
import { TectonicGA } from '../tectonic-ga';
import { CLUSTER_NAME } from '../cluster-config';
import { LoaderInline } from './loader';

const UPLOAD_ERROR_NAME = 'UPLOAD_ERROR_NAME';

const handleUpload = (blob, cb) => dispatch => {
  TectonicGA.sendEvent('Installer Button', 'click', 'User uploads progress file');

  if (!blob) {
    return;
  }

  const uploadError = error => dispatch({
    type: eventErrorsActionTypes.ERROR,
    payload: {error, name: UPLOAD_ERROR_NAME},
  });

  readFile(blob)
    .then(result => {
      try {
        const restoreState = JSON.parse(result);

        const platform = _.get(restoreState, 'clusterConfig.platformType');
        if (!isSupported(platform)) {
          uploadError(`Platform ${platform} is not supported`);
          return;
        }
        if (!isEnabled(platform)) {
          uploadError(`Platform ${platform} is disabled`);
          return;
        }

        dispatch({
          type: restoreActionTypes.RESTORE_STATE,
          payload: restoreState,
        });
        uploadError(null);

        // The restored state may contain errors, so we don't want to use an old version.
        dispatch(validateAllFields(cb));
      } catch (e) {
        uploadError("File doesn't seem to be a saved installer state");
      }
    })
    .catch(err => {
      console.error(err);
      uploadError("Can't read installer state file");
    });
};

const stateToProps = ({eventErrors, clusterConfig}) => ({
  clusterName: clusterConfig[CLUSTER_NAME],
  uploadError: eventErrors[UPLOAD_ERROR_NAME],
});

const Modal_ = connect(stateToProps, {handleUpload})(
  class Modal_Inner extends React.Component {
    constructor(props) {
      super(props);
      this.state = {done: false, inProgress: false};
      this.onKeyDown = event => event.keyCode === 27 && this.close();
    }

    componentDidMount() {
      window.addEventListener('keydown', this.onKeyDown, true);
    }

    componentWillUnmount() {
      this.unmounted = true;
      window.removeEventListener('keydown', this.onKeyDown);
    }

    handleUpload (e) {
      if (e.target.files.length) {
        this.setState({done: false, inProgress: true});
        this.props.handleUpload(e.target.files[0], () => {
          if (!this.unmounted) {
            this.setState({done: true, inProgress: false});
          }
        });
      }
    }

    close () {
      setTimeout(() => ReactDom.unmountComponentAtNode(document.getElementById('tectonic-modal')), 0);
      if (this.state.done && this.props.cb) {
        this.props.cb();
      }
    }

    render () {
      const {clusterName, uploadError} = this.props;

      return (
        <Modal isOpen={true} className="tectonic-modal" overlayClassName="tectonic-modal-overlay" shouldCloseOnOverlayClick={false}>
          <div className="modal-header">
            <h2 className="modal-title">Restore Progress</h2>
            <p>Pre-fill all of the inputs based on a "progress file", which can be downloaded at the end of the installer. Note that future compatibility is not guaranteed.</p>
          </div>
          <div className="modal-body" style={{minHeight: 100}}>
            <input id="upload-state" type="file" onChange={this.handleUpload.bind(this)} />
            {uploadError && <p className="wiz-error-message">{uploadError}</p>}
            {!uploadError && this.state.inProgress && <span><LoaderInline /> Restoring...</span>}
            {this.state.done && <p className="alert alert-info">Restored state for {clusterName} cluster.</p>}
          </div>
          <div className="modal-footer tectonic-modal-footer">
            <button className="btn btn-default" onClick={this.close.bind(this)}>Close</button>
          </div>
        </Modal>
      );
    }
  });

export const restoreModal = (cb) => ReactDom.render(<Modal_ cb={cb} store={store} />, document.getElementById('tectonic-modal'));
