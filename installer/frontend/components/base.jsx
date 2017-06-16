import classNames from 'classnames';
import React from 'react';
import { saveAs } from 'file-saver';
import { connect } from 'react-redux';

import { savable } from '../reducer';
import * as trail from '../trail';

import { Loader } from './loader';
import { ResetButton } from './reset-button';
import { restoreModal } from './restore';
import { WithTooltip } from './tooltip';
import { PLATFORM_TYPE } from '../cluster-config';
import { TectonicGA } from '../tectonic-ga';
import { Header } from './header';
import { Footer } from './footer';

const downloadState = (state) => {
  const toSave = savable(state);
  const saved = JSON.stringify(toSave, null, 2);
  const stateBlob = new Blob([saved], {type: 'application/json'});
  saveAs(stateBlob, 'tectonic.progress');
  TectonicGA.sendEvent('Installer Link', 'click', 'User downloads progress file', state.clusterConfig[PLATFORM_TYPE]);
};

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
                <button className='wiz-wizard__nav__link btn btn-link btn-link-ordinary'
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

const Pager = ({showPrev, showNext, disableNext, loadingNext, navigatePrevious, navigateNext, resetBtn}) => {
  const nextLinkClasses = classNames('btn', 'btn-primary', {
    disabled: disableNext || loadingNext,
  });

  return (
    <div className="wiz-form__actions">
      {
        showPrev &&
        <button onClick={navigatePrevious}
                className="btn btn-default wiz-form__actions__prev"
                >Previous Step</button>
      }
      { resetBtn && <div className="wiz-form__actions__prev">
        <ResetButton />
      </div>
      }
      {
        showNext &&
          <div className="wiz-form__actions__next">
            <WithTooltip text="All fields are required unless specified." shouldShow={disableNext}>
              <button onClick={navigateNext}
                      className={nextLinkClasses}>
                {loadingNext &&
                 <span><i className="fa fa-spin fa-circle-o-notch"></i>{' '}</span>}
                 Next Step
              </button>
            </WithTooltip>
          </div>
      }
    </div>
  );
};

const stateToProps = (state) => {
  const t = trail.trail(state);
  const currentPage = t.pageByPath.get(state.path);
  return {
    currentPage,
    nextPage: t.nextFrom(currentPage),
    prevPage: t.previousFrom(currentPage),
    state,
    t,
    title: `${currentPage.title}${window.config.devMode ? ' (dev)' : ''}`,
  };
};

const Wizard = connect(stateToProps)(
class extends React.Component {
  static get contextTypes() {
    return {
      router: React.PropTypes.object.isRequired,
    };
  }

  navigate (currentPage, nextPage, state) {
    const {router} = this.context;

    if (currentPage.path === '/define/cluster-type' && nextPage !== currentPage && state) {
      TectonicGA.sendEvent('Platform Selected', 'user input', state.clusterConfig[PLATFORM_TYPE], state.clusterConfig[PLATFORM_TYPE]);
    }

    if (nextPage === currentPage) {
      return;
    }

    if (state) {
      TectonicGA.sendEvent('Page Navigation Next', 'click', 'next on', state.clusterConfig[PLATFORM_TYPE]);
    }
    router.push(nextPage.path);
  }

  componentDidMount() {
    document.title = `Tectonic - ${this.props.title}`;
  }

  componentWillReceiveProps (nextProps) {
    if (nextProps.title === this.props.title) {
      return;
    }
    document.title = `Tectonic - ${nextProps.title}`;
  }

  render() {
    const {children, t, currentPage, prevPage, nextPage, state} = this.props;

    const navigatePrevious = () => this.navigate(currentPage, prevPage, state);
    const navigateNext = () => this.navigate(currentPage, nextPage, state);
    const nav = page => this.navigate(currentPage, page);

    const kids = React.Children.map(children, el => {
      return React.cloneElement(el, {navigatePrevious, navigateNext});
    });

    const canNavigateForward = currentPage.component.canNavigateForward || (() => true);
    return (
      <div className="tectonic">
        <Header />
        <div className="tectonic-installer">
          <div className="wiz-wizard">
            <div className="wiz-wizard__cell wiz-wizard__nav">
              <NavSection
                  title="1. Choose Cluster Type"
                  navTrail={t}
                  sections={[trail.sections.choose]}
                  currentPage={currentPage}
                  handlePage={nav} />
              <NavSection
                  title="2. Define Cluster"
                  navTrail={t}
                  sections={[trail.sections.defineBaremetal, trail.sections.defineAWS]}
                  currentPage={currentPage}
                  handlePage={nav} />
              <NavSection
                  title="3. Boot Cluster"
                  navTrail={t}
                  sections={[
                    trail.sections.bootBaremetal,
                    trail.sections.bootAWS,
                    trail.sections.bootAWSTF,
                    trail.sections.bootDryRun,
                  ]}
                  currentPage={currentPage}
                  handlePage={nav} />
            </div>
            <div className="wiz-wizard__content wiz-wizard__cell">
              <div className="wiz-form__header">
                <span className="wiz-form__header__title">{this.props.title}</span>
                {currentPage.showRestore &&
                  <span className="wiz-form__header__control">
                    <a onClick={restoreModal}><i className="fa fa-upload"></i>&nbsp;&nbsp;Restore progress</a>
                  </span>
                }
                {currentPage.hideSave ||
                 <span className="wiz-form__header__control">
                   <a onClick={() => downloadState(state)}><i className="fa fa-download"></i>&nbsp;&nbsp;Save progress</a>
                 </span>
                }
              </div>
              <div className="wiz-wizard__content__body">
                {kids}
              </div>
              {
                currentPage.hidePager ||
                <Pager
                    showPrev={!!prevPage}
                    showNext={!!nextPage}
                    disableNext={!canNavigateForward(state)}
                    navigatePrevious={navigatePrevious}
                    resetBtn={t.canReset}
                    navigateNext={navigateNext} />
              }
            </div>
          </div>
        </div>
        <Footer />
      </div>
    );
  }
});

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
)(({loaded, failed, children}) => {
  if (loaded && !failed) {
    return <Wizard>{children}</Wizard>;
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
