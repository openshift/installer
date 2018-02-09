exports.command = function (fileSelector, textSelector, validText, invalidText, invalidMsg) {
  this
    .setField(textSelector, 'abc')
    .expectValidationErrorContains(invalidMsg)
    .setFileField(fileSelector, validText)
    .expectNoValidationError()
    .expect.element(textSelector).to.have.value.that.equals(validText);

  this
    .setFileField(fileSelector, invalidText)
    .expectValidationErrorContains(invalidMsg)
    .expect.element(textSelector).to.have.value.that.equals(invalidText);

  this
    .setFileField(fileSelector, validText)
    .expectNoValidationError()
    .expect.element(textSelector).to.have.value.that.equals(validText);

  return this;
};
