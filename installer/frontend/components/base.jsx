import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { Route, Switch, withRouter } from 'react-router-dom';

import { withNav } from '../nav';
import { trail, trailSections } from '../trail';

import { Loader } from './loader';
import { Header } from './header';
import { Footer } from './footer';

const NavSection = withRouter(({canNavigateForward, currentPage, history, sections, navTrail, title}) => {
  const section = _.find(navTrail.sections, s => sections.includes(s) && s.includes(currentPage));
  let disabled = false;

  return <ul className="wiz-wizard__nav__section">
    <li className="wiz-wizard__nav__heading">{title}</li>
    {_.map(section, page => {
      const classes = classNames('wiz-wizard__nav__step', {'wiz-wizard__nav__step--active': page === currentPage});
      const link = <li className={classes} key={page.path}>
        <button
          className="wiz-wizard__nav__link btn btn-link btn-link-ordinary"
          disabled={disabled}
          onClick={() => history.push(page.path)}
        >{page.title}</button>
      </li>;
      disabled = disabled || !canNavigateForward(page);
      return link;
    })}
  </ul>;
});

const NextButton = withNav(
  ({disabled, navNext}) => <div className="withtooltip">
    <button onClick={navNext} className={`btn btn-primary ${disabled ? 'disabled' : ''}`}>Next Step</button>
    {disabled && <div className="tooltip">All fields are required unless specified.</div>}
  </div>
);

const PreviousButton = withNav(
  ({navPrevious}) => <button onClick={navPrevious} className="btn btn-default">Previous Step</button>
);

const ResetButton = () => <button onClick={() => {
  // eslint-disable-next-line no-alert
  (window.config.devMode || window.confirm('Do you really want to start over?')) && window.reset();
}} className="btn btn-link">
  <i className="fa fa-refresh"></i>&nbsp;&nbsp;Start Over
</button>;

const stateToProps = (state, {history}) => {
  const navTrail = trail(state);
  const currentPage = navTrail.pageByPath.get(history.location.pathname);

  return {
    canNavigateForward: page => page.component.canNavigateForward ? page.component.canNavigateForward(state) : true,
    canReset: _.has(currentPage, 'component.canReset') && currentPage.component.canReset(state),
    currentPage,
    navTrail,
    title: `${_.get(currentPage, 'title')}${window.config.devMode ? ' (dev)' : ''}`,
  };
};

// No components have the same path, so this is safe.
const routes = _.uniq(_.flatMap(trailSections));

const Wizard = withRouter(connect(stateToProps)(
  class Wizard_ extends React.Component {
    fixLocation () {
      const {canNavigateForward, history, navTrail} = this.props;
      const page = navTrail.pageByPath.get(history.location.pathname);

      if (!page) {
        // Path does not exist in the current trail, so navigate to the trail's first page instead
        history.push(navTrail.pages[0].path);
        return;
      }

      // If the Next button on a previous page is disabled, you shouldn't be able to see this page. Show the first
      // page, before or including this one, that won't allow forward navigation.
      const fixed = navTrail.pages.find(p => p === page || !canNavigateForward(p));
      if (fixed !== page) {
        history.push(fixed.path);
      }
    }

    componentWillMount () {
      this.fixLocation();
      document.title = `Tectonic - ${this.props.title}`;
    }

    componentWillReceiveProps (nextProps) {
      this.fixLocation();
      if (nextProps.title !== this.props.title) {
        document.title = `Tectonic - ${nextProps.title}`;
      }
    }

    render () {
      const {canNavigateForward, canReset, currentPage, navTrail, title} = this.props;
      if (!currentPage) {
        return null;
      }

      const navProps = {canNavigateForward, currentPage, navTrail};

      return (
        <div className="tectonic">
          <Header />
          <div className="tectonic-installer">
            <div className="wiz-wizard">
              <div className="wiz-wizard__cell wiz-wizard__nav">
                <NavSection
                  {...navProps}
                  title="1. Choose Cluster Type"
                  sections={[trailSections.choose]} />
                <NavSection
                  {...navProps}
                  title="2. Define Cluster"
                  sections={[trailSections.defineBaremetal, trailSections.defineAWS]} />
                <NavSection
                  {...navProps}
                  title="3. Boot Cluster"
                  sections={[trailSections.bootBaremetalTF, trailSections.bootAWSTF, trailSections.bootDryRun]} />
              </div>
              <div className="wiz-wizard__content wiz-wizard__cell">
                <div className="wiz-form__header">
                  <span className="wiz-form__header__title">{title}</span>
                </div>
                <div className="wiz-wizard__content__body">
                  <Switch>
                    {routes.map(r => <Route exact key={r.path} path={r.path} render={() => <r.component />} />)}
                  </Switch>
                </div>
                {currentPage.hidePager || <div className="wiz-form__actions">
                  <div className="wiz-form__actions__prev">
                    {navTrail.previousFrom(currentPage) && <PreviousButton />}
                    {canReset && <ResetButton />}
                  </div>
                  <div className="wiz-form__actions__next">
                    {navTrail.nextFrom(currentPage) && <NextButton disabled={!canNavigateForward(currentPage)} />}
                  </div>
                </div>
                }
              </div>
            </div>
          </div>
          <Footer />
        </div>
      );
    }
  }
));

export const Base = connect(
  ({cluster, serverFacts}) => ({
    loaded: cluster.loaded && serverFacts.loaded,
    failed: serverFacts.error !== null,
  }),
  undefined, // mapDispatchToProps
  undefined, // mergeProps
  {pure: false} // Base isn't pure because Wizard isn't pure
)(({loaded, failed}) => {
  if (!loaded) {
    return <Loader />;
  }
  if (failed) {
    return <div className="wiz-wizard">
      <div className="wiz-wizard__cell wiz-wizard__content">
        The Tectonic Installer has encountered an error. Please contact Tectonic support.
      </div>
    </div>;
  }
  return <Wizard />;
});
