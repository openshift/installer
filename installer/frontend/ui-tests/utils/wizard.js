const nextStep = '.wiz-form__actions__next button.btn-primary';
const prevStep = 'button.wiz-form__actions__prev';

const testPage = (page, json, nextInitiallyDisabled = true) => {
  page.expect.element(prevStep).to.be.visible;
  page.expect.element(prevStep).to.not.have.attribute('class').which.contains('disabled');

  page.expect.element(nextStep).to.be.visible;
  if (nextInitiallyDisabled) {
    page.expect.element(nextStep).to.have.attribute('class').which.contains('disabled');
  }
  if (page.test) {
    page.test(json);
  }
  page.waitForElementNotPresent(`${nextStep}.disabled`, 10000);
  page.expect.element(nextStep).to.have.attribute('class').which.not.contains('disabled');
  return page.click(nextStep);
};

module.exports = {
  nextStep,
  prevStep,
  testPage,
};
