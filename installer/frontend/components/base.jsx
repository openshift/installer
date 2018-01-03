import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { Route, Switch, withRouter } from 'react-router-dom';

import { withNav } from '../nav';
import { sections, trail } from '../trail';

import { Loader } from './loader';
import { PLATFORM_TYPE } from '../cluster-config';
import { TectonicGA } from '../tectonic-ga';
import { Header } from './header';
import { Footer } from './footer';

const NavSection = ({canNavigateForward, currentPage, navigate, navSections, navTrail, title}) => {
  const section = _.find(navTrail.sections, s => navSections.includes(s) && s.includes(currentPage));
  let disabled = false;

  return <ul className="wiz-wizard__nav__section">
    <li className="wiz-wizard__nav__heading">{title}</li>
    {_.map(section, page => {
      const classes = classNames('wiz-wizard__nav__step', {'wiz-wizard__nav__step--active': page === currentPage});
      const link = <li className={classes} key={page.path}>
        <button
          className="wiz-wizard__nav__link btn btn-link btn-link-ordinary"
          disabled={disabled}
          onClick={() => navigate(page)}
        >{page.title}</button>
      </li>;
      disabled = disabled || !canNavigateForward(page);
      return link;
    })}
  </ul>;
};

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
  const t = trail(state);
  const currentPage = t.pageByPath.get(history.location.pathname);
  return {
    canReset: _.has(currentPage, 'component.canReset') && currentPage.component.canReset(state),
    currentPage,
    state,
    t,
    title: `${_.get(currentPage, 'title')}${window.config.devMode ? ' (dev)' : ''}`,
  };
};

// No components have the same path, so this is safe.
const routes = _.uniq(_.flatMap(sections));

const Wizard = withNav(withRouter(connect(stateToProps)(
  class Wizard_ extends React.Component {
    navigate (nextPage, state) {
      const {currentPage, history} = this.props;

      if (currentPage.path === '/define/cluster-type' && nextPage !== currentPage && state) {
        TectonicGA.sendEvent('Platform Selected', 'user input', state.clusterConfig[PLATFORM_TYPE], state.clusterConfig[PLATFORM_TYPE]);
      }

      if (nextPage === currentPage) {
        return;
      }

      if (state) {
        TectonicGA.sendEvent('Page Navigation Next', 'click', 'next on', state.clusterConfig[PLATFORM_TYPE]);
      }
      history.push(nextPage.path);
    }

    canNavigateForward (page) {
      return page.component.canNavigateForward ? page.component.canNavigateForward(this.props.state) : true;
    }

    fixLocation () {
      const {history, t} = this.props;
      const page = t.pageByPath.get(history.location.pathname);

      if (!page) {
        // Path does not exist in the current trail, so navigate to the trail's first page instead
        history.push(t.pages[0].path);
        return;
      }

      // If the Next button on a previous page is disabled, you shouldn't be able to see this page. Show the first
      // page, before or including this one, that won't allow forward navigation.
      const fixed = t.pages.find(p => p === page || !this.canNavigateForward(p)).path;
      if (fixed !== history.location.pathname) {
        history.push(fixed);
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

    render() {
      const {t, canReset, currentPage, title} = this.props;
      if (!currentPage) {
        return null;
      }

      const navProps = {
        canNavigateForward: page => this.canNavigateForward(page),
        currentPage,
        navigate: page => this.navigate(page),
        navTrail: t,
      };

      return (
        <div className="tectonic">
          <Header />
          <div className="tectonic-installer">
            <div className="wiz-wizard">
              <div className="wiz-wizard__cell wiz-wizard__nav">
                <NavSection
                  {...navProps}
                  title="1. Choose Cluster Type"
                  navSections={[sections.choose]} />
                <NavSection
                  {...navProps}
                  title="2. Define Cluster"
                  navSections={[sections.defineBaremetal, sections.defineAWS]} />
                <NavSection
                  {...navProps}
                  title="3. Boot Cluster"
                  navSections={[sections.bootBaremetalTF, sections.bootAWSTF, sections.bootDryRun]} />
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
                    {t.previousFrom(currentPage) && <PreviousButton />}
                    {canReset && <ResetButton />}
                  </div>
                  <div className="wiz-form__actions__next">
                    {t.nextFrom(currentPage) && <NextButton disabled={!this.canNavigateForward(currentPage)} />}
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
)));

export const Base = connect(
  ({cluster, serverFacts}) => {
    return {
      loaded: cluster.loaded && serverFacts.loaded,
      failed: serverFacts.error !== null,
    };
  },
  undefined, // mapDispatchToProps
  undefined, // mergeProps
  {pure: false} // base isn't pure because Wizard isn't pure
)(({loaded, failed}) => {
  if (loaded && !failed) {
    return <Wizard />;
  }

  if (loaded) {
    return (
      <div className="wiz-wizard">
        <div className="wiz-wizard__cell wiz-wizard__content">
          The Tectonic Installer has encountered an error. Please contact
          Tectonic support.
        </div>
      </div>
    );
  }

  return <Loader />;
});
