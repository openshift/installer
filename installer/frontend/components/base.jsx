import _ from 'lodash';
import classNames from 'classnames';
import React from 'react';
import { connect } from 'react-redux';
import { Route, Switch, withRouter } from 'react-router-dom';

import { fixLocation } from '../app';
import { withNav } from '../nav';
import { sections as trailSections, trail } from '../trail';

import { Loader } from './loader';
import { PLATFORM_TYPE } from '../cluster-config';
import { TectonicGA } from '../tectonic-ga';
import { Header } from './header';
import { Footer } from './footer';

const NavSection = connect(state => ({state}))(
  ({title, navTrail, sections, currentPage, handlePage, state}) => {
    const currentSection = navTrail.sectionFor(currentPage);
    const section = sections.find(s => s === currentSection);

    return (
      <ul className="wiz-wizard__nav__section">
        <li className="wiz-wizard__nav__heading">{title}</li>
        {
          section &&
          section.map(page => {
            const classes = classNames('wiz-wizard__nav__step', {
              'wiz-wizard__nav__step--active': page === currentPage,
            });
            return (
              <li className={classes} key={page.path}>
                <button className="wiz-wizard__nav__link btn btn-link btn-link-ordinary"
                  onClick={() => handlePage(page)}
                  disabled={!navTrail.navigable(page) || !navTrail.canNavigateForward(currentPage, page, state)}
                >{page.title}</button>
              </li>
            );
          })
        }
      </ul>
    );
  }
);

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
    currentPage,
    state,
    t,
    title: `${_.get(currentPage, 'title')}${window.config.devMode ? ' (dev)' : ''}`,
  };
};

// No components have the same path, so this is safe.
const routes = _.uniq(_.flatMap(trailSections));

const Wizard = withNav(withRouter(connect(stateToProps)(
  class extends React.Component {
    navigate (currentPage, nextPage, state) {
      if (currentPage.path === '/define/cluster-type' && nextPage !== currentPage && state) {
        TectonicGA.sendEvent('Platform Selected', 'user input', state.clusterConfig[PLATFORM_TYPE], state.clusterConfig[PLATFORM_TYPE]);
      }

      if (nextPage === currentPage) {
        return;
      }

      if (state) {
        TectonicGA.sendEvent('Page Navigation Next', 'click', 'next on', state.clusterConfig[PLATFORM_TYPE]);
      }
      this.props.history.push(nextPage.path);
    }

    componentDidMount() {
      document.title = `Tectonic - ${this.props.title}`;
    }

    componentWillReceiveProps (nextProps) {
      if (!nextProps.currentPage) {
        fixLocation();
      }
      if (nextProps.title === this.props.title) {
        return;
      }
      document.title = `Tectonic - ${nextProps.title}`;
    }

    render() {
      const {t, currentPage, state, title} = this.props;
      if (!currentPage) {
        return null;
      }

      const nav = page => this.navigate(currentPage, page);

      const {canNavigateForward, canReset} = currentPage.component;
      const disableNext = canNavigateForward ? !canNavigateForward(state) : false;

      return (
        <div className="tectonic">
          <Header />
          <div className="tectonic-installer">
            <div className="wiz-wizard">
              <div className="wiz-wizard__cell wiz-wizard__nav">
                <NavSection
                  title="1. Choose Cluster Type"
                  navTrail={t}
                  sections={[trailSections.choose]}
                  currentPage={currentPage}
                  handlePage={nav} />
                <NavSection
                  title="2. Define Cluster"
                  navTrail={t}
                  sections={[trailSections.defineBaremetal, trailSections.defineAWS]}
                  currentPage={currentPage}
                  handlePage={nav} />
                <NavSection
                  title="3. Boot Cluster"
                  navTrail={t}
                  sections={[
                    trailSections.bootBaremetalTF,
                    trailSections.bootAWSTF,
                    trailSections.bootDryRun,
                  ]}
                  currentPage={currentPage}
                  handlePage={nav} />
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
                    {canReset && canReset(state) && <ResetButton />}
                  </div>
                  <div className="wiz-form__actions__next">
                    {t.nextFrom(currentPage) && <NextButton disabled={disableNext} />}
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
