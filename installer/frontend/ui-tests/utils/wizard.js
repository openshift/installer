const nextStep = '.wiz-form__actions__next button.btn-primary';
const prevStep = 'button.wiz-form__actions__prev';

const testPage = (page, platform, json, nextInitiallyDisabled = true) => {
  page.expect.element(prevStep).to.be.visible;
  page.expect.element(prevStep).to.have.attribute('class').which.not.contains('disabled');

  page.expect.element(nextStep).to.be.visible;
  if (nextInitiallyDisabled) {
    page.expect.element(nextStep).to.have.attribute('class').which.contains('disabled');
  }

  // Sidebar link for the current page should be highlighted and enabled
  page.expect.element('.wiz-wizard__nav__step--active button').to.not.have.attribute('disabled');

  // The next sidebar links should be disabled if the next button is disabled
  const nextNavLink = page.expect.element('.wiz-wizard__nav__step--active + .wiz-wizard__nav__step button');
  if (nextInitiallyDisabled) {
    nextNavLink.to.have.attribute('disabled');

    // If the next button is disabled, all sidebar links for later screens should be disabled
    page.expect.element('.wiz-wizard__nav__step--active ~ .wiz-wizard__nav__step button:not([disabled])').to.not.be.present;
  } else {
    nextNavLink.to.not.have.attribute('disabled');

    // If the next button is enabled, the next sidebar link should be enabled too
    page.expect.element('.wiz-wizard__nav__step--active + .wiz-wizard__nav__step button').to.not.have.attribute('disabled');
  }

  // Save progress link exists
  page.expect.element('.wiz-form__header a').text.which.contains('Save progress');

  if (page.test) {
    page.test(json, platform);
  }
  page.expect.element(nextStep).to.have.attribute('class').which.not.contains('disabled');
  return page.click(nextStep);
};

module.exports = {
  nextStep,
  prevStep,
  testPage,
};
