import { connect } from 'react-redux';
import React from 'react';

import { ClusterInfo } from './cluster-info';
import { AWS_Tags, tagsForm } from './aws-tags';
import { CLUSTER_NAME } from '../cluster-config';

const TagsWithPlaceholder = connect(({clusterConfig}) => ({clusterName: clusterConfig[CLUSTER_NAME] || 'myclustername'}))(({clusterName}) => <AWS_Tags placeholder={`e.g. ${clusterName}`} />);

export const AWS_ClusterInfo = () => <div>
  <ClusterInfo />
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="tags">AWS Tags</label>
    </div>
    <div className="col-xs-9">
      <TagsWithPlaceholder />
    </div>
  </div>
  <tagsForm.Errors />
</div>;

AWS_ClusterInfo.canNavigateForward = state => tagsForm.canNavigateForward(state) &&
  ClusterInfo.canNavigateForward(state);
