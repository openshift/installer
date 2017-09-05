import React from 'react';
import semver from 'semver';
import { Dropdown } from './ui';


const fetchLatestRelease = () => {
  return fetch('/releases/latest',
    { method: 'GET' }).then((response) => {
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

const hasNewVersion = (latestRelease) => {
  const lrMajor = semver.major(latestRelease);
  const gtMajor = semver.major(GIT_TAG);
  if (lrMajor !== gtMajor) {
    return lrMajor > gtMajor;
  }

  const lrMinor = semver.minor(latestRelease);
  const gtMinor = semver.minor(GIT_TAG);
  if (lrMinor !== gtMinor) {
    return lrMinor > gtMinor;
  }

  const lrPatch = semver.patch(latestRelease);
  const gtPatch = semver.patch(GIT_TAG);
  if (lrPatch !== gtPatch) {
    return lrPatch > gtPatch;
  }

  const latestReleasePr = semver.prerelease(latestRelease);
  const gitTagPr = semver.prerelease(GIT_TAG);

  //No rc string in the latest release
  if (latestReleasePr.length < gitTagPr.length) {
    return true;
  }

  //compares rc-x string
  if (latestReleasePr[1] !== gitTagPr[1]) {
    return latestReleasePr[1] > gitTagPr[1];
  }

  //rc version
  if (latestReleasePr[2] !== gitTagPr[2]) {
    return latestReleasePr[2] > gitTagPr[2];
  }

  return false;
};

export class Header extends React.Component {
  constructor (props) {
    super(props);
    this.state = { latestRelease: null };
  }

  componentDidMount() {
    fetchLatestRelease().then((release) => {
      this.setState({ latestRelease: parseLatestVersion(release) });
    }).catch((err) => {
      console.error('Error retrieving latest version of tectonic ', err.message);
      this.setState({ latestRelease: null });
    });

  }

  render() {
    const latestRelease = this.state.latestRelease || null;

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
            <Dropdown items={productDdItems} header="Product" />
          </li>
          <li>
            <Dropdown items={openSourceDdItems} header="Open Source" />
          </li>
          <li className="tectonic-dropdown-menu-title">
            {/* eslint-disable react/jsx-no-target-blank */}
            <a href="https://coreos.com/tectonic/docs/latest/" rel="noopener" target="_blank" className="tectonic-dropdown-menu-title__link">Tectonic Docs</a>
            {/* eslint-enable react/jsx-no-target-blank */}
          </li>
        </ul>
        <div className="co-navbar--right">
          <ul className="co-navbar-nav">
            {latestRelease && hasNewVersion(latestRelease) && <li className="co-navbar-nav-item__version">
              <span className="co-navbar-nav-item__version--new">
                {/* eslint-disable react/jsx-no-target-blank */}
                New installer version: <a href="https://coreos.com/tectonic/releases/" rel="noopener" target="_blank">Release notes {latestRelease}</a>
                {/* eslint-enable react/jsx-no-target-blank */}
              </span>
            </li>}
            <li className="co-navbar-nav-item__version">
              <span>Version: {GIT_TAG} ({GIT_COMMIT})</span>
            </li>
          </ul>
        </div>
      </div>
    </div>;
  }
}
