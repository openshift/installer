import React from 'react';

import { ClusterInfo } from './cluster-info';
import { AWS_Tags, tagsForm } from './aws-tags';

export const AWS_ClusterInfo = () => <div>
  <ClusterInfo />
  <div className="row form-group">
    <div className="col-xs-3">
      <label htmlFor="tags">AWS Tags</label>
    </div>
    <div className="col-xs-9">
      <AWS_Tags />
    </div>
  </div>
</div>;

AWS_ClusterInfo.canNavigateForward = state => tagsForm.canNavigateForward(state) &&
  ClusterInfo.canNavigateForward(state);
