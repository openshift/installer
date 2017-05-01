import React from 'react';
import { connect } from 'react-redux';
import { saveAs } from 'file-saver';

import { CLUSTER_NAME } from '../cluster-config';


export const SaveAssets = connect(
  ({clusterConfig}) => {
    return {
      clusterName: clusterConfig[CLUSTER_NAME],
    };
  }
)(({blob, clusterName}) => {
  if (!blob) {
    return <p>Assets not available.</p>;
  }
  const filename = `assets-${clusterName}.zip`;
  return (
    <button className="btn btn-primary wiz-giant-button" onClick={() => saveAs(blob, filename)}>
      <i className="fa fa-download"></i>&nbsp;&nbsp;Download assets
    </button>
  );
});
