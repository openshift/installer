import _ from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import semver from 'semver';

import { A, DocsA, DropdownMixin } from './ui';

const fetchLatestRelease = () => {
  const opts = {
    credentials: 'same-origin',
    method: 'GET',
  };
  return fetch('/releases/latest', opts).then((response) => {
    if (response.ok) {
      return response.text();
    }
    return response.text().then(err => Promise.reject(err));
  });
};

const parseLatestVersion = (html) => {
  const parser = new DOMParser();
  const htmlDoc = parser.parseFromString(html, 'text/html');
  return htmlDoc.getElementsByClassName('latestReleaseTag')[0].textContent.trim();
};

const hasNewVersion = (currentRelease, latestRelease) => {
  if (!semver.valid(latestRelease)) {
    return false;
  }
  if (!semver.valid(currentRelease)) {
    return true;
  }
  return latestRelease > currentRelease;
};

class MenuDropdown extends DropdownMixin {
  render () {
    const {active} = this.state;
    const {items, header} = this.props;

    const children = _.map(items, (href, title) => {
      const rel = href.includes('coreos.com') ? 'noopener' : 'noopener noreferrer';
      return <li className="tectonic-dropdown-menu-item" key={title}>
        <A className="tectonic-dropdown-menu-item__link" href={href} rel={rel}>{title}</A>
      </li>;
    });

    return (
      <div ref={el => this.dropdownElement = el} className="dropdown" onClick={this.toggle}>
        <a className="tectonic-dropdown-menu-title">{header}&nbsp;&nbsp;<i className="fa fa-angle-down"></i></a>
        <ul className="dropdown-menu tectonic-dropdown-menu" style={{display: active ? 'block' : 'none'}}>
          {children}
        </ul>
      </div>
    );
  }
}

export const Header = connect(
  ({serverFacts: {buildTime, version}}) => ({buildTime, version})
)(class Header_ extends React.Component {
  constructor (props) {
    super(props);
    this.state = { latestRelease: null };
  }

  componentDidMount () {
    fetchLatestRelease().then((release) => {
      this.setState({ latestRelease: parseLatestVersion(release) });
    }).catch((err) => {
      console.error('Error retrieving latest version of Tectonic ', err.message);
      this.setState({ latestRelease: null });
    });
  }

  render () {
    const {buildTime, version} = this.props;
    const {latestRelease} = this.state;

    const productDdItems = {
      'Tectonic - Kubernetes': 'https://coreos.com/tectonic/',
      'Quay - Registry': 'https://coreos.com/quay-enterprise',
      'Container Linux Support': 'https://coreos.com/products/container-linux-subscription/',
      'Training': 'https://coreos.com/training/',
    };

    const openSourceDdItems = {
      'Open Source Docs': 'https://coreos.com/docs/',
      'Kubernetes': 'https://coreos.com/kubernetes',
      'Operators': 'https://coreos.com/operators/',
      'Container Linux': 'https://coreos.com/os/docs/latest',
      'rkt': 'https://coreos.com/rkt/',
      'etcd': 'https://coreos.com/etcd/',
      'Clair': 'https://coreos.com/clair/',
      'flannel': 'https://coreos.com/flannel/',
      'Ignition': 'https://coreos.com/ignition/',
      'Matchbox': 'https://coreos.com/matchbox/',
      '90+ more on Github': 'https://github.com/coreos/',
    };

    return <div className="co-navbar">
      <div className="co-navbar__header">
        <a href="/" className="co-navbar__logo-link">
          <img className="co-navbar__logo" src="/frontend/lib/mochi/img/logo/coreos/logo.svg" />
        </a>
      </div>
      <div>
        <ul className="co-navbar-nav">
          <li>
            <MenuDropdown items={productDdItems} header="Product" />
          </li>
          <li>
            <MenuDropdown items={openSourceDdItems} header="Open Source" />
          </li>
          <li className="tectonic-dropdown-menu-title">
            <DocsA className="tectonic-dropdown-menu-title__link" path="/">Tectonic Docs</DocsA>
          </li>
        </ul>
        <div className="co-navbar--right">
          <ul className="co-navbar-nav">
            {hasNewVersion(version, latestRelease) && <li className="co-navbar-nav-item__version">
              <span className="co-navbar-nav-item__version--new">
                New installer version: <A href="https://coreos.com/tectonic/releases/" rel="noopener">Release notes {latestRelease}</A>
              </span>
            </li>}
            {version && <li className="co-navbar-nav-item__version" title={buildTime && (new Date(buildTime)).toString()}>Version: {version}</li>}
          </ul>
        </div>
      </div>
    </div>;
  }
});
