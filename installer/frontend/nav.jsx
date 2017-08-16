import _ from 'lodash';
import React from 'react';
import { withRouter } from 'react-router-dom';

const getDisplayName = Component => Component.displayName || Component.name || 'Component';

export const withNav = Wrapped => {
  const WithNav = withRouter(props => <Wrapped
    {..._.omit(props, 'history')}
    navNext={() => props.history.push(props.history.location.pathname, {next: true})}
    navPrevious={() => props.history.push(props.history.location.pathname, {previous: true})}
  />);
  WithNav.displayName = `WithNav(${getDisplayName(Wrapped)})`;
  return WithNav;
};
